// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	models "gophkeeper/pkg/models"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// StorageModel is an autogenerated mock type for the StorageModel type
type StorageModel struct {
	mock.Mock
}

// GetLastSyncDate provides a mock function with given fields:
func (_m *StorageModel) GetLastSyncDate() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// GetNotes provides a mock function with given fields: all
func (_m *StorageModel) GetNotes(all bool) []models.Note {
	ret := _m.Called(all)

	var r0 []models.Note
	if rf, ok := ret.Get(0).(func(bool) []models.Note); ok {
		r0 = rf(all)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Note)
		}
	}

	return r0
}

// SetLastSyncDate provides a mock function with given fields: date
func (_m *StorageModel) SetLastSyncDate(date time.Time) {
	_m.Called(date)
}

// SetNotes provides a mock function with given fields: notes
func (_m *StorageModel) SetNotes(notes []models.Note) {
	_m.Called(notes)
}
