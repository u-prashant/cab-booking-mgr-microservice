package data

import (
	"bookingmgr/internal/config"
	"context"
	"database/sql"
	"errors"

	// import the MySQL Driver
	_ "github.com/go-sql-driver/mysql"

	log "github.com/sirupsen/logrus"
)

const (
	// default user id (returned on error)
	defaultUserID = 0

	// SQL statements
	sqlGetAllBookings = `SELECT b.id, b.time, d.id, d.full_name, c.id, c.license_plate, b.source_latitude, 
						 b.source_longitude, b.destination_latitude, b.destination_longitude 
						 FROM users AS u INNER JOIN bookings AS b on u.id = b.user_id
						 INNER JOIN drivers AS d ON b.driver_id = d.id
						 INNER JOIN cabs AS c ON c.id = b.cab_id
						 WHERE u.id = ?`

	sqlGetAllCabsNearby = `SELECT c.id, c.license_plate, c.latitude, c.longitude, d.id, d.full_name
						   FROM cabs AS c INNER JOIN drivers AS d ON d.cab_id = c.id
						   WHERE c.availability = ?
						   ORDER BY (sqrt(power((c.latitude-?),2)+power((c.longitude-?),2))) LIMIT 10`

	sqlUpdateCabAvailability = `UPDATE cabs SET availability = ? WHERE id = ?`

	sqlInsertBooking = `INSERT INTO bookings (user_id, time, driver_id, cab_id, source_latitude, source_longitude,
						destination_latitude, destination_longitude) VALUES (?,?,?,?,?,?,?,?)`
)

var (
	db *sql.DB

	// ErrNotFound is returned when the no records were matched by the query
	ErrNotFound = errors.New("not found")
)

var getDB = func() (*sql.DB, error) {
	if db == nil {
		if config.App == nil {
			return nil, errors.New("config is not initialized")
		}

		var err error
		db, err = sql.Open("mysql", config.App.DSN)
		if err != nil {
			panic(err.Error())
		}
	}

	return db, nil
}

func init() {
	// ensure the config is loaded and the db initialised
	_, _ = getDB()
}

// Booking is the data object for any booking
type Booking struct {
	// ID is the booking id
	ID int64
	// Date is the booking date of the cab
	Date string
	// Source is the location from where the pickup was done
	Source Location
	// Destionation is the location from where the drop was done
	Destination Location
	// Driver contains the driver details who drove the cab
	Driver Driver
	// Cab contains the cab details in which the user commuted
	Cab Cab
}

// Driver is the data object for driver's details
type Driver struct {
	// ID is the driver id
	ID int
	// Name is the driver name
	Name string
}

// Cab is the data object for cab's details
type Cab struct {
	// ID is the cab id
	ID int
	// Name is the cab name
	LicensePlate string
	// Locations is current coordinates of the cab
	Location Location
	// Driver contains the driver details
	Driver Driver
}

// Location contains the coordinates of any location
type Location struct {
	Latitude  int
	Longitude int
}

// GetAllBookings will attempt to load all the bookings from the database.
// It will return ErrNotFound when there are no bookings in the database.
// Any other errors returned are caused by the underlying database or our connection to it.
func GetAllBookings(ID int) ([]*Booking, error) {
	db, err := getDB()
	if err != nil {
		log.Errorf("failed to get DB connection. err: %s", err)
		return nil, err
	}

	// perform DB select
	rows, err := db.Query(sqlGetAllBookings, ID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var bookings []*Booking

	for rows.Next() {
		booking := &Booking{}
		err := rows.Scan(&booking.ID,
			&booking.Date,
			&booking.Driver.ID,
			&booking.Driver.Name,
			&booking.Cab.ID,
			&booking.Cab.LicensePlate,
			&booking.Source.Latitude,
			&booking.Source.Longitude,
			&booking.Destination.Latitude,
			&booking.Destination.Longitude)
		if err != nil {
			log.Errorf("failed to scan query result. err: %s", err)
			return nil, err
		}

		bookings = append(bookings, booking)
	}

	if len(bookings) == 0 {
		log.Debugf("no bookings found in the database. err: %s", err)
		return nil, ErrNotFound
	}

	return bookings, nil
}

// GetAllCabsNearby will attempt to load all the cabs nearby from the database.
// It will return ErrNotFound when there are no bookings found.
// Any other errors returned are caused by the underlying database or our connection to it.
func GetAllCabsNearby(latitude int, longitude int, availability bool) ([]*Cab, error) {
	db, err := getDB()
	if err != nil {
		log.Errorf("failed to get DB connection. err: %s", err)
		return nil, err
	}

	// perform DB select
	rows, err := db.Query(sqlGetAllCabsNearby, availability, latitude, longitude)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var cabs []*Cab

	for rows.Next() {
		cab := &Cab{}
		err := rows.Scan(&cab.ID,
			&cab.LicensePlate,
			&cab.Location.Latitude,
			&cab.Location.Longitude,
			&cab.Driver.ID,
			&cab.Driver.Name)
		if err != nil {
			log.Errorf("failed to scan query result. err: %s", err)
			return nil, err
		}

		cabs = append(cabs, cab)
	}

	if len(cabs) == 0 {
		log.Debugf("no nearby cabs found in the database. err: %s", err)
		return nil, ErrNotFound
	}

	return cabs, nil
}

// SaveBooking will save the booking into the database.
// Any errors returned are caused by the underlying database or our connection to it.
func SaveBooking(userID int, b *Booking) error {
	db, err := getDB()
	if err != nil {
		log.Errorf("failed to get DB connection. err: %s", err)
		return err
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("failed to create transaction for db. err: %s", err)
		return err
	}
	_, err = tx.ExecContext(ctx, sqlUpdateCabAvailability, 1, b.Cab.ID)
	if err != nil {
		log.Errorf("failed to update cab availability. err: %s", err)
		tx.Rollback()
		return err
	}

	result, err := tx.ExecContext(ctx, sqlInsertBooking, userID, b.Date, b.Driver.ID, b.Cab.ID,
		b.Source.Latitude, b.Source.Longitude, b.Destination.Latitude, b.Destination.Longitude)
	if err != nil {
		log.Errorf("failed to insert the booking. err: %s", err)
		tx.Rollback()
		return err
	}
	b.ID, err = result.LastInsertId()
	if err != nil {
		log.Errorf("failed to retrieve booking id. err: %s", err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("failed to commit the transaction. err: %s", err)
		return err
	}

	return nil
}
