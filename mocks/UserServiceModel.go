// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"
	authorization "gophkeeper/internal/app/server/authorization"

	mock "github.com/stretchr/testify/mock"

	models "gophkeeper/internal/common/models"
)

// UserServiceModel is an autogenerated mock type for the UserServiceModel type
type UserServiceModel struct {
	mock.Mock
}

// AuthUser provides a mock function with given fields: ctx, user
func (_m *UserServiceModel) AuthUser(ctx context.Context, user models.User) (*authorization.TokenInfo, error) {
	ret := _m.Called(ctx, user)

	var r0 *authorization.TokenInfo
	if rf, ok := ret.Get(0).(func(context.Context, models.User) *authorization.TokenInfo); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*authorization.TokenInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *UserServiceModel) CreateUser(ctx context.Context, user models.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RefreshToken provides a mock function with given fields: ctx, token
func (_m *UserServiceModel) RefreshToken(ctx context.Context, token string) (*authorization.TokenInfo, error) {
	ret := _m.Called(ctx, token)

	var r0 *authorization.TokenInfo
	if rf, ok := ret.Get(0).(func(context.Context, string) *authorization.TokenInfo); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*authorization.TokenInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
