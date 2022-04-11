/*
	Class that holds the data access functions for any reservations.
*/
package dataAccess

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

//CheckReservation checks that a reservation for a given room exists
func CheckReservation(roomId int32, db *pgxpool.Pool) (bool, error) {
	if db == nil {
		return false, errors.New("Database instance empty")
	}

	//query for reservations on the room
	rows, err := db.Query(context.Background(), `
		SELECT 
			id 
		FROM
			reservation
		Where
			room_id = $1 
		AND 
			expired = false
	`, roomId)
	if err != nil {
		return false, err
	}

	// rows.Next() will be false if no rows exist and true if one or more does
	isNext := rows.Next()
	rows.Close()
	return isNext, nil
}

// Reserve creates a reservation for a room. If an expiry time is supplied,
// a thread opens that will count down the time from the reservation
func Reserve(roomId int32, expiryTime int, db *pgxpool.Pool) (int32, error) {
	id := -1
	// set the start time of the reservation
	startTime := time.Now()
	err := db.QueryRow(context.Background(), `
	INSERT INTO 
		reservation (room_id, start_time)
	VALUES 
		($1, $2)
	RETURNING id
	`, roomId, startTime).Scan(&id)

	if err != nil {
		panic("Error reserving room: " + err.Error())
	}

	// if the expiry time is provided then fire off a thread that
	// will handle expiring the reservation
	if expiryTime > 0 {
		go expireReservation(expiryTime, int32(id), db)
	}
	// return the id of the reservation
	return int32(id), nil
}

// DeleteReservation deletes a reservation for a room
func DeleteReservation(roomId int32, db *pgxpool.Pool) error {

	_, err := db.Exec(context.Background(), `
	DELETE 
	FROM 
		reservation
	WHERE 
		room_id = $1
	`, roomId)

	if err != nil {
		return err
	}

	return nil
}

// expireReservation starts a timer to expire a reservation after the supplied
// number of minutes. Should only be opened in a thread
func expireReservation(numMins int, reservationId int32, db *pgxpool.Pool) {
	for {
		select {
		// thread will run and after the number of minutes provided, will
		// expire the reservation
		case <-time.After(time.Minute * time.Duration(numMins)):
			endTime := time.Now()
			_, err := db.Exec(context.Background(), `
				UPDATE reservation
				SET expired = true, end_time = $2
				WHERE id = $1
			`, reservationId, endTime)
			if err != nil {
				panic("Error updating the reservations expiry: " + err.Error())
			}
			return
		}
	}
}

// CreateFutureReservation starts a timer to create a reservation after the supplied
// number of minutes. Should only be opened in a thread. Calls Reserve, which handles
// the rest of the reserve workflow
func CreateFutureReservation(timeInFuture float64, roomId int32, expiryTime int, db *pgxpool.Pool) {
	fmt.Println(time.Duration(timeInFuture))
	for {
		select {
		// Create the reservation in the future
		case <-time.After(time.Minute * time.Duration(timeInFuture)):
			exists, err := CheckReservation(roomId, db)
			if exists {
				fmt.Println("Reservation already exists")
				return
			}
			_, err = Reserve(roomId, expiryTime, db)
			if err != nil {
				panic("Error creating reservation: " + err.Error())
			}
			return
		}
	}
}

// CheckRoomExists checks if a room id supplied is in the database to reserve
func CheckRoomExists(id int32, db *pgxpool.Pool) (bool, error) {
	//query for reservations on the room
	rows, err := db.Query(context.Background(), `
		SELECT 
			id 
		FROM
			room
		Where
			id = $1
	`, id)
	if err != nil {
		return false, err
	}

	next := rows.Next()
	rows.Close()

	return next, nil
}
