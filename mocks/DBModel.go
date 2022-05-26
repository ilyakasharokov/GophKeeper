// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "gophkeeper/internal/common/models"

	time "time"
)

// DBModel is an autogenerated mock type for the DBModel type
type DBModel struct {
	mock.Mock
}

// AddNote provides a mock function with given fields: ctx, userID, note
func (_m *DBModel) AddNote(ctx context.Context, userID string, note models.Note) (string, error) {
	ret := _m.Called(ctx, userID, note)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, models.Note) string); ok {
		r0 = rf(ctx, userID, note)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, models.Note) error); ok {
		r1 = rf(ctx, userID, note)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckUserPassword provides a mock function with given fields: ctx, user
func (_m *DBModel) CheckUserPassword(ctx context.Context, user models.User) (string, error) {
	ret := _m.Called(ctx, user)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, models.User) string); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(string)
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
func (_m *DBModel) CreateUser(ctx context.Context, user models.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteNote provides a mock function with given fields: ctx, userID, id
func (_m *DBModel) DeleteNote(ctx context.Context, userID string, id string) error {
	ret := _m.Called(ctx, userID, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, userID, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetNotes provides a mock function with given fields: ctx, userID
func (_m *DBModel) GetNotes(ctx context.Context, userID string) ([]models.Note, error) {
	ret := _m.Called(ctx, userID)

	var r0 []models.Note
	if rf, ok := ret.Get(0).(func(context.Context, string) []models.Note); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Note)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUpdates provides a mock function with given fields: ctx, userID, after
func (_m *DBModel) GetUpdates(ctx context.Context, userID string, after time.Time) ([]models.Note, error) {
	ret := _m.Called(ctx, userID, after)

	var r0 []models.Note
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) []models.Note); ok {
		r0 = rf(ctx, userID, after)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Note)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time) error); ok {
		r1 = rf(ctx, userID, after)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateNote provides a mock function with given fields: ctx, userID, note
func (_m *DBModel) UpdateNote(ctx context.Context, userID string, note models.Note) error {
	ret := _m.Called(ctx, userID, note)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, models.Note) error); ok {
		r0 = rf(ctx, userID, note)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
