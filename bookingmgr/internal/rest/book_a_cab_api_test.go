package rest

import (
	"bookingmgr/internal/data"
	"bookingmgr/internal/rest/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getRequest(t *testing.T, id int, bodyJSON bookingRequest) *http.Request {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(bodyJSON)
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/v1/users/%d/book", id), b)
	require.NoError(t, err)
	if id == 0 {
		mux.SetURLVars(req, nil)
	}
	return mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func TestBookingReqHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		modelMock func() *mocks.BookingReqModel
	}
	type args struct {
		bodyFormat bookingRequest
		userID     int
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
				modelMock: func() *mocks.BookingReqModel {
					output := &data.Booking{
						ID:          1,
						Source:      data.Location{Latitude: 1, Longitude: 1},
						Destination: data.Location{Latitude: 2, Longitude: 2},
						Driver:      data.Driver{ID: 1, Name: "Prashant"},
						Cab:         data.Cab{ID: 1, LicensePlate: "HR 26 AJ 2827"},
					}
					mockBookingModel := &mocks.BookingReqModel{}
					mockBookingModel.On("Do", 1, 1, 1, 2, 2).Return(output, nil)
					return mockBookingModel
				},
			},
			args: args{
				bodyFormat: bookingRequest{
					Source: location{
						Latitude:  1,
						Longitude: 1,
					},
					Destination: location{
						Latitude:  2,
						Longitude: 2,
					},
				},
				userID: 1,
			},
			expectedStatus:  http.StatusOK,
			expectedPayload: `{"ID":1,"Date":"","Source":{"Latitude":1,"Longitude":1},"Destination":{"Latitude":2,"Longitude":2},"Driver":{"ID":1,"Name":"Prashant"},"Cab":{"ID":1,"LicensePlate":"HR 26 AJ 2827","Location":{"Latitude":0,"Longitude":0},"Driver":{"ID":0,"Name":""}}}` + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// build handler
			h := NewBookingReqHandler(tt.fields.modelMock())

			// perform request
			response := httptest.NewRecorder()
			h.ServeHTTP(response, getRequest(t, tt.args.userID, tt.args.bodyFormat))

			// validate outputs
			require.Equal(t, tt.expectedStatus, response.Code, tt.name)
			responsePayload := &data.Booking{}
			decoder := json.NewDecoder(response.Body)
			_ = decoder.Decode(responsePayload)
			responsePayload.Date = ""
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(responsePayload)
			assert.Equal(t, tt.expectedPayload, b.String(), tt.name)
		})
	}
}
