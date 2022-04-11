/*
	Class that holds the data access functions for any reservations.
*/
package dataAccess

import (
	database "avaros/database"
	test "avaros/test"

	"testing"
	"time"
)

func TestReservation(t *testing.T) {
	db := test.NewDatabase()
	defer test.CloseDb(db)

	database.Seed(db)

	reservationExists, err := CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a room: %s", err.Error())
	}

	if reservationExists {
		t.Errorf("No reservations should exist")
	}

	_, err = Reserve(1, 0, db)
	if err != nil {
		t.Errorf("Error reserving a room: %s", err.Error())
	}

	reservationExists, err = CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a room: %s", err.Error())
	}

	if !reservationExists {
		t.Errorf("Reservation for room 1 should exist")
	}
}

func TestDeleteReservation(t *testing.T) {
	db := test.NewDatabase()
	defer test.CloseDb(db)

	database.Seed(db)

	_, err := Reserve(1, 0, db)
	if err != nil {
		t.Errorf("Error reserving a room: %s", err.Error())
	}

	reservationExists, err := CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a reservation: %s", err.Error())
	}

	if !reservationExists {
		t.Errorf("Reservation for room 1 should exist")
	}

	err = DeleteReservation(1, db)
	if err != nil {
		t.Errorf("Error deleting a reservation: %s", err.Error())
	}

	reservationExists, err = CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a reservation: %s", err.Error())
	}

	if reservationExists {
		t.Errorf("Reservation for room 1 should have been deleted")
	}
}

func TestCheckRoomExists(t *testing.T) {
	db := test.NewDatabase()
	defer test.CloseDb(db)

	database.Seed(db)

	exists, err := CheckRoomExists(1, db)
	if err != nil {
		t.Errorf("Error checking a room: %s", err.Error())
	}

	if !exists {
		t.Errorf("Room 1 should now exist")
	}
}

func TestCreateFutureReservation(t *testing.T) {
	db := test.NewDatabase()
	defer test.CloseDb(db)

	database.Seed(db)

	go CreateFutureReservation(1, 1, 0, db)

	reservationExists, err := CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a reservation: %s", err.Error())
	}

	if reservationExists {
		t.Errorf("Reservation for room 1 should not exist")
	}

	time.Sleep(70 * time.Minute)

	reservationExists, err = CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a reservation: %s", err.Error())
	}

	if !reservationExists {
		t.Errorf("Reservation for room 1 should exist")
	}
}

func TestExpireReservation(t *testing.T) {
	db := test.NewDatabase()
	defer test.CloseDb(db)

	database.Seed(db)

	_, err := Reserve(1, 1, db)
	if err != nil {
		t.Errorf("Error reserving a room: %s", err.Error())
	}

	reservationExists, err := CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a reservation: %s", err.Error())
	}

	if !reservationExists {
		t.Errorf("Reservation for room 1 should exist")
	}

	time.Sleep(80 * time.Second)

	reservationExists, err = CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a reservation: %s", err.Error())
	}

	if reservationExists {
		t.Errorf("Reservation for room 1 should have expired")
	}
}
