package repo_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
)

// TestCreateCompany_IGDBIDNull_ShouldBeNoError tests case when we add company without igdb id, and there should be no error
func TestCreateCompany_IGDBIDIsNull_ShouldBeNoError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	company := repo.Company{
		Name: td.String(),
	}

	_, err := s.CreateCompany(context.Background(), company)
	require.NoError(t, err, "err should be nil")
}

// TestGetCompanies_DataExists_ShouldBeEqual tests case when we add one company, then fetch first company, and they should be equal
func TestGetCompanies_DataExists_ShouldBeEqual(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	company := repo.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	}

	id, err := s.CreateCompany(context.Background(), company)
	require.NoError(t, err, "err should be nil")

	companies, err := s.GetCompanies(context.Background())
	require.NoError(t, err, "err should be nil")
	require.Equal(t, len(companies), 1, "companies len should be 1")

	want := company
	got := companies[0]
	require.Equal(t, id, got.ID, "id should be equal")
	require.Equal(t, want.Name, got.Name, "name should be equal")
	require.Equal(t, want.IGDBID, got.IGDBID, "igdb id should be equal")
}

// GetCompanyIDByName_CompanyExists_ShouldReturnID tests case when we add one company, then get id by name
func TestGetCompanyIDByName_CompanyExists_ShouldReturnID(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	company := repo.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	}

	id, err := s.CreateCompany(context.Background(), company)
	require.NoError(t, err, "err should be nil")

	gotID, err := s.GetCompanyIDByName(context.Background(), strings.ToUpper(company.Name))
	require.NoError(t, err, "err should be nil")
	require.Equal(t, id, gotID, "id should be equal")
}

// GetCompanyIDByName_CompanyNotExist_ShouldReturnErrNotFound tests case when we add one company, then get id by another name
func TestGetCompanyIDByName_CompanyNotExist_ShouldReturnErrNotFound(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	company := repo.Company{
		Name:   td.String(),
		IGDBID: sql.NullInt64{Int64: td.Int64(), Valid: true},
	}

	_, err := s.CreateCompany(context.Background(), company)
	require.NoError(t, err, "err should be nil")

	randomName := td.String()
	gotID, err := s.GetCompanyIDByName(context.Background(), randomName)
	require.ErrorIs(t, err, repo.ErrNotFound[string]{Entity: "company", ID: randomName}, "err should be NotFound")
	require.Zero(t, gotID, "got id should be 0")
}
