// Code generated by MockGen. DO NOT EDIT.
// Source: source.go

// Package mock_sources is a generated GoMock package.
package mock_sources

import (
	models "Effective-Mobile-Music-Library/internal/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSource is a mock of Source interface.
type MockSource struct {
	ctrl     *gomock.Controller
	recorder *MockSourceMockRecorder
}

// MockSourceMockRecorder is the mock recorder for MockSource.
type MockSourceMockRecorder struct {
	mock *MockSource
}

// NewMockSource creates a new mock instance.
func NewMockSource(ctrl *gomock.Controller) *MockSource {
	mock := &MockSource{ctrl: ctrl}
	mock.recorder = &MockSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSource) EXPECT() *MockSourceMockRecorder {
	return m.recorder
}

// SongWithDetails mocks base method.
func (m *MockSource) SongWithDetails(ctx context.Context, song *models.Song) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SongWithDetails", ctx, song)
	ret0, _ := ret[0].(error)
	return ret0
}

// SongWithDetails indicates an expected call of SongWithDetails.
func (mr *MockSourceMockRecorder) SongWithDetails(ctx, song interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SongWithDetails", reflect.TypeOf((*MockSource)(nil).SongWithDetails), ctx, song)
}