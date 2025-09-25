package taskprocessor_test

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/OutOfStack/game-library/internal/client/s3"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"go.uber.org/mock/gomock"
)

func (s *TestSuite) TestStartFetchIGDBGames_Success() {
	lastReleasedAt := time.Now()

	task := model.Task{
		Name:     "fetch_igdb_games",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: []byte(fmt.Sprintf(`{"lastReleasedAt":"%s"}`, lastReleasedAt.Format(time.RFC3339))),
	}

	platforms := []model.Platform{
		{ID: td.Int31(), IGDBID: td.Int64()},
	}

	logoURL, logoFileName, screenshotURL, screenshotFileName, contentType := td.String(), td.String(), td.String(), td.String(), td.String()
	developerID, developerIGDBID, developerName := td.Int31(), td.Int64(), td.String()
	publisherID, publisherIGDBID, publisherName := td.Int31(), td.Int64(), td.String()
	genreID, genreIGDBID, genreName := td.Int31(), td.Int64(), td.String()

	igdbGame := igdbapi.TopRatedGames{
		ID:               td.Int64(),
		Name:             td.String(),
		TotalRating:      td.Float64n(100),
		TotalRatingCount: td.Int32(),
		Cover: igdbapi.URL{
			URL: fmt.Sprintf("https://%s.com/cover.jpg", td.String()),
		},
		FirstReleaseDate: time.Now().Add(-time.Minute).Unix(),
		Genres: []igdbapi.IDName{{
			ID:   genreIGDBID,
			Name: genreName,
		}},
		InvolvedCompanies: []igdbapi.Company{{
			Company: igdbapi.IDName{
				ID:   developerIGDBID,
				Name: developerName,
			},
			Developer: true,
		}, {
			Company: igdbapi.IDName{
				ID:   publisherIGDBID,
				Name: publisherName,
			},
			Publisher: true,
		}},
		Platforms: []int64{platforms[0].IGDBID},
		Screenshots: []igdbapi.URL{
			{URL: fmt.Sprintf("https://%s.com/screenshot.png", td.String())},
		},
		Slug:    td.String(),
		Summary: td.String(),
		Websites: []igdbapi.Website{
			{URL: td.String(), Type: igdbapi.WebsiteTypeSteam},
			{URL: td.String(), Type: int8(-1)},
		},
	}

	s.storageMock.EXPECT().RunWithTx(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
		return f(ctx)
	})
	s.storageMock.EXPECT().GetTask(gomock.Any(), task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetPlatforms(gomock.Any()).Return(platforms, nil)
	s.storageMock.EXPECT().GetCompanies(gomock.Any()).Return(nil, nil)
	s.storageMock.EXPECT().GetGenres(gomock.Any()).Return(nil, nil)

	s.igdbClientMock.EXPECT().GetTopRatedGames(gomock.Any(), []int64{platforms[0].IGDBID}, gomock.Cond(func(x time.Time) bool { return x.Sub(lastReleasedAt) < time.Second }), gomock.Any(), int64(50), gomock.Any()).
		Return([]igdbapi.TopRatedGames{igdbGame}, nil).Times(1)
	// next iterations - return no games in order to stop
	s.igdbClientMock.EXPECT().GetTopRatedGames(gomock.Any(), []int64{platforms[0].IGDBID}, time.Unix(igdbGame.FirstReleaseDate, 0), gomock.Any(), int64(50), gomock.Any()).
		Return(nil, nil).Times(4)
	s.storageMock.EXPECT().GetGameIDByIGDBID(gomock.Any(), igdbGame.ID).Return(int32(0), apperr.NewNotFoundError("game", igdbGame.ID))
	s.storageMock.EXPECT().CreateCompany(gomock.Any(), model.Company{Name: developerName, IGDBID: sql.NullInt64{Valid: true, Int64: developerIGDBID}}).
		Return(developerID, nil)
	s.storageMock.EXPECT().CreateCompany(gomock.Any(), model.Company{Name: publisherName, IGDBID: sql.NullInt64{Valid: true, Int64: publisherIGDBID}}).Return(publisherID, nil)
	s.storageMock.EXPECT().CreateGenre(gomock.Any(), model.Genre{Name: genreName, IGDBID: genreIGDBID}).Return(genreID, nil)
	s.igdbClientMock.EXPECT().GetImageByURL(gomock.Any(), igdbGame.Cover.URL, igdbapi.ImageTypeCoverBig2xAlias).Return(
		igdbapi.GetImageResp{FileName: logoFileName, ContentType: contentType}, nil)
	s.s3ClientMock.EXPECT().Upload(gomock.Any(), gomock.Any(), contentType, map[string]string{"fileName": logoFileName, "game": igdbGame.Name}).
		Return(s3.UploadResult{FileURL: logoURL}, nil)
	s.igdbClientMock.EXPECT().GetImageByURL(gomock.Any(), igdbGame.Screenshots[0].URL, igdbapi.ImageTypeScreenshotBigAlias).Return(
		igdbapi.GetImageResp{FileName: screenshotFileName, ContentType: contentType}, nil)
	s.s3ClientMock.EXPECT().Upload(gomock.Any(), gomock.Any(), contentType, map[string]string{"fileName": screenshotFileName, "game": igdbGame.Name}).
		Return(s3.UploadResult{FileURL: screenshotURL}, nil)
	s.storageMock.EXPECT().CreateGame(gomock.Any(), model.CreateGameData{
		Name:             igdbGame.Name,
		DevelopersIDs:    []int32{developerID},
		PublishersIDs:    []int32{publisherID},
		ReleaseDate:      time.Unix(igdbGame.FirstReleaseDate, 0).Format("2006-01-02"),
		GenresIDs:        []int32{genreID},
		LogoURL:          logoURL,
		Summary:          igdbGame.Summary,
		Slug:             igdbGame.Slug,
		PlatformsIDs:     []int32{platforms[0].ID},
		Screenshots:      []string{screenshotURL},
		Websites:         []string{igdbGame.Websites[0].URL},
		IGDBRating:       igdbGame.TotalRating,
		IGDBRatingCount:  igdbGame.TotalRatingCount,
		IGDBID:           igdbGame.ID,
		ModerationStatus: model.ModerationStatusReady,
	}).Return(int32(1), nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil)

	err := s.provider.StartFetchIGDBGames()

	s.Require().NoError(err)
}
