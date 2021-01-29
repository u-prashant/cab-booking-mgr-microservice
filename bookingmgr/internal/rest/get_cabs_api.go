package rest

import (
	"bookingmgr/internal/data"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// GetCabsModel will load a registration
type GetCabsModel interface {
	Do(latitude int, longitude int, availability bool) ([]*data.Cab, error)
}

// NewGetCabsHandler is the constructor for GetCabsHandler
func NewGetCabsHandler(model GetCabsModel) *GetCabsHandler {
	return &GetCabsHandler{
		getter: model,
	}
}

// GetCabsHandler is the HTTP handler for the "Get Cabs" endpoint
type GetCabsHandler struct {
	getter GetCabsModel
}

// ServerHTTP implements http.Handler
func (h *GetCabsHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// extract params from the request
	latitude, longitude, availability, err := h.extractParams(request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// attempt get
	cabs, err := h.getter.Do(latitude, longitude, availability)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if no error, call to http.ResponseWriter.Write() will cause HTTP OK (200)
	err = json.NewEncoder(response).Encode(cabs)
	if err != nil {
		log.Errorf("error in writing the response. err: %s", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *GetCabsHandler) extractParams(request *http.Request) (int, int, bool, error) {
	// latitude
	lati := request.FormValue("latitude")
	if lati == "" {
		return 0, 0, false, errors.New("latitude not found in request")
	}
	latitude, err := strconv.Atoi(lati)
	if err != nil {
		err = fmt.Errorf("[get cabs nearby] failed to convert latitude into a number. err: %s", err)
		log.Error(err.Error())
		return 0, 0, false, err
	}

	// longitude
	longi := request.FormValue("longitude")
	if longi == "" {
		return 0, 0, false, errors.New("longitude not found in request")
	}
	longitude, err := strconv.Atoi(longi)
	if err != nil {
		err = fmt.Errorf("[get cabs nearby] failed to convert longitude into a number. err: %s", err)
		log.Error(err.Error())
		return 0, 0, false, err
	}

	// availability
	availability := false
	available := request.FormValue("availability")
	if available == "true" {
		availability = true
	}

	return latitude, longitude, availability, nil
}
