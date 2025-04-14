// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/game-library-api/facade/provider.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/game-library-api/facade/provider.go -destination=internal/app/game-library-api/facade/mocks/provider.go -package=facade_mock
//

// Package facade_mock is a generated GoMock package.
package facade_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
	isgomock struct{}
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddRating mocks base method.
func (m *MockStorage) AddRating(ctx context.Context, cr model.CreateRating) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRating", ctx, cr)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRating indicates an expected call of AddRating.
func (mr *MockStorageMockRecorder) AddRating(ctx, cr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRating", reflect.TypeOf((*MockStorage)(nil).AddRating), ctx, cr)
}

// CreateCompany mocks base method.
func (m *MockStorage) CreateCompany(ctx context.Context, c model.Company) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCompany", ctx, c)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCompany indicates an expected call of CreateCompany.
func (mr *MockStorageMockRecorder) CreateCompany(ctx, c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCompany", reflect.TypeOf((*MockStorage)(nil).CreateCompany), ctx, c)
}

// CreateGame mocks base method.
func (m *MockStorage) CreateGame(ctx context.Context, cg model.CreateGameData) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGame", ctx, cg)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGame indicates an expected call of CreateGame.
func (mr *MockStorageMockRecorder) CreateGame(ctx, cg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGame", reflect.TypeOf((*MockStorage)(nil).CreateGame), ctx, cg)
}

// DeleteGame mocks base method.
func (m *MockStorage) DeleteGame(ctx context.Context, id int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGame", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGame indicates an expected call of DeleteGame.
func (mr *MockStorageMockRecorder) DeleteGame(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGame", reflect.TypeOf((*MockStorage)(nil).DeleteGame), ctx, id)
}

// GetCompanies mocks base method.
func (m *MockStorage) GetCompanies(ctx context.Context) ([]model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanies", ctx)
	ret0, _ := ret[0].([]model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanies indicates an expected call of GetCompanies.
func (mr *MockStorageMockRecorder) GetCompanies(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanies", reflect.TypeOf((*MockStorage)(nil).GetCompanies), ctx)
}

// GetCompanyByID mocks base method.
func (m *MockStorage) GetCompanyByID(ctx context.Context, id int32) (model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyByID", ctx, id)
	ret0, _ := ret[0].(model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyByID indicates an expected call of GetCompanyByID.
func (mr *MockStorageMockRecorder) GetCompanyByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyByID", reflect.TypeOf((*MockStorage)(nil).GetCompanyByID), ctx, id)
}

// GetCompanyIDByName mocks base method.
func (m *MockStorage) GetCompanyIDByName(ctx context.Context, name string) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyIDByName", ctx, name)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyIDByName indicates an expected call of GetCompanyIDByName.
func (mr *MockStorageMockRecorder) GetCompanyIDByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyIDByName", reflect.TypeOf((*MockStorage)(nil).GetCompanyIDByName), ctx, name)
}

// GetGameByID mocks base method.
func (m *MockStorage) GetGameByID(ctx context.Context, id int32) (model.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGameByID", ctx, id)
	ret0, _ := ret[0].(model.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGameByID indicates an expected call of GetGameByID.
func (mr *MockStorageMockRecorder) GetGameByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGameByID", reflect.TypeOf((*MockStorage)(nil).GetGameByID), ctx, id)
}

// GetGames mocks base method.
func (m *MockStorage) GetGames(ctx context.Context, pageSize, page int, filter model.GamesFilter) ([]model.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGames", ctx, pageSize, page, filter)
	ret0, _ := ret[0].([]model.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGames indicates an expected call of GetGames.
func (mr *MockStorageMockRecorder) GetGames(ctx, pageSize, page, filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGames", reflect.TypeOf((*MockStorage)(nil).GetGames), ctx, pageSize, page, filter)
}

// GetGamesCount mocks base method.
func (m *MockStorage) GetGamesCount(ctx context.Context, filter model.GamesFilter) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGamesCount", ctx, filter)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGamesCount indicates an expected call of GetGamesCount.
func (mr *MockStorageMockRecorder) GetGamesCount(ctx, filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGamesCount", reflect.TypeOf((*MockStorage)(nil).GetGamesCount), ctx, filter)
}

// GetGenreByID mocks base method.
func (m *MockStorage) GetGenreByID(ctx context.Context, id int32) (model.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenreByID", ctx, id)
	ret0, _ := ret[0].(model.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenreByID indicates an expected call of GetGenreByID.
func (mr *MockStorageMockRecorder) GetGenreByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenreByID", reflect.TypeOf((*MockStorage)(nil).GetGenreByID), ctx, id)
}

// GetGenres mocks base method.
func (m *MockStorage) GetGenres(ctx context.Context) ([]model.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenres", ctx)
	ret0, _ := ret[0].([]model.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenres indicates an expected call of GetGenres.
func (mr *MockStorageMockRecorder) GetGenres(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenres", reflect.TypeOf((*MockStorage)(nil).GetGenres), ctx)
}

// GetPlatformByID mocks base method.
func (m *MockStorage) GetPlatformByID(ctx context.Context, id int32) (model.Platform, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlatformByID", ctx, id)
	ret0, _ := ret[0].(model.Platform)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlatformByID indicates an expected call of GetPlatformByID.
func (mr *MockStorageMockRecorder) GetPlatformByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlatformByID", reflect.TypeOf((*MockStorage)(nil).GetPlatformByID), ctx, id)
}

// GetPlatforms mocks base method.
func (m *MockStorage) GetPlatforms(ctx context.Context) ([]model.Platform, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlatforms", ctx)
	ret0, _ := ret[0].([]model.Platform)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlatforms indicates an expected call of GetPlatforms.
func (mr *MockStorageMockRecorder) GetPlatforms(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlatforms", reflect.TypeOf((*MockStorage)(nil).GetPlatforms), ctx)
}

// GetTopDevelopers mocks base method.
func (m *MockStorage) GetTopDevelopers(ctx context.Context, limit int64) ([]model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopDevelopers", ctx, limit)
	ret0, _ := ret[0].([]model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopDevelopers indicates an expected call of GetTopDevelopers.
func (mr *MockStorageMockRecorder) GetTopDevelopers(ctx, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopDevelopers", reflect.TypeOf((*MockStorage)(nil).GetTopDevelopers), ctx, limit)
}

// GetTopGenres mocks base method.
func (m *MockStorage) GetTopGenres(ctx context.Context, limit int64) ([]model.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopGenres", ctx, limit)
	ret0, _ := ret[0].([]model.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopGenres indicates an expected call of GetTopGenres.
func (mr *MockStorageMockRecorder) GetTopGenres(ctx, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopGenres", reflect.TypeOf((*MockStorage)(nil).GetTopGenres), ctx, limit)
}

// GetTopPublishers mocks base method.
func (m *MockStorage) GetTopPublishers(ctx context.Context, limit int64) ([]model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopPublishers", ctx, limit)
	ret0, _ := ret[0].([]model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopPublishers indicates an expected call of GetTopPublishers.
func (mr *MockStorageMockRecorder) GetTopPublishers(ctx, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopPublishers", reflect.TypeOf((*MockStorage)(nil).GetTopPublishers), ctx, limit)
}

// GetUserRatings mocks base method.
func (m *MockStorage) GetUserRatings(ctx context.Context, userID string) (map[int32]uint8, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserRatings", ctx, userID)
	ret0, _ := ret[0].(map[int32]uint8)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserRatings indicates an expected call of GetUserRatings.
func (mr *MockStorageMockRecorder) GetUserRatings(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserRatings", reflect.TypeOf((*MockStorage)(nil).GetUserRatings), ctx, userID)
}

// RemoveRating mocks base method.
func (m *MockStorage) RemoveRating(ctx context.Context, rr model.RemoveRating) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveRating", ctx, rr)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveRating indicates an expected call of RemoveRating.
func (mr *MockStorageMockRecorder) RemoveRating(ctx, rr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRating", reflect.TypeOf((*MockStorage)(nil).RemoveRating), ctx, rr)
}

// UpdateGame mocks base method.
func (m *MockStorage) UpdateGame(ctx context.Context, id int32, ug model.UpdateGameData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGame", ctx, id, ug)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGame indicates an expected call of UpdateGame.
func (mr *MockStorageMockRecorder) UpdateGame(ctx, id, ug any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGame", reflect.TypeOf((*MockStorage)(nil).UpdateGame), ctx, id, ug)
}

// UpdateGameRating mocks base method.
func (m *MockStorage) UpdateGameRating(ctx context.Context, id int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGameRating", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGameRating indicates an expected call of UpdateGameRating.
func (mr *MockStorageMockRecorder) UpdateGameRating(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGameRating", reflect.TypeOf((*MockStorage)(nil).UpdateGameRating), ctx, id)
}
