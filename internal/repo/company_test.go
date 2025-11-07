package repo_test

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
)

// TestCreateCompany_IGDBIDNull_ShouldBeNoError tests case when we add company without igdb id, and there should be no error
func TestCreateCompany_IGDBIDIsNull_ShouldBeNoError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	company := model.Company{
		Name: td.String(),
	}

	_, err := s.CreateCompany(t.Context(), company)
	require.NoError(t, err)
}

// TestGetCompanies_DataExists_ShouldBeEqual tests case when we add one company, then fetch first company, and they should be equal
func TestGetCompanies_DataExists_ShouldBeEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	company := model.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	}

	id, err := s.CreateCompany(ctx, company)
	require.NoError(t, err)

	companies, err := s.GetCompanies(ctx)
	require.NoError(t, err)
	require.Len(t, companies, 1, "companies len should be 1")

	want := company
	got := companies[0]
	require.Equal(t, id, got.ID, "id should be equal")
	require.Equal(t, want.Name, got.Name, "name should be equal")
	require.Equal(t, want.IGDBID, got.IGDBID, "igdb id should be equal")
}

// TestGetCompanyIDByName_CompanyExists_ShouldReturnID tests case when we add one company, then get id by name
func TestGetCompanyIDByName_CompanyExists_ShouldReturnID(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	company := model.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	}

	id, err := s.CreateCompany(ctx, company)
	require.NoError(t, err)

	gotID, err := s.GetCompanyIDByName(ctx, strings.ToUpper(company.Name))
	require.NoError(t, err)
	require.Equal(t, id, gotID, "id should be equal")
}

// TestGetCompanyIDByName_CompanyNotExist_ShouldReturnNotFoundError tests case when we add one company, then get id by another name
func TestGetCompanyIDByName_CompanyNotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	company := model.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	}
	randomName := td.String()

	_, err := s.CreateCompany(ctx, company)
	require.NoError(t, err)

	gotID, err := s.GetCompanyIDByName(ctx, randomName)
	require.ErrorIs(t, err, apperr.NewNotFoundError("company", randomName), "err should be NotFound")
	require.Zero(t, gotID, "got id should be 0")
}

func TestGetCompanyByID_CompanyExists_ShouldReturnCompany(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	company := model.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	}

	id, err := s.CreateCompany(ctx, company)
	require.NoError(t, err)

	gotCompany, err := s.GetCompanyByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, gotCompany.ID, "id should be equal")
	require.Equal(t, company.Name, gotCompany.Name, "name should be equal")
	require.Equal(t, company.IGDBID.Int64, gotCompany.IGDBID.Int64, "igdb id should be equal")
}

func TestGetCompanyByID_CompanyNotExist_ShouldReturnNotFoundError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	company := model.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	}

	_, err := s.CreateCompany(ctx, company)
	require.NoError(t, err)

	randomID := td.Int32()
	gotCompany, err := s.GetCompanyByID(ctx, randomID)
	require.ErrorIs(t, err, apperr.NewNotFoundError("company", randomID), "err should be NotFound")
	require.Zero(t, gotCompany.ID, "got id should be 0")
}

func TestGetTopDevelopers_Ok(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	// create 3 developers and 4 games
	developer1ID, err := s.CreateCompany(ctx, model.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	})
	require.NoError(t, err)

	developer2ID, err := s.CreateCompany(ctx, model.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	})
	require.NoError(t, err)

	developer3ID, err := s.CreateCompany(ctx, model.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	})
	require.NoError(t, err)

	cg1, cg2, cg3, cg4 := getCreateGameData(), getCreateGameData(), getCreateGameData(), getCreateGameData()

	// developer 1 developed 1 game, developer 2 developed 3 games, developer 3 developed 2 games
	cg1.DevelopersIDs = []int32{developer1ID}
	cg2.DevelopersIDs = []int32{developer2ID}
	cg3.DevelopersIDs = []int32{developer2ID, developer3ID}
	cg4.DevelopersIDs = []int32{developer2ID, developer3ID}

	_, err = s.CreateGame(ctx, cg1)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg2)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg3)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg4)
	require.NoError(t, err)

	top, err := s.GetTopDevelopers(ctx, 5)
	require.NoError(t, err)

	require.Len(t, top, 3, "len of top developers should be 3")

	require.Equal(t, developer2ID, top[0].ID, "top 1 developer should be developer 2")
	require.Equal(t, developer3ID, top[1].ID, "top 2 developer should be developer 3")
	require.Equal(t, developer1ID, top[2].ID, "top 3 developer should be developer 1")
}

func TestGetTopPublishers_Ok(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	// create 2 publishers and 4 games
	publisher1ID, err := s.CreateCompany(ctx, model.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	})
	require.NoError(t, err)

	publisher2ID, err := s.CreateCompany(ctx, model.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	})
	require.NoError(t, err)

	cg1, cg2, cg3, cg4 := getCreateGameData(), getCreateGameData(), getCreateGameData(), getCreateGameData()

	// publisher 1 published 2 games, publisher 2 published 3 games
	cg1.PublishersIDs = []int32{publisher1ID}
	cg2.PublishersIDs = []int32{publisher2ID}
	cg3.PublishersIDs = []int32{publisher2ID}
	cg4.PublishersIDs = []int32{publisher1ID, publisher2ID}

	_, err = s.CreateGame(ctx, cg1)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg2)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg3)
	require.NoError(t, err)

	_, err = s.CreateGame(ctx, cg4)
	require.NoError(t, err)

	top, err := s.GetTopPublishers(ctx, 5)
	require.NoError(t, err)

	require.Len(t, top, 2, "len of top publishers should be 2")

	require.Equal(t, publisher2ID, top[0].ID, "top 1 publisher should be publisher 2")
	require.Equal(t, publisher1ID, top[1].ID, "top 2 publisher should be publisher 1")
}
