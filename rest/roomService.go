/*
	The rooom rest service.
*/

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"avaros/dataAccess"
	"avaros/models"

	"github.com/gocraft/web"
)

type RoomService struct {
	RestObj RestServiceObject
}

type ReservationResponse struct {
	Result bool    `json:"result"`
	Reason string  `json:"reason"`
	Ids    []int32 `json:"ids"`
}

// The request obejct when making a reservation. Can contain the the start time the reservation
// is for and the duration in the number of minutes
type ReservationRequest struct {
	StartTime         time.Time `json:"startTime"`
	ReservationLength int       `json:"reservationLength"`
}

// Init initialises the service and starts listening for its paths
func (rs *RoomService) Init() error {
	if rs.RestObj.Router == nil {
		return errors.New("A router must be present for the service to listen on")
	}

	rs.RestObj.Router.Post("/room/reserve/:id", rs.reserveRoom)
	rs.RestObj.Router.Delete("/room/delete-reservation/:id", rs.deleteReservation)
	rs.RestObj.Router.Get("/room/check-reservation/:id", rs.checkReservation)
	return nil
}

// reserveRoom reserves a room. If a start time is supplied it will not reserve the room until
// that time
func (rs *RoomService) reserveRoom(rw web.ResponseWriter, req *web.Request) {
	// Get the id from the url parameters
	roomId := getIdAsInt(req.PathParams["id"])

	roomExists, err := dataAccess.CheckRoomExists(roomId, rs.RestObj.Db)
	if err != nil {
		panic("Error determining if room exists: " + err.Error())
	}

	if !roomExists {
		panic(fmt.Sprintf("Room with id %d does not exist", roomId))
	}
	// read the request body to get the values passed in, if any
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		panic("Error reading request body: " + err.Error())
	}

	// Unmarshal the request body
	var resReq ReservationRequest
	err = json.Unmarshal(b, &resReq)
	if err != nil {
		panic("Error unmarshalling request body: " + err.Error())
	}
	// check if a reservation exists
	reservationExists, err := dataAccess.CheckReservation(roomId, rs.RestObj.Db)
	if err != nil {
		panic("Error checking reservation: " + err.Error())
	}

	resRsp := ReservationResponse{}
	// if a reservation exists then you cannot reserve the room.
	if reservationExists {
		resRsp.Result = false
		resRsp.Reason = "Reservation already exists."
	} else {
		// If there is no start time provided reserve now
		if resReq.StartTime.IsZero() {
			// call reserve and handle any error passed back
			reservationId, err := dataAccess.Reserve(roomId, resReq.ReservationLength, rs.RestObj.Db)
			if err != nil {
				panic("Error reserving room: " + err.Error())
			}
			resRsp.Result = true
			resRsp.Ids = []int32{reservationId}
		} else {
			// get the time difference for when to create the future reservation
			startTimeDelay := math.Abs(time.Now().Sub(resReq.StartTime).Minutes())
			//fmt.Println(startTimeDelay)
			// fire off a thread that will handle creating that reservation
			// at the correct time
			go dataAccess.CreateFutureReservation(startTimeDelay, roomId, resReq.ReservationLength, rs.RestObj.Db)
			resRsp.Result = true
		}
	}
	// send the response back to the client
	sendResponse(resRsp, rw)
}

// deleteReservation deletes a reservation
func (rs *RoomService) deleteReservation(rw web.ResponseWriter, req *web.Request) {
	// get the id of the room to delete the reservation for
	roomId := getIdAsInt(req.PathParams["id"])
	// check if the room has a reservation
	reservationExists, err := dataAccess.CheckReservation(roomId, rs.RestObj.Db)
	if err != nil {
		panic("Error checking reservation: " + err.Error())
	}

	resRsp := ReservationResponse{
		Result: reservationExists,
	}
	// if it doesnt return a message indicating so
	if !reservationExists {
		resRsp.Reason = fmt.Sprintf("Reservation for room %d does not exist exists.", roomId)
	} else {
		// otherwise delete the reservation
		err := dataAccess.DeleteReservation(roomId, rs.RestObj.Db)
		if err != nil {
			panic("Error deleting reservation: " + err.Error())
		}
	}

	sendResponse(resRsp, rw)
}

// checkReservation checks if a room has a reservation currently or not
func (rs *RoomService) checkReservation(rw web.ResponseWriter, req *web.Request) {
	roomId := getIdAsInt(req.PathParams["id"])
	// get the reservation status for the room passed in
	reservationExists, err := dataAccess.CheckReservation(roomId, rs.RestObj.Db)
	if err != nil {
		panic("Error checking reservation: " + err.Error())
	}

	resRsp := ReservationResponse{
		Result: reservationExists,
	}
	// if a reservation exists return a message saying so
	if reservationExists {
		resRsp.Reason = "Reservation already exists."
	}

	sendResponse(resRsp, rw)
}

// getIdAsInt converts the id passed in to the api from a string to a number
func getIdAsInt(roomIdStr string) int32 {
	// do not panic, log and leave function
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		panic("Error getting room id: " + err.Error())
	}

	return int32(roomId)
}

// sendResponse is used to send the response back to the client
func sendResponse(resRsp ReservationResponse, rw web.ResponseWriter) {
	rsp := models.RestResponse{
		Result:     resRsp,
		StatusCode: http.StatusOK,
	}

	err := json.NewEncoder(rw).Encode(rsp.Result)
	if err != nil {
		panic("Error sending response: " + err.Error())
	}
}
