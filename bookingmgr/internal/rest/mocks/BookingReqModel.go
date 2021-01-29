// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import data "bookingmgr/internal/data"
import mock "github.com/stretchr/testify/mock"

// BookingReqModel is an autogenerated mock type for the BookingReqModel type
type BookingReqModel struct {
	mock.Mock
}

// Do provides a mock function with given fields: userID, sourceLatitude, sourceLongitude, destLatitude, destLongitude
func (_m *BookingReqModel) Do(userID int, sourceLatitude int, sourceLongitude int, destLatitude int, destLongitude int) (*data.Booking, error) {
	ret := _m.Called(userID, sourceLatitude, sourceLongitude, destLatitude, destLongitude)

	var r0 *data.Booking
	if rf, ok := ret.Get(0).(func(int, int, int, int, int) *data.Booking); ok {
		r0 = rf(userID, sourceLatitude, sourceLongitude, destLatitude, destLongitude)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data.Booking)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int, int, int, int) error); ok {
		r1 = rf(userID, sourceLatitude, sourceLongitude, destLatitude, destLongitude)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}