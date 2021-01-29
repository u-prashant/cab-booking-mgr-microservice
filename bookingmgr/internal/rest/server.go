package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Server is the HTTP REST server
type Server struct {
	address string
	server  *http.Server

	handlerGetBookings   http.Handler
	handlerGetCabsNearby http.Handler
	handlerBookingReq    http.Handler
}

// configure the endpoints to handlers
func (s *Server) buildRouter() http.Handler {
	router := mux.NewRouter()

	// map URL endpoints to HTTP handlers
	router.Handle("/api/v1/users/{id}/bookings", s.handlerGetBookings).Methods("GET")
	router.Handle("/api/v1/cabs", s.handlerGetCabsNearby).Methods("GET")
	router.Handle("/api/v1/users/{id}/book", s.handlerBookingReq).Methods("POST")

	return router
}

// New will create and initialize the server
func New(address string, getBookingsModel GetBookingsModel, getCabsModel GetCabsModel, bookingReqModel BookingReqModel) *Server {
	return &Server{
		address:              address,
		handlerGetBookings:   NewGetBookingsHandler(getBookingsModel),
		handlerGetCabsNearby: NewGetCabsHandler(getCabsModel),
		handlerBookingReq:    NewBookingReqHandler(bookingReqModel),
	}
}

// Listen for start a HTTP rest for this service
func (s *Server) Listen(stop <-chan struct{}) {
	router := s.buildRouter()

	// create the HTTP server
	s.server = &http.Server{
		Handler: router,
		Addr:    s.address,
	}

	// listen for shutdown
	go func() {
		// wait for shutdown signal
		<-stop

		_ = s.server.Close()
	}()

	_ = s.server.ListenAndServe()
}
