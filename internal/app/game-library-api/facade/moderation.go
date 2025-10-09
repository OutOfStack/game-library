package facade

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/client/openaiapi"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"go.uber.org/zap"
)

const (
	maxModerationAttempts = 5
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
		Screenshots: g.Screenshots,
		Websites:    g.Websites,
	}, nil
}

// ProcessModeration processes moderation for a game using a two-phase approach:
// 1. Basic moderation: Uses OpenAI moderation API to check for policy violations in text and images
// 2. Gaming-specific analysis: Uses vision model to analyze images for gaming appropriateness and relevance
func (p *Provider) ProcessModeration(ctx context.Context, gameID int32) error {
	p.log.Info("processing moderation for game", zap.Int32("game_id", gameID))

	moderation, err := p.storage.GetModerationRecordByGameID(ctx, gameID)
	if err != nil {
		return fmt.Errorf("get moderation record: %w", err)
	}
	if moderation.Attempts >= maxModerationAttempts {
		p.log.Error("exceeded maximum moderation attempts", zap.Int32("game_id", gameID))
		return p.storage.SetModerationRecordsStatus(ctx, []int32{moderation.ID}, model.ModerationStatusFailed)
	}

	// get game data for moderation
	game, err := p.storage.GetGameByID(ctx, gameID)
	if err != nil {
		return fmt.Errorf("get game %d: %w", gameID, err)
	}

	moderationData, err := p.mapGameToModerationData(ctx, &game)
	if err != nil {
		return fmt.Errorf("map game %d to moderation data: %w", gameID, err)
	}

	// phase 1: Basic moderation with OpenAI moderation API
	moderationResp, err := p.openAIClient.ModerateText(ctx, moderationData)
	if err != nil {
		p.log.Error("openai moderation api error", zap.Error(err), zap.Int32("game_id", gameID))
		return fmt.Errorf("moderation api: %w", err)
	}

	// check if basic moderation flagged content
	if hasViolations(moderationResp) {
		details := getViolationDetails(moderationResp)
		violationTypes := getViolationTypes(moderationResp)

		p.log.Info("game moderation declined - policy violations",
			zap.Int32("game_id", gameID),
			zap.Strings("violations", violationTypes))

		return p.saveModerationResult(ctx, gameID, model.ModerationStatusDeclined, "Content violates safety policies", details, violationTypes)
	}

	// phase 2: Gaming-specific image analysis
	if len(moderationData.Screenshots) > 0 || moderationData.LogoURL != "" {
		visionResult, vErr := p.openAIClient.AnalyzeGameImages(ctx, moderationData)
		if vErr != nil {
			return fmt.Errorf("image analysis failed for game %d: %w", gameID, vErr)
		}
		if !visionResult.Approved || !visionResult.GamingAppropriate || !visionResult.ContentRelevant {
			details := fmt.Sprintf("Gaming appropriate: %t, Content relevant: %t",
				visionResult.GamingAppropriate, visionResult.ContentRelevant)

			p.log.Info("game moderation declined - image analysis",
				zap.Int32("game_id", gameID),
				zap.String("reason", visionResult.Reason))

			return p.saveModerationResult(ctx, gameID, model.ModerationStatusDeclined, visionResult.Reason, details, nil)
		}
	}

	// all checks passed - approve
	p.log.Info("game moderation approved", zap.Int32("game_id", gameID))

	err = p.saveModerationResult(ctx, gameID, model.ModerationStatusReady, "Content approved", "All moderation checks passed", nil)
	if err != nil {
		return fmt.Errorf("save moderation result for game %d: %w", gameID, err)
	}

	// invalidate cache on successful moderation
	go func() {
		bCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()

		// invalidate game in case moderation happens after game update and original game data is still cached
		key := getGameKey(gameID)
		err = cache.Delete(bCtx, p.cache, key)
		if err != nil {
			p.log.Error("remove game cache by key", zap.String("key", key), zap.Error(err))
		}
		// recache game
		err = cache.Get(bCtx, p.cache, key, new(model.Game), func() (model.Game, error) {
			return p.storage.GetGameByID(bCtx, gameID)
		}, 0)
		if err != nil {
			p.log.Error("cache game with id", zap.Int32("id", gameID), zap.Error(err))
		}
	}()

	return nil
}

// saveModerationResult saves moderation result to database
func (p *Provider) saveModerationResult(ctx context.Context, gameID int32, status model.ModerationStatus,
	reason, details string, violationTypes []string) error {
	resultDetails := fmt.Sprintf("%s. %s", reason, details)
	if len(violationTypes) > 0 {
		resultDetails += " Violations: " + strings.Join(violationTypes, ", ")
	}

	err := p.storage.SetModerationRecordResultByGameID(ctx, gameID, model.UpdateModerationResult{
		ResultStatus: status,
		Details:      resultDetails,
	})
	if err != nil {
		return fmt.Errorf("update moderation result for game %d: %w", gameID, err)
	}

	return nil
}

// hasViolations checks if moderation response has any flagged content
func hasViolations(resp *openaiapi.ModerationResponse) bool {
	for _, result := range resp.Results {
		if result.Flagged {
			return true
		}
	}
	return false
}

// getViolationDetails returns details about violations
func getViolationDetails(resp *openaiapi.ModerationResponse) string {
	var violations []string
	inputTypes := []string{"name", "summary", "developers", "publisher", "websites", "logo", "screenshot"}

	for i, result := range resp.Results {
		if result.Flagged {
			inputType := "content"
			if i < len(inputTypes) {
				inputType = inputTypes[i]
			}

			var categories []string
			for category, flagged := range result.Categories {
				if flagged {
					categories = append(categories, category)
				}
			}
			violations = append(violations, fmt.Sprintf("%s: %s", inputType, strings.Join(categories, ", ")))
		}
	}
	return strings.Join(violations, "; ")
}

// getViolationTypes returns array of violation types
func getViolationTypes(resp *openaiapi.ModerationResponse) []string {
	types := make(map[string]bool)
	for _, result := range resp.Results {
		if result.Flagged {
			for category, flagged := range result.Categories {
				if flagged {
					types[category] = true
				}
			}
		}
	}

	result := make([]string, 0, len(types))
	for category := range types {
		result = append(result, category)
	}
	return result
}
