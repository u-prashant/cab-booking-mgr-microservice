package cmd

import (
	"bookingmgr/internal/data"
	"errors"
)

var (
	// error thrown when the request user is not in the database
	errBookingsNotFound = errors.New("bookings not found")
)

// BookingsGetter attempts to get bookings.
type BookingsGetter struct {
}

// Do will perform the get
func (bg *BookingsGetter) Do(ID int) ([]*data.Booking, error) {
	// load bookings from the data layer
	bookings, err := getAllBookings(ID)
	if err != nil {
		if err == data.ErrNotFound {
			return nil, errBookingsNotFound
		}
		return nil, err
	}

	return bookings, err
}

var getAllBookings = data.GetAllBookings
