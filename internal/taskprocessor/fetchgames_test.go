package taskprocessor_test

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/OutOfStack/game-library/internal/client/s3"
	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"go.uber.org/mock/gomock"
)

func (s *TestSuite) TestStartFetchIGDBGames_Success() {
	lastReleasedAt := time.Now()

	task := model.Task{
		Name:     "fetch_igdb_games",
		Status:   model.IdleTaskStatus,
		RunCount: 0,
		Settings: fmt.Appendf(nil, `{"lastReleasedAt":"%s"}`, lastReleasedAt.Format(time.RFC3339)),
	}

	platforms := []model.Platform{
		{ID: td.Int31(), IGDBID: td.Int64()},
	}

	logoURL, logoFileName, screenshotURL, screenshotFileName, contentType := td.String(), td.String(), td.String(), td.String(), td.String()
	developerID, developerIGDBID, developerName := td.Int31(), td.Int64(), td.String()
	publisherID, publisherIGDBID, publisherName := td.Int31(), td.Int64(), td.String()
	genreID, genreIGDBID, genreName := td.Int31(), td.Int64(), td.String()
	relatedIGDBID, relatedGameID := td.Int64(), td.Int31()
	igdbGameID := td.Int64()

	igdbGame := igdbapi.TopRatedGames{
		ID:               igdbGameID,
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
		Videos:       []igdbapi.Video{{VideoID: "trailer"}},
		GameModes:    []int32{1, 3},
		SimilarGames: []int64{igdbGameID, relatedIGDBID},
		Franchise:    igdbapi.Franchise{Games: []int64{relatedIGDBID}},
	}

	s.storageMock.EXPECT().RunWithTx(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error {
		return f(ctx)
	})
	s.storageMock.EXPECT().GetTask(gomock.Any(), task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil)

	s.storageMock.EXPECT().GetPlatforms(gomock.Any()).Return(platforms, nil)
	s.storageMock.EXPECT().GetCompanies(gomock.Any()).Return(nil, nil)
	s.storageMock.EXPECT().GetGenres(gomock.Any()).Return(nil, nil)
	s.storageMock.EXPECT().GetGamesIGDBIDMap(gomock.Any()).Return(map[int64]int32{relatedIGDBID: relatedGameID}, nil)

	s.igdbClientMock.EXPECT().GetTopRatedGames(gomock.Any(), []int64{platforms[0].IGDBID}, gomock.Cond(func(x time.Time) bool { return x.Sub(lastReleasedAt) < time.Second }), gomock.Any(), int64(50), gomock.Any()).
		Return([]igdbapi.TopRatedGames{igdbGame}, nil).Times(1)
	// next iterations - return no games in order to stop
	s.igdbClientMock.EXPECT().GetTopRatedGames(gomock.Any(), []int64{platforms[0].IGDBID}, time.Unix(igdbGame.FirstReleaseDate, 0), gomock.Any(), int64(50), gomock.Any()).
		Return(nil, nil).Times(4)
	s.gameFacadeMock.EXPECT().CreateCompany(gomock.Any(), model.Company{Name: developerName, IGDBID: sql.NullInt64{Valid: true, Int64: developerIGDBID}}).
		Return(developerID, nil)
	s.gameFacadeMock.EXPECT().CreateCompany(gomock.Any(), model.Company{Name: publisherName, IGDBID: sql.NullInt64{Valid: true, Int64: publisherIGDBID}}).Return(publisherID, nil)
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
		Trailers:         []string{"https://www.youtube.com/watch?v=trailer"},
		GameModeIDs:      []int32{1, 3},
		SimilarGameIDs:   []int32{relatedGameID},
		FranchiseGameIDs: []int32{relatedGameID},
		IGDBRating:       igdbGame.TotalRating,
		IGDBRatingCount:  igdbGame.TotalRatingCount,
		IGDBID:           igdbGame.ID,
		ModerationStatus: model.ModerationStatusReady,
	}).Return(int32(1), nil)

	s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil)

	err := s.provider.StartFetchIGDBGames()

	s.Require().NoError(err)
}

func (s *TestSuite) TestStartFetchIGDBGames_ResolvesGamesInsertedEarlierInBatch() {
	lastReleasedAt := time.Now()
	task := model.Task{
		Name:     "fetch_igdb_games",
		Status:   model.IdleTaskStatus,
		Settings: fmt.Appendf(nil, `{"lastReleasedAt":"%s"}`, lastReleasedAt.Format(time.RFC3339)),
	}
	first := igdbapi.TopRatedGames{
		ID:               101,
		Name:             "First",
		Cover:            igdbapi.URL{URL: "https://example.com/first.jpg"},
		FirstReleaseDate: lastReleasedAt.Add(-2 * time.Minute).Unix(),
	}
	second := igdbapi.TopRatedGames{
		ID:               102,
		Name:             "Second",
		Cover:            igdbapi.URL{URL: "https://example.com/second.jpg"},
		FirstReleaseDate: lastReleasedAt.Add(-time.Minute).Unix(),
		SimilarGames:     []int64{first.ID},
	}

	s.storageMock.EXPECT().RunWithTx(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func(context.Context) error) error { return f(ctx) })
	s.storageMock.EXPECT().GetTask(gomock.Any(), task.Name).Return(task, nil)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil)
	s.storageMock.EXPECT().GetGamesIGDBIDMap(gomock.Any()).Return(map[int64]int32{}, nil)
	s.storageMock.EXPECT().GetPlatforms(gomock.Any()).Return(nil, nil)
	s.storageMock.EXPECT().GetCompanies(gomock.Any()).Return(nil, nil)
	s.storageMock.EXPECT().GetGenres(gomock.Any()).Return(nil, nil)
	s.igdbClientMock.EXPECT().GetTopRatedGames(gomock.Any(), nil, gomock.Any(), gomock.Any(), int64(50), gomock.Any()).Return([]igdbapi.TopRatedGames{first, second}, nil)
	s.igdbClientMock.EXPECT().GetTopRatedGames(gomock.Any(), nil, gomock.Any(), gomock.Any(), int64(50), gomock.Any()).Return(nil, nil).Times(4)
	s.igdbClientMock.EXPECT().GetImageByURL(gomock.Any(), first.Cover.URL, igdbapi.ImageTypeCoverBig2xAlias).Return(igdbapi.GetImageResp{FileName: "first.jpg", ContentType: "image/jpeg"}, nil)
	s.s3ClientMock.EXPECT().Upload(gomock.Any(), gomock.Any(), "image/jpeg", gomock.Any()).Return(s3.UploadResult{FileURL: "first.jpg"}, nil)
	s.igdbClientMock.EXPECT().GetImageByURL(gomock.Any(), second.Cover.URL, igdbapi.ImageTypeCoverBig2xAlias).Return(igdbapi.GetImageResp{FileName: "second.jpg", ContentType: "image/jpeg"}, nil)
	s.s3ClientMock.EXPECT().Upload(gomock.Any(), gomock.Any(), "image/jpeg", gomock.Any()).Return(s3.UploadResult{FileURL: "second.jpg"}, nil)

	var created []model.CreateGameData
	s.storageMock.EXPECT().CreateGame(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, game model.CreateGameData) (int32, error) {
		created = append(created, game)
		return int32(len(created)), nil
	}).Times(2)
	s.storageMock.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).Return(nil)

	err := s.provider.StartFetchIGDBGames()

	s.Require().NoError(err)
	s.Require().Len(created, 2)
	s.Equal([]int32{1}, created[1].SimilarGameIDs)
}
