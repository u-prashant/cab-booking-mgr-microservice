package cmd

import (
	"bookingmgr/internal/data"
	"reflect"
	"testing"
)

func TestBookingsGetter_Do(t *testing.T) {
	type args struct {
		ID int
	}
	tests := []struct {
		name           string
		bg             *BookingsGetter
		args           args
		want           []*data.Booking
		wantErr        bool
		getAllBookings func(ID int) ([]*data.Booking, error)
	}{
		{
			name: "sunny scenario",
			args: args{1},
			getAllBookings: func(ID int) ([]*data.Booking, error) {
				return []*data.Booking{
					{
						ID:          1,
						Date:        "12/12/2020",
						Source:      data.Location{Latitude: 1, Longitude: 1},
						Destination: data.Location{Latitude: 2, Longitude: 2},
						Driver:      data.Driver{ID: 1, Name: "Prashant"},
						Cab:         data.Cab{ID: 1, LicensePlate: "HR 26 AJ 2827"},
					},
				}, nil
			},
			want: []*data.Booking{
				{
					ID:          1,
					Date:        "12/12/2020",
					Source:      data.Location{Latitude: 1, Longitude: 1},
					Destination: data.Location{Latitude: 2, Longitude: 2},
					Driver:      data.Driver{ID: 1, Name: "Prashant"},
					Cab:         data.Cab{ID: 1, LicensePlate: "HR 26 AJ 2827"},
				},
			},
			wantErr: false,
		},
		{
			name: "no records found error",
			args: args{1},
			getAllBookings: func(ID int) ([]*data.Booking, error) {
				return nil, data.ErrNotFound
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// copy the original handler for restoring
			original := getAllBookings
			bg := &BookingsGetter{}
			getAllBookings = tt.getAllBookings
			got, err := bg.Do(tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BookingsGetter.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BookingsGetter.Do() = %v, want %v", got, tt.want)
			}
			// restore the original handler
			getAllBookings = original
		})
	}
}
