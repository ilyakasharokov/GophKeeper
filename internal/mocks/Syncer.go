// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	models "gophkeeper/pkg/models"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Syncer is an autogenerated mock type for the Syncer type
type Syncer struct {
	mock.Mock
}

// GetLastSyncDate provides a mock function with given fields:
func (_m *Syncer) GetLastSyncDate() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// GetNonSyncNotes provides a mock function with given fields:
func (_m *Syncer) GetNonSyncNotes() []models.Note {
	ret := _m.Called()

	var r0 []models.Note
	if rf, ok := ret.Get(0).(func() []models.Note); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Note)
		}
	}

	return r0
}

// Sync provides a mock function with given fields:
func (_m *Syncer) Sync() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLastSyncDate provides a mock function with given fields:
func (_m *Syncer) UpdateLastSyncDate() {
	_m.Called()
}