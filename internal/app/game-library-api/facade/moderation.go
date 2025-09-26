package facade

import (
	"context"
	"fmt"
	"slices"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
)

// CreateModerationRecord creates a moderation record for a game
func (p *Provider) CreateModerationRecord(ctx context.Context, gameID int32) (int32, error) {
	// get game
	game, err := p.storage.GetGameByID(ctx, gameID)
	if err != nil {
		return 0, fmt.Errorf("get game %d: %w", gameID, err)
	}

	moderationData, err := p.mapGameToModerationData(ctx, &game)
	if err != nil {
		return 0, fmt.Errorf("get game %d moderation data: %w", gameID, err)
	}

	var moderationID int32
	txErr := p.storage.RunWithTx(ctx, func(ctx context.Context) error {
		// create moderation record
		moderationID, err = p.storage.CreateModerationRecord(ctx, model.NewCreateModeration(gameID, moderationData))
		if err != nil {
			return fmt.Errorf("create moderation record for game %d: %w", gameID, err)
		}

		// set moderation id to game
		err = p.storage.UpdateGameModerationID(ctx, gameID, moderationID)
		if err != nil {
			return fmt.Errorf("update game %d moderation id: %w", gameID, err)
		}

		return nil
	})
	if txErr != nil {
		return 0, txErr
	}

	return moderationID, nil
}

// GetGameModerations returns all moderations for a game, ensuring the caller is its publisher
func (p *Provider) GetGameModerations(ctx context.Context, gameID int32, publisher string) ([]model.Moderation, error) {
	game, err := p.storage.GetGameByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	publisherID, err := p.storage.GetCompanyIDByName(ctx, publisher)
	if err != nil {
		return nil, fmt.Errorf("get publisher id by name %s: %w", publisher, err)
	}
	if !slices.Contains(game.PublishersIDs, publisherID) {
		return nil, apperr.NewForbiddenError("game", gameID)
	}
	list, err := p.storage.GetModerationRecordsByGameID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("get moderations by game %d: %w", gameID, err)
	}
	return list, nil
}

// mapGameToModerationData maps game to moderation data
func (p *Provider) mapGameToModerationData(ctx context.Context, g *model.Game) (model.ModerationData, error) {
	if len(g.PublishersIDs) == 0 {
		return model.ModerationData{}, fmt.Errorf("no publisher in game %d", g.ID)
	}

	companies, err := p.GetCompaniesMap(ctx)
	if err != nil {
		return model.ModerationData{}, err
	}
	developers := make([]string, len(g.DevelopersIDs))
	for i, id := range g.DevelopersIDs {
		developers[i] = companies[id].Name
	}

	genresMap, err := p.GetGenresMap(ctx)
	if err != nil {
		return model.ModerationData{}, err
	}
	genres := make([]string, len(g.GenresIDs))
	for i, id := range g.GenresIDs {
		genres[i] = genresMap[id].Name
	}

	return model.ModerationData{
		Name:        g.Name,
		Developers:  developers,
		Publisher:   companies[g.PublishersIDs[0]].Name,
		ReleaseDate: g.ReleaseDate.String(),
		Genres:      genres,
		LogoURL:     g.LogoURL,
		Summary:     g.Summary,
		Slug:        g.Slug,
		Screenshots: g.Screenshots,
		Websites:    g.Websites,
	}, nil
}
