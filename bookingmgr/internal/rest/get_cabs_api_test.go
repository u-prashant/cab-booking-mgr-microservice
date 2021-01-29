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

func TestGetCabsHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		inModelMock func() *mocks.GetCabsModel
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
			fields: fields{inModelMock: func() *mocks.GetCabsModel {
				output := []*data.Cab{{ID: 1, LicensePlate: "HR 26 AJ 2827", Location: data.Location{Latitude: 2, Longitude: 2}, Driver: data.Driver{ID: 1, Name: "Prashant"}}}
				mockGetCabsModel := &mocks.GetCabsModel{}
				mockGetCabsModel.On("Do", 2, 2, true).Return(output, nil).Once()
				return mockGetCabsModel
			}},
			args: args{
				request: func() *http.Request {
					req, err := http.NewRequest("GET", "/api/v1/cabs?latitude=2&longitude=2&availability=true", nil)
					require.NoError(t, err)
					return mux.SetURLVars(req, nil)
				},
			},
			expectedStatus:  http.StatusOK,
			expectedPayload: `[{"ID":1,"LicensePlate":"HR 26 AJ 2827","Location":{"Latitude":2,"Longitude":2},"Driver":{"ID":1,"Name":"Prashant"}}]` + "\n",
		},
		{
			name: "bad request",
			fields: fields{inModelMock: func() *mocks.GetCabsModel {
				// expect the model not to be called
				mockGetCabsModel := &mocks.GetCabsModel{}
				return mockGetCabsModel
			}},
			args: args{
				request: func() *http.Request {
					req, err := http.NewRequest("GET", "/api/v1/cabs?latitude=2&availability=true", nil)
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
			h := NewGetCabsHandler(tt.fields.inModelMock())

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
