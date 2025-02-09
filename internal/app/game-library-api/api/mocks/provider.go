// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/game-library-api/api/provider.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/game-library-api/api/provider.go -destination=internal/app/game-library-api/api/mocks/provider.go -package=api_mock
//

// Package api_mock is a generated GoMock package.
package api_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	gomock "go.uber.org/mock/gomock"
)

// MockGameFacade is a mock of GameFacade interface.
type MockGameFacade struct {
	ctrl     *gomock.Controller
	recorder *MockGameFacadeMockRecorder
	isgomock struct{}
}

// MockGameFacadeMockRecorder is the mock recorder for MockGameFacade.
type MockGameFacadeMockRecorder struct {
	mock *MockGameFacade
}

// NewMockGameFacade creates a new mock instance.
func NewMockGameFacade(ctrl *gomock.Controller) *MockGameFacade {
	mock := &MockGameFacade{ctrl: ctrl}
	mock.recorder = &MockGameFacadeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGameFacade) EXPECT() *MockGameFacadeMockRecorder {
	return m.recorder
}

// CreateGame mocks base method.
func (m *MockGameFacade) CreateGame(ctx context.Context, cg model.CreateGame) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGame", ctx, cg)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGame indicates an expected call of CreateGame.
func (mr *MockGameFacadeMockRecorder) CreateGame(ctx, cg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGame", reflect.TypeOf((*MockGameFacade)(nil).CreateGame), ctx, cg)
}

// DeleteGame mocks base method.
func (m *MockGameFacade) DeleteGame(ctx context.Context, id int32, publisher string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGame", ctx, id, publisher)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGame indicates an expected call of DeleteGame.
func (mr *MockGameFacadeMockRecorder) DeleteGame(ctx, id, publisher any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGame", reflect.TypeOf((*MockGameFacade)(nil).DeleteGame), ctx, id, publisher)
}

// GetCompaniesMap mocks base method.
func (m *MockGameFacade) GetCompaniesMap(ctx context.Context) (map[int32]model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompaniesMap", ctx)
	ret0, _ := ret[0].(map[int32]model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompaniesMap indicates an expected call of GetCompaniesMap.
func (mr *MockGameFacadeMockRecorder) GetCompaniesMap(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompaniesMap", reflect.TypeOf((*MockGameFacade)(nil).GetCompaniesMap), ctx)
}

// GetGameByID mocks base method.
func (m *MockGameFacade) GetGameByID(ctx context.Context, id int32) (model.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGameByID", ctx, id)
	ret0, _ := ret[0].(model.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGameByID indicates an expected call of GetGameByID.
func (mr *MockGameFacadeMockRecorder) GetGameByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGameByID", reflect.TypeOf((*MockGameFacade)(nil).GetGameByID), ctx, id)
}

// GetGames mocks base method.
func (m *MockGameFacade) GetGames(ctx context.Context, page, pageSize int, filter model.GamesFilter) ([]model.Game, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGames", ctx, page, pageSize, filter)
	ret0, _ := ret[0].([]model.Game)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetGames indicates an expected call of GetGames.
func (mr *MockGameFacadeMockRecorder) GetGames(ctx, page, pageSize, filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGames", reflect.TypeOf((*MockGameFacade)(nil).GetGames), ctx, page, pageSize, filter)
}

// GetGenres mocks base method.
func (m *MockGameFacade) GetGenres(ctx context.Context) ([]model.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenres", ctx)
	ret0, _ := ret[0].([]model.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenres indicates an expected call of GetGenres.
func (mr *MockGameFacadeMockRecorder) GetGenres(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenres", reflect.TypeOf((*MockGameFacade)(nil).GetGenres), ctx)
}

// GetGenresMap mocks base method.
func (m *MockGameFacade) GetGenresMap(ctx context.Context) (map[int32]model.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenresMap", ctx)
	ret0, _ := ret[0].(map[int32]model.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenresMap indicates an expected call of GetGenresMap.
func (mr *MockGameFacadeMockRecorder) GetGenresMap(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenresMap", reflect.TypeOf((*MockGameFacade)(nil).GetGenresMap), ctx)
}

// GetPlatforms mocks base method.
func (m *MockGameFacade) GetPlatforms(ctx context.Context) ([]model.Platform, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlatforms", ctx)
	ret0, _ := ret[0].([]model.Platform)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlatforms indicates an expected call of GetPlatforms.
func (mr *MockGameFacadeMockRecorder) GetPlatforms(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlatforms", reflect.TypeOf((*MockGameFacade)(nil).GetPlatforms), ctx)
}

// GetPlatformsMap mocks base method.
func (m *MockGameFacade) GetPlatformsMap(ctx context.Context) (map[int32]model.Platform, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlatformsMap", ctx)
	ret0, _ := ret[0].(map[int32]model.Platform)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlatformsMap indicates an expected call of GetPlatformsMap.
func (mr *MockGameFacadeMockRecorder) GetPlatformsMap(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlatformsMap", reflect.TypeOf((*MockGameFacade)(nil).GetPlatformsMap), ctx)
}

// GetTopCompanies mocks base method.
func (m *MockGameFacade) GetTopCompanies(ctx context.Context, companyType string, limit int64) ([]model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopCompanies", ctx, companyType, limit)
	ret0, _ := ret[0].([]model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopCompanies indicates an expected call of GetTopCompanies.
func (mr *MockGameFacadeMockRecorder) GetTopCompanies(ctx, companyType, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopCompanies", reflect.TypeOf((*MockGameFacade)(nil).GetTopCompanies), ctx, companyType, limit)
}

// GetTopGenres mocks base method.
func (m *MockGameFacade) GetTopGenres(ctx context.Context, limit int64) ([]model.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopGenres", ctx, limit)
	ret0, _ := ret[0].([]model.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopGenres indicates an expected call of GetTopGenres.
func (mr *MockGameFacadeMockRecorder) GetTopGenres(ctx, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopGenres", reflect.TypeOf((*MockGameFacade)(nil).GetTopGenres), ctx, limit)
}

// GetUserRatings mocks base method.
func (m *MockGameFacade) GetUserRatings(ctx context.Context, userID string) (map[int32]uint8, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserRatings", ctx, userID)
	ret0, _ := ret[0].(map[int32]uint8)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserRatings indicates an expected call of GetUserRatings.
func (mr *MockGameFacadeMockRecorder) GetUserRatings(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserRatings", reflect.TypeOf((*MockGameFacade)(nil).GetUserRatings), ctx, userID)
}

// RateGame mocks base method.
func (m *MockGameFacade) RateGame(ctx context.Context, gameID int32, userID string, rating uint8) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RateGame", ctx, gameID, userID, rating)
	ret0, _ := ret[0].(error)
	return ret0
}

// RateGame indicates an expected call of RateGame.
func (mr *MockGameFacadeMockRecorder) RateGame(ctx, gameID, userID, rating any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RateGame", reflect.TypeOf((*MockGameFacade)(nil).RateGame), ctx, gameID, userID, rating)
}

// UpdateGame mocks base method.
func (m *MockGameFacade) UpdateGame(ctx context.Context, id int32, publisher string, upd model.UpdatedGame) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGame", ctx, id, publisher, upd)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGame indicates an expected call of UpdateGame.
func (mr *MockGameFacadeMockRecorder) UpdateGame(ctx, id, publisher, upd any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGame", reflect.TypeOf((*MockGameFacade)(nil).UpdateGame), ctx, id, publisher, upd)
}
