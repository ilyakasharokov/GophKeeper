// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "gophkeeper/internal/common/models"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// GRPCClientModel is an autogenerated mock type for the GRPCClientModel type
type GRPCClientModel struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *GRPCClientModel) Close() {
	_m.Called()
}

// Login provides a mock function with given fields: login, pwd
func (_m *GRPCClientModel) Login(login string, pwd string) (string, error) {
	ret := _m.Called(login, pwd)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(login, pwd)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(login, pwd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RefreshToken provides a mock function with given fields: ctx, refreshToken
func (_m *GRPCClientModel) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	ret := _m.Called(ctx, refreshToken)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, refreshToken)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string) string); ok {
		r1 = rf(ctx, refreshToken)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, refreshToken)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Registration provides a mock function with given fields: login, pwd
func (_m *GRPCClientModel) Registration(login string, pwd string) (string, string, error) {
	ret := _m.Called(login, pwd)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(login, pwd)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string, string) string); ok {
		r1 = rf(login, pwd)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(login, pwd)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SyncData provides a mock function with given fields: notes, lastSync
func (_m *GRPCClientModel) SyncData(notes []models.Note, lastSync time.Time) ([]models.Note, time.Time, error) {
	ret := _m.Called(notes, lastSync)

	var r0 []models.Note
	if rf, ok := ret.Get(0).(func([]models.Note, time.Time) []models.Note); ok {
		r0 = rf(notes, lastSync)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Note)
		}
	}

	var r1 time.Time
	if rf, ok := ret.Get(1).(func([]models.Note, time.Time) time.Time); ok {
		r1 = rf(notes, lastSync)
	} else {
		r1 = ret.Get(1).(time.Time)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func([]models.Note, time.Time) error); ok {
		r2 = rf(notes, lastSync)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}