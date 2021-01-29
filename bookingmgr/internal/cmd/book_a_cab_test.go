package cmd

import (
	"bookingmgr/internal/data"
	"reflect"
	"testing"
)

func TestCabBooker_Do(t *testing.T) {
	type args struct {
		userID          int
		sourceLatitude  int
		sourceLongitude int
		destLatitude    int
		destLongitude   int
	}
	tests := []struct {
		name        string
		args        args
		want        *data.Booking
		wantErr     bool
		getCabs     func(latitude int, longitude int, availability bool) ([]*data.Cab, error)
		saveBooking func(userID int, b *data.Booking) error
	}{
		{
			name: "sunny day scenario",
			args: args{
				userID:          1,
				sourceLatitude:  1,
				sourceLongitude: 1,
				destLatitude:    2,
				destLongitude:   2,
			},
			getCabs: func(latitude int, longitude int, availability bool) ([]*data.Cab, error) {
				return []*data.Cab{
					{
						ID:           1,
						LicensePlate: "HR 26 AJ 2827",
						Driver:       data.Driver{ID: 1, Name: "Prashant"},
					},
				}, nil
			},
			saveBooking: func(userID int, b *data.Booking) error {
				b.ID = 1
				return nil
			},
			want: &data.Booking{
				ID:          1,
				Source:      data.Location{Latitude: 1, Longitude: 1},
				Destination: data.Location{Latitude: 2, Longitude: 2},
				Driver:      data.Driver{ID: 1, Name: "Prashant"},
				Cab:         data.Cab{ID: 1, LicensePlate: "HR 26 AJ 2827", Driver: data.Driver{ID: 1, Name: "Prashant"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// copy the original handler for restoring
			originalGetCabs := getCabs
			getCabs = tt.getCabs

			originalSaveBooking := saveBooking
			cb := &CabBooker{}
			saveBooking = tt.saveBooking

			got, err := cb.Do(tt.args.userID, tt.args.sourceLatitude, tt.args.sourceLongitude, tt.args.destLatitude, tt.args.destLongitude)
			if (err != nil) != tt.wantErr {
				t.Errorf("CabBooker.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				tt.want.Date = got.Date
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CabBooker.Do() = %v, want %v", got, tt.want)
			}

			getCabs = originalGetCabs
			saveBooking = originalSaveBooking
		})
	}
}
