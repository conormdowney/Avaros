/*
	The rooom rest service.
*/

package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gocraft/web"
	"github.com/jackc/pgx/v4/pgxpool"

	dataAccess "avaros/dataAccess"
	database "avaros/database"
	router "avaros/router"
	test "avaros/test"
)

func TestCheckReservation(t *testing.T) {

	db, router := setup()
	defer test.CloseDb(db)

	req, err := http.NewRequest("GET", "/room/check-reservation/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "userId",
		Value: "1",
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatal(err)
	}

	resRsp := ReservationResponse{}
	json.Unmarshal(rr.Body.Bytes(), &resRsp)

	if resRsp.Result {
		t.Fatal("Reservation should not exist")
	}

	_, err = dataAccess.Reserve(1, 0, db)
	if err != nil {
		t.Errorf("Error reserving a room: %s", err.Error())
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatal(err)
	}

	resRsp = ReservationResponse{}
	json.Unmarshal(rr.Body.Bytes(), &resRsp)

	if !resRsp.Result {
		t.Fatal("Reservation should exist")
	}
}

func TestDeleteReservation(t *testing.T) {

	db, router := setup()
	defer test.CloseDb(db)

	_, err := dataAccess.Reserve(1, 0, db)
	if err != nil {
		t.Errorf("Error reserving a room: %s", err.Error())
	}

	req, err := http.NewRequest("DELETE", "/room/delete-reservation/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "userId",
		Value: "1",
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatal(err)
	}

	exists, err := dataAccess.CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a reservation: %s", err.Error())
	}

	if exists {
		t.Errorf("Reservation should have been deleted")
	}
}

func TestReservation(t *testing.T) {
	db, router := setup()
	defer test.CloseDb(db)

	resReq := ReservationRequest{}
	jsonStr, err := json.Marshal(resReq)

	req, err := http.NewRequest("POST", "/room/reserve/1", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "userId",
		Value: "1",
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatal(err)
	}

	exists, err := dataAccess.CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a reservation: %s", err.Error())
	}

	if !exists {
		t.Errorf("Reservation should exist")
	}
}

func TestReservationExpiry(t *testing.T) {
	db, router := setup()
	defer test.CloseDb(db)

	resReq := ReservationRequest{
		ReservationLength: 1,
	}
	jsonStr, err := json.Marshal(resReq)

	req, err := http.NewRequest("POST", "/room/reserve/1", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "userId",
		Value: "1",
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatal(err)
	}

	time.Sleep(70 * time.Second)
	exists, err := dataAccess.CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a reservation: %s", err.Error())
	}

	if exists {
		t.Errorf("Reservation should have expired")
	}
}

func TestFutureReservation(t *testing.T) {
	db, router := setup()
	defer test.CloseDb(db)

	resReq := ReservationRequest{
		StartTime: time.Now().Add(time.Minute * 2),
	}
	jsonStr, err := json.Marshal(resReq)

	req, err := http.NewRequest("POST", "/room/reserve/1", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "userId",
		Value: "1",
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatal(err)
	}

	exists, err := dataAccess.CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a reservation: %s", err.Error())
	}

	if exists {
		t.Errorf("Reservation should not exist yet")
	}

	time.Sleep(70 * time.Second)
	exists, err = dataAccess.CheckReservation(1, db)
	if err != nil {
		t.Errorf("Error checking a reservation: %s", err.Error())
	}

	if !exists {
		t.Errorf("Reservation should exist")
	}
}

func setup() (*pgxpool.Pool, *web.Router) {
	db := test.NewDatabase()
	database.Seed(db)
	router := router.NewRouter()

	RestObj := RestServiceObject{
		router,
		db,
	}

	restService := &RoomService{RestObj}

	err := restService.Init()
	if err != nil {
		panic(err)
	}

	return db, router
}
