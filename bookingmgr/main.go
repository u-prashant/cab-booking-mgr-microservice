package main

import (
	"bookingmgr/internal/cmd"
	"bookingmgr/internal/config"
	"bookingmgr/internal/rest"
	"context"
)

func main() {
	// bind stop channel to context
	ctx := context.Background()

	// build model layer
	getBookingsModel := &cmd.BookingsGetter{}
	getCabsModel := &cmd.CabsGetter{}
	bookingReqModel := &cmd.CabBooker{}

	// start REST server
	server := rest.New(config.App.Address, getBookingsModel, getCabsModel, bookingReqModel)
	server.Listen(ctx.Done())
}
