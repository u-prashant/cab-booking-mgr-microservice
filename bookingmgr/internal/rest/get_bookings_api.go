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

const (
	// default user id (returned on error)
	defaultUserID = 0
)

// GetBookingsModel will load the bookings
type GetBookingsModel interface {
	Do(ID int) ([]*data.Booking, error)
}

// NewGetBookingsHandler is the constructor for GetBookingsHandler
func NewGetBookingsHandler(model GetBookingsModel) *GetBookingsHandler {
	return &GetBookingsHandler{
		getter: model,
	}
}

// GetBookingsHandler is the HTTP handler for the "Get Bookings" endpoint
type GetBookingsHandler struct {
	getter GetBookingsModel
}

// ServeHTTP implements http.Handler
func (h *GetBookingsHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	// extract user id from request
	id, err := h.extractID(request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// attempt get
	bookings, err := h.getter.Do(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if no error, call to http.ResponseWriter.Write() will cause HTTP OK (200)
	err = json.NewEncoder(response).Encode(bookings)
	if err != nil {
		log.Errorf("error in writing the response. err: %s", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// extract the userID from the request
func (h *GetBookingsHandler) extractID(request *http.Request) (int, error) {
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
