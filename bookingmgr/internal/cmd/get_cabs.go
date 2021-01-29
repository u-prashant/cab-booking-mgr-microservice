package cmd

import (
	"bookingmgr/internal/data"
	"errors"
)

var (
	// error thrown when there are no cabs
	errCabsNotFound = errors.New("cabs not found")
)

// CabsGetter will attempt to load cabs.
// It can return an error when the cab is not found.
type CabsGetter struct {
}

// Do will perform the get
func (cg *CabsGetter) Do(latitude int, longitude int, availability bool) ([]*data.Cab, error) {
	// load cabs from the data layer
	cabs, err := getCabs(latitude, longitude, availability)
	if err != nil {
		if err == data.ErrNotFound {
			return nil, errCabsNotFound
		}
		return nil, err
	}

	return cabs, err
}

// this function as a variable allow us to monkey patch during testing
var getCabs = data.GetAllCabsNearby
