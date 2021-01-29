package rest

import (
	"bookingmgr/internal/data"
	"bookingmgr/internal/rest/mocks"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBookingsHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		inModelMock func() *mocks.GetBookingsModel
	}
	type args struct {
		request func() *http.Request
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectedStatus  int
		expectedPayload string
	}{
		{
			name: "sunny day scenario",
			fields: fields{
				inModelMock: func() *mocks.GetBookingsModel {
					output := []*data.Booking{
						{
							ID:          1,
							Date:        "10/12/2020",
							Source:      data.Location{Latitude: 101, Longitude: 101},
							Destination: data.Location{Latitude: 202, Longitude: 202},
							Driver:      data.Driver{ID: 1, Name: "Prashant"},
							Cab:         data.Cab{ID: 1, LicensePlate: "HR 26 AJ 2827"},
						},
					}
					mockGetBookingsModel := &mocks.GetBookingsModel{}
					mockGetBookingsModel.On("Do", 1).Return(output, nil).Once()
					return mockGetBookingsModel
				}},
			args: args{
				request: func() *http.Request {
					req, err := http.NewRequest("GET", "/api/v1/users/1/bookings", nil)
					require.NoError(t, err)
					return mux.SetURLVars(req, map[string]string{"id": "1"})
				},
			},
			expectedStatus:  http.StatusOK,
			expectedPayload: `[{"ID":1,"Date":"10/12/2020","Source":{"Latitude":101,"Longitude":101},"Destination":{"Latitude":202,"Longitude":202},"Driver":{"ID":1,"Name":"Prashant"},"Cab":{"ID":1,"LicensePlate":"HR 26 AJ 2827","Location":{"Latitude":0,"Longitude":0},"Driver":{"ID":0,"Name":""}}}]` + "\n",
		},
		{
			name: "bad request",
			fields: fields{
				inModelMock: func() *mocks.GetBookingsModel {
					// expect the model not to be called
					mockGetBookingsModel := &mocks.GetBookingsModel{}
					return mockGetBookingsModel
				}},
			args: args{
				request: func() *http.Request {
					req, err := http.NewRequest("GET", "/api/v1/users/0/bookings", nil)
					require.NoError(t, err)
					return mux.SetURLVars(req, nil)
				},
			},
			expectedStatus:  http.StatusBadRequest,
			expectedPayload: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// build handler
			h := NewGetBookingsHandler(tt.fields.inModelMock())

			// perform request
			response := httptest.NewRecorder()
			h.ServeHTTP(response, tt.args.request())

			// validate outputs
			require.Equal(t, tt.expectedStatus, response.Code, tt.name)

			payload, _ := ioutil.ReadAll(response.Body)
			assert.Equal(t, tt.expectedPayload, string(payload), tt.name)
		})
	}
}
