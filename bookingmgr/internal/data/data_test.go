package data

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllBookings(t *testing.T) {
	type args struct {
		ID int
	}
	tests := []struct {
		name            string
		args            args
		configureMockDB func(sqlmock.Sqlmock)
		want            []*Booking
		wantErr         bool
	}{
		{
			name: "sunny scenario",
			args: args{ID: 1},
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlGetAllBookings)
				dbMock.ExpectQuery(queryRegex).WillReturnRows(
					sqlmock.NewRows([]string{"booking.id", "booking.name", "driver.id", "driver.name", "cab.id", "cab.name",
						"booking.source_latitude", "booking.source_longitude", "booking.destination_latitude", "booking.destination_longitude"}).
						AddRow(1, "12/12/2020", 1, "Prashant", 1, "HR 26 AJ 2827", 1, 1, 2, 2))
			},
			want: []*Booking{
				{
					ID:          1,
					Date:        "12/12/2020",
					Source:      Location{1, 1},
					Destination: Location{2, 2},
					Driver:      Driver{ID: 1, Name: "Prashant"},
					Cab:         Cab{ID: 1, LicensePlate: "HR 26 AJ 2827"},
				},
			},
			wantErr: false,
		},
		{
			name: "any db error",
			args: args{ID: 1},
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlGetAllBookings)
				dbMock.ExpectQuery(queryRegex).WillReturnError(errors.New("some db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		// define a mock db
		testDb, dbMock, err := sqlmock.New()
		require.NoError(t, err)

		// configure the mock db
		tt.configureMockDB(dbMock)

		// monkey db for this test
		original := *db
		db = testDb

		// call function
		result, err := GetAllBookings(tt.args.ID)

		// validate results
		assert.Equal(t, tt.want, result, tt.name)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.NoError(t, dbMock.ExpectationsWereMet())

		// restore original DB (after test)
		db = &original
		testDb.Close()
	}
}

// convert SQL string to regex by treating the entire query as a literal
func convertSQLToRegex(in string) string {
	return `\Q` + in + `\E`
}

func TestGetAllCabsNearby(t *testing.T) {
	type args struct {
		latitude     int
		longitude    int
		availability bool
	}
	tests := []struct {
		name            string
		args            args
		configureMockDB func(sqlmock.Sqlmock)
		want            []*Cab
		wantErr         bool
	}{
		{
			name: "sunny day scenario",
			args: args{
				latitude:     2,
				longitude:    2,
				availability: true,
			},
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlGetAllCabsNearby)
				dbMock.ExpectQuery(queryRegex).WillReturnRows(
					sqlmock.NewRows([]string{"cab.id", "cab.license_plate", "cab.latitude", "cab.longitude", "driver.id", "driver.name"}).
						AddRow(1, "HR 26 AJ 2827", 2, 2, 1, "Prashant"))
			},
			want: []*Cab{
				{
					ID:           1,
					LicensePlate: "HR 26 AJ 2827",
					Location:     Location{Latitude: 2, Longitude: 2},
					Driver:       Driver{ID: 1, Name: "Prashant"},
				},
			},
			wantErr: false,
		},
		{
			name: "any db error",
			args: args{
				latitude:     2,
				longitude:    2,
				availability: true,
			},
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlGetAllCabsNearby)
				dbMock.ExpectQuery(queryRegex).WillReturnError(errors.New("some db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		// define a mock db
		testDb, dbMock, err := sqlmock.New()
		require.NoError(t, err)

		// configure the mock db
		tt.configureMockDB(dbMock)

		// monkey db for this test
		original := *db
		db = testDb

		// call function
		result, err := GetAllCabsNearby(tt.args.latitude, tt.args.longitude, tt.args.availability)

		// validate results
		assert.Equal(t, tt.want, result, tt.name)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.NoError(t, dbMock.ExpectationsWereMet())

		// restore original DB (after test)
		db = &original
		testDb.Close()
	}
}
