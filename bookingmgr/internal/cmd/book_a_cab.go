package cmd

import (
	"bookingmgr/internal/data"
	"errors"
	"time"
)

var (
	// error thrown when there are no cabs
	errBookingFailed = errors.New("booking failed")
)

// CabBooker will attempt to book a cab.
// It can return an error when the requested person is not found.
type CabBooker struct {
}

// Do will perform the get
func (cb *CabBooker) Do(userID int, sourceLatitude int, sourceLongitude int, destLatitude int, destLongitude int) (*data.Booking, error) {
	// load cabs from the data layer
	cabsGetter := &CabsGetter{}
	cabs, err := cabsGetter.Do(sourceLatitude, sourceLongitude, true)
	if err != nil || len(cabs) == 0 {
		return nil, errBookingFailed
	}

	// selecting the 1st cab
	booking := &data.Booking{
		Source:      data.Location{Latitude: sourceLatitude, Longitude: sourceLongitude},
		Destination: data.Location{Latitude: destLatitude, Longitude: destLongitude},
		Driver:      cabs[0].Driver,
		Cab:         *cabs[0],
		Date:        getCurrentTimestamp(),
	}
	err = saveBooking(userID, booking)
	if err != nil {
		return nil, errBookingFailed
	}

	return booking, nil
}

// this function as a variable allow us to monkey patch during testing
var saveBooking = data.SaveBooking

func getCurrentTimestamp() string {
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	return ts
}
