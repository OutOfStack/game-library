package taskprocessor_test

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/client/igdb"
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

	logoURL, screenshotURL := td.String(), td.String()
	developerID, developerIGDBID, developerName := td.Int31(), td.Int64(), td.String()
	publisherID, publisherIGDBID, publisherName := td.Int31(), td.Int64(), td.String()
	genreID, genreIGDBID, genreName := td.Int31(), td.Int64(), td.String()

	igdbGame := igdb.TopRatedGamesResp{
		ID:               td.Int64(),
		Name:             td.String(),
		TotalRating:      td.Float64(),
		TotalRatingCount: td.Int64(),
		Cover: igdb.URL{
			URL: fmt.Sprintf("https://%s.com/cover.jpg", td.String()),
		},
		FirstReleaseDate: time.Now().Add(-time.Minute).Unix(),
		Genres: []igdb.IDName{{
			ID:   genreIGDBID,
			Name: genreName,
		}},
		InvolvedCompanies: []igdb.Company{{
			Company: igdb.IDName{
				ID:   developerIGDBID,
				Name: developerName,
			},
			Developer: true,
		}, {
			Company: igdb.IDName{
				ID:   publisherIGDBID,
				Name: publisherName,
			},
			Publisher: true,
		}},
		Platforms: []int64{platforms[0].IGDBID},
		Screenshots: []igdb.URL{
			{URL: fmt.Sprintf("https://%s.com/screenshot.png", td.String())},
		},
		Slug:    td.String(),
		Summary: td.String(),
		Websites: []igdb.Website{
			{URL: td.String(), Category: int8(igdb.WebsiteCategorySteam)},
			{URL: td.String(), Category: int8(-1)},
		},
	}

	s.storageMock.EXPECT().BeginTx(gomock.Any()).Return(s.tx, nil)
	s.storageMock.EXPECT().GetTask(gomock.Any(), s.tx, task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), s.tx, gomock.Any()).Return(nil)
	s.tx.EXPECT().Commit(gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetPlatforms(gomock.Any()).Return(platforms, nil)
	s.storageMock.EXPECT().GetCompanies(gomock.Any()).Return(nil, nil)
	s.storageMock.EXPECT().GetGenres(gomock.Any()).Return(nil, nil)

	s.igdbClientMock.EXPECT().GetTopRatedGames(gomock.Any(), []int64{platforms[0].IGDBID}, gomock.Cond(func(x time.Time) bool { return x.Sub(lastReleasedAt) < time.Second }), gomock.Any(), int64(60), gomock.Any()).
		Return([]igdb.TopRatedGamesResp{igdbGame}, nil).Times(1)
	// next iterations - return no games in order to stop
	s.igdbClientMock.EXPECT().GetTopRatedGames(gomock.Any(), []int64{platforms[0].IGDBID}, time.Unix(igdbGame.FirstReleaseDate, 0), gomock.Any(), int64(60), gomock.Any()).
		Return(nil, nil).Times(4)
	s.storageMock.EXPECT().GetGameIDByIGDBID(gomock.Any(), igdbGame.ID).Return(int32(0), apperr.NewNotFoundError("game", igdbGame.ID))
	s.storageMock.EXPECT().CreateCompany(gomock.Any(), model.Company{Name: developerName, IGDBID: sql.NullInt64{Valid: true, Int64: developerIGDBID}}).
		Return(developerID, nil)
	s.storageMock.EXPECT().CreateCompany(gomock.Any(), model.Company{Name: publisherName, IGDBID: sql.NullInt64{Valid: true, Int64: publisherIGDBID}}).Return(publisherID, nil)
	s.storageMock.EXPECT().CreateGenre(gomock.Any(), model.Genre{Name: genreName, IGDBID: genreIGDBID}).Return(genreID, nil)
	s.uploadcareMock.EXPECT().UploadImageFromURL(gomock.Any(), igdbGame.Cover.URL).Return(logoURL, nil)
	s.uploadcareMock.EXPECT().UploadImageFromURL(gomock.Any(), igdbGame.Screenshots[0].URL).Return(screenshotURL, nil)
	s.storageMock.EXPECT().CreateGame(gomock.Any(), model.CreateGame{
		Name:          igdbGame.Name,
		DevelopersIDs: []int32{developerID},
		PublishersIDs: []int32{publisherID},
		ReleaseDate:   time.Unix(igdbGame.FirstReleaseDate, 0).Format("2006-01-02"),
		Genres:        []int32{genreID},
		LogoURL:       logoURL,
		Summary:       igdbGame.Summary,
		Slug:          igdbGame.Slug,
		Platforms:     []int32{platforms[0].ID},
		Screenshots:   []string{screenshotURL},
		Websites:      []string{igdbGame.Websites[0].URL},
		IGDBRating:    igdbGame.TotalRating,
		IGDBID:        igdbGame.ID,
	}).Return(int32(1), nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), nil, gomock.Any()).Return(nil)

	err := s.provider.StartFetchIGDBGames()

	s.Require().NoError(err)
}
