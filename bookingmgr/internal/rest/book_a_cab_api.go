package rest

import (
	"bookingmgr/internal/data"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// BookingReqModel will book a cab
type BookingReqModel interface {
	Do(userID int, sourceLatitude int, sourceLongitude int, destLatitude int, destLongitude int) (*data.Booking, error)
}

// NewBookingReqHandler is the constructor for BookingReqHandler
func NewBookingReqHandler(model BookingReqModel) *BookingReqHandler {
	return &BookingReqHandler{
		model: model,
	}
}

// BookingReqHandler is the HTTP handler
type BookingReqHandler struct {
	model BookingReqModel
}

// ServeHTTP implements http.Handler
func (h *BookingReqHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	// extract user id from request
	userID, err := h.extractID(request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// extract payload from request
	requestPayload, err := h.extractPayload(request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// book a cab
	bookingInfo, err := h.model.Do(userID, requestPayload.Source.Latitude, requestPayload.Source.Longitude,
		requestPayload.Destination.Latitude, requestPayload.Destination.Longitude)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// writing the response
	err = json.NewEncoder(response).Encode(bookingInfo)
	if err != nil {
		log.Errorf("error in writing the response. err: %s", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// extract payload from request
func (h *BookingReqHandler) extractPayload(request *http.Request) (*bookingRequest, error) {
	requestPayload := &bookingRequest{}

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(requestPayload)
	return requestPayload, err
}

// extract the userID from the request
func (h *BookingReqHandler) extractID(request *http.Request) (int, error) {
	vars := mux.Vars(request)
	idAsString, exists := vars["id"]
	if !exists {
		err := errors.New("[get bookings] user id missing from request")
		log.Warn(err.Error())
		return defaultUserID, err
	}

	// convert ID to int
	id, err := strconv.Atoi(idAsString)
	if err != nil {
		err = fmt.Errorf("[get bookings] failed to convert user id into a number. err: %s", err)
		log.Error(err.Error())
		return defaultUserID, err
	}

	return id, nil
}

// booking request format
type bookingRequest struct {
	Source      location
	Destination location
}

type location struct {
	Latitude  int
	Longitude int
}
