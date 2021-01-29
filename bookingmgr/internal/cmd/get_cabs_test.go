package cmd

import (
	"bookingmgr/internal/data"
	"reflect"
	"testing"
)

func TestCabsGetter_Do(t *testing.T) {
	type args struct {
		latitude     int
		longitude    int
		availability bool
	}
	tests := []struct {
		name    string
		args    args
		want    []*data.Cab
		wantErr bool
		getCabs func(latitude int, longitude int, availability bool) ([]*data.Cab, error)
	}{
		{
			name: "sunny scenario",
			args: args{2, 2, true},
			getCabs: func(latitude int, longitude int, availability bool) ([]*data.Cab, error) {
				return []*data.Cab{
					{
						ID:           1,
						LicensePlate: "HR 26 AJ 2827",
						Location:     data.Location{Latitude: 2, Longitude: 2},
					},
				}, nil
			},
			want: []*data.Cab{
				{
					ID:           1,
					LicensePlate: "HR 26 AJ 2827",
					Location:     data.Location{Latitude: 2, Longitude: 2},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// copy the original handler for restoring
			original := getCabs
			cg := &CabsGetter{}
			getCabs = tt.getCabs
			got, err := cg.Do(tt.args.latitude, tt.args.longitude, tt.args.availability)
			if (err != nil) != tt.wantErr {
				t.Errorf("CabsGetter.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CabsGetter.Do() = %v, want %v", got, tt.want)
			}
			// revert the original handler
			getCabs = original
		})
	}
}
