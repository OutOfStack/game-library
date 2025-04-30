package facade

import (
	"context"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"go.uber.org/zap"
)

const (
	// MaxGamesPerPublisherPerMonth is the maximum number of games a publisher can create in a month
	MaxGamesPerPublisherPerMonth = 2
)

// GetGames returns games and count with pagination
func (p *Provider) GetGames(ctx context.Context, page, pageSize int, filter model.GamesFilter) (games []model.Game, count uint64, err error) {
	err = cache.Get(ctx, p.cache, getGamesKey(int64(pageSize), int64(page), filter), &games, func() ([]model.Game, error) {
		return p.storage.GetGames(ctx, pageSize, page, filter)
	}, 0)
	if err != nil {
		return nil, 0, fmt.Errorf("get games: %w", err)
	}

	err = cache.Get(ctx, p.cache, getGamesCountKey(filter), &count, func() (uint64, error) {
		return p.storage.GetGamesCount(ctx, filter)
	}, 0)
	if err != nil {
		return nil, 0, fmt.Errorf("get games count: %w", err)
	}

	return games, count, nil
}

// GetGameByID returns game by id
func (p *Provider) GetGameByID(ctx context.Context, id int32) (model.Game, error) {
	var game model.Game
	err := cache.Get(ctx, p.cache, getGameKey(id), &game, func() (model.Game, error) {
		return p.storage.GetGameByID(ctx, id)
	}, 0)
	if err != nil {
		if apperr.IsStatusCode(err, apperr.NotFound) {
			return model.Game{}, err
		}
		return model.Game{}, fmt.Errorf("get game by id %d: %w", id, err)
	}

	return game, nil
}

// CreateGame creates new game
func (p *Provider) CreateGame(ctx context.Context, cg model.CreateGame) (id int32, err error) {
	// get developer id or create developer
	developerID, err := p.storage.GetCompanyIDByName(ctx, cg.Developer)
	if err != nil && !apperr.IsStatusCode(err, apperr.NotFound) {
		return 0, fmt.Errorf("get company id with name %s: %w", cg.Developer, err)
	}
	if developerID == 0 {
		developerID, err = p.storage.CreateCompany(ctx, model.Company{
			Name: cg.Developer,
		})
		if err != nil {
			return 0, fmt.Errorf("create company %s: %w", cg.Developer, err)
		}
	}

	// get publisher id or create publisher
	publisherID, err := p.storage.GetCompanyIDByName(ctx, cg.Publisher)
	if err != nil && !apperr.IsStatusCode(err, apperr.NotFound) {
		return 0, fmt.Errorf("get company id by name %s: %w", cg.Publisher, err)
	}
	if publisherID == 0 {
		publisherID, err = p.storage.CreateCompany(ctx, model.Company{
			Name: cg.Publisher,
		})
		if err != nil {
			return 0, fmt.Errorf("create company %s: %w", cg.Publisher, err)
		}
	}

	// check if publisher has reached the monthly limit
	if err = p.checkPublisherMonthlyLimit(ctx, publisherID); err != nil {
		return 0, err
	}

	create := cg.MapToCreateGameData(publisherID, developerID)

	id, err = p.storage.CreateGame(ctx, create)
	if err != nil {
		return 0, fmt.Errorf("add new game: %w", err)
	}

	// invalidate cache
	go func() {
		bCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
		defer cancel()

		// invalidate games cache
		key := gamesKey
		err = cache.DeleteByStartsWith(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}
		// invalidate games count cache
		key = gamesCountKey
		err = cache.DeleteByStartsWith(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}
		// cache game
		key = getGameKey(id)
		err = cache.Get(bCtx, p.cache, key, new(model.Game), func() (model.Game, error) {
			return p.storage.GetGameByID(bCtx, id)
		}, 0)
		if err != nil {
			p.log.Error("cache game with id", zap.Int32("id", id), zap.Error(err))
		}
		// invalidate companies
		key = getCompaniesKey()
		err = cache.Delete(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove companies cache", zap.String("key", key), zap.Error(err))
		}
	}()

	return id, nil
}

// UpdateGame updates game
func (p *Provider) UpdateGame(ctx context.Context, id int32, upd model.UpdateGame) error {
	game, err := p.storage.GetGameByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get game by id %d: %w", id, err)
	}

	// check game ownership by publisher
	publisherID, err := p.storage.GetCompanyIDByName(ctx, upd.Publisher)
	if err != nil {
		return fmt.Errorf("get company id by name %s: %w", upd.Publisher, err)
	}

	if len(game.PublishersIDs) != 1 || game.PublishersIDs[0] != publisherID {
		return apperr.NewForbiddenError("game", id)
	}

	developer := upd.Developer
	developersIDs := game.DevelopersIDs
	if developer != nil {
		if *developer == "" {
			developersIDs = []int32{}
		} else {
			// get id or create developer
			developerID, cErr := p.storage.GetCompanyIDByName(ctx, *developer)
			if cErr != nil && !apperr.IsStatusCode(cErr, apperr.NotFound) {
				return fmt.Errorf("get developer id by name %s: %v", *developer, cErr)
			}
			if developerID == 0 {
				developerID, err = p.storage.CreateCompany(ctx, model.Company{
					Name: *developer,
				})
				if err != nil {
					return fmt.Errorf("create developer %s: %v", *developer, err)
				}
			}
			developersIDs = []int32{developerID}
		}
	}

	update := upd.MapToUpdateGameData(game, developersIDs)

	err = p.storage.UpdateGame(ctx, id, update)
	if err != nil {
		if apperr.IsStatusCode(err, apperr.NotFound) {
			return err
		}
		return fmt.Errorf("update game with id %v: %v", id, err)
	}

	// invalidate cache
	go func() {
		bCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
		defer cancel()

		// invalidate games cache
		key := gamesKey
		err = cache.DeleteByStartsWith(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}

		// invalidate game cache
		key = getGameKey(id)
		err = cache.Delete(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by key", zap.String("key", key), zap.Error(err))
		}
		// recache game
		err = cache.Get(bCtx, p.cache, key, new(model.Game), func() (model.Game, error) {
			return p.storage.GetGameByID(bCtx, id)
		}, 0)
		if err != nil {
			p.log.Error("recache game with id", zap.Int32("id", id), zap.Error(err))
		}

		// invalidate companies
		key = getCompaniesKey()
		err = cache.Delete(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove companies cache", zap.String("key", key), zap.Error(err))
		}
	}()

	return nil
}

// DeleteGame deletes game by id
func (p *Provider) DeleteGame(ctx context.Context, id int32, publisher string) error {
	// check game ownership by publisher
	publisherID, err := p.storage.GetCompanyIDByName(ctx, publisher)
	if err != nil {
		return fmt.Errorf("get company id by name %s: %w", publisher, err)
	}

	game, err := p.storage.GetGameByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get game by id %d: %w", id, err)
	}

	if len(game.PublishersIDs) != 1 || game.PublishersIDs[0] != publisherID {
		return apperr.NewForbiddenError("game", id)
	}

	err = p.storage.DeleteGame(ctx, id)
	if err != nil {
		if apperr.IsStatusCode(err, apperr.NotFound) {
			return err
		}
		return fmt.Errorf("delete game with id %v: %v", id, err)
	}

	// invalidate cache
	go func() {
		bCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
		defer cancel()

		// invalidate games cache
		key := gamesKey
		err = cache.DeleteByStartsWith(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}
		// invalidate games count cache
		key = gamesCountKey
		err = cache.DeleteByStartsWith(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove cache by matching key", zap.String("key", key), zap.Error(err))
		}
		// invalidate game cache
		key = getGameKey(id)
		err = cache.Delete(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove game cache by key", zap.String("key", key), zap.Error(err))
		}
	}()

	return nil
}

// Checks if a publisher has reached the monthly limit for creating games
func (p *Provider) checkPublisherMonthlyLimit(ctx context.Context, publisherID int32) error {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Millisecond)

	count, err := p.storage.GetPublisherGamesCount(ctx, publisherID, startOfMonth, endOfMonth)
	if err != nil {
		return fmt.Errorf("check publisher monthly limit: %v", err)
	}

	if count >= MaxGamesPerPublisherPerMonth {
		return apperr.NewTooManyRequestsError("game", fmt.Sprintf("publishing monthly limit of %d reached", MaxGamesPerPublisherPerMonth))
	}

	return nil
}
