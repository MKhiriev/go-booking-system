package handlers

import (
	"encoding/json"
	"humoBooking/internal/models"
	"humoBooking/pkg"
	"log"
	"net/http"
	"strconv"
)

func (h *Handlers) DeleteRoom(w http.ResponseWriter, r *http.Request) { // DONE
	// get room_id from query path
	id := r.URL.Query().Get("room_id")
	// convert room_id param string to int
	roomId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "RoomHandler.DeleteRoom(): room_id should be an integer!")
		return
	}

	// delete room
	_, err = h.service.RoomService.Delete(roomId)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "RoomService.Delete(): error occured", err.Error())
		return
	}

	// return success message
	pkg.Response(w, "success")
}

func (h *Handlers) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	// get all rooms
	rooms := h.service.RoomService.GetAll()

	// return all rooms
	pkg.Response(w, rooms)
}

func (h *Handlers) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var roomParams models.Room

	// convert JSON to models.Room type
	err := json.NewDecoder(r.Body).Decode(&roomParams)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "RoomHandler.CreateRoom(): cannot convert JSON to models.Room struct", err.Error())
		return
	}

	// validate passed room data
	validator := NewRoomValidator(&roomParams)
	if validator.AllRoomFieldsValid != true {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "RoomHandler.CreateRoom(): Room data is not valid!", validator)
		return
	}

	// create room
	createdRoom, err := h.service.RoomService.Create(roomParams)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "RoomService.Create(): error occured", err.Error())
		return
	}

	// return created room
	pkg.Response(w, createdRoom)
}

func (h *Handlers) GetRoomById(w http.ResponseWriter, r *http.Request) {
	// get room_id from query path
	roomIdStr := r.URL.Query().Get("room_id")
	// convert room_id param string to int
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "RoomHandler.GetRoomById(): room_id should be an integer!", err.Error())
		return
	}

	// get Room from services
	room, err := h.service.RoomService.GetRoomById(roomId)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "RoomService.GetRoomById(): error occured", err.Error())
		return
	}

	// return found room
	pkg.Response(w, room)
}

func (h *Handlers) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	// get room_id from query path
	_ = r.URL.Query().Get("room_id")

	var roomParamsToUpdate models.Room
	// convert JSON to models.Room type
	err := json.NewDecoder(r.Body).Decode(&roomParamsToUpdate)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "RoomHandler.UpdateRoom(): cannot convert JSON to models.Room struct", err.Error())
		return
	}

	// validate passed room data
	validator := NewRoomValidator(&roomParamsToUpdate)
	if validator.AllRoomFieldsValid != true {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "RoomHandler.UpdateRoom(): Room data is not valid!", validator)
		return
	}

	// update room
	updatedRoom, err := h.service.RoomService.Update(roomParamsToUpdate)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "RoomService.Update(): error occured", err.Error())
		return
	}

	// return updated room
	pkg.Response(w, updatedRoom)
}

type RoomValidator struct {
	RoomToValidate     *models.Room      `json:"passed_room"`
	ValidationErrors   map[string]string `json:"validation_errors"`
	IsNumberValid      bool              `json:"is_number_valid"`
	IsCapacityValid    bool              `json:"is_capacity_valid"`
	AllRoomFieldsValid bool              `json:"all_room_fields_valid"`
}

func NewRoomValidator(room *models.Room) *RoomValidator {
	validationErrors := map[string]string{
		"number_error":   "Room.Number: should not be empty string",
		"capacity_error": "Room.Capacity: should not be negative integer or zero",
	}

	validator := &RoomValidator{RoomToValidate: room, ValidationErrors: validationErrors, AllRoomFieldsValid: false}
	validator.IsRoomValid()

	return validator
}

func (r *RoomValidator) IsRoomValid() {
	r.ValidateFields()

	if r.IsNumberValid && r.IsCapacityValid {
		r.AllRoomFieldsValid = true
	}
}

func (r *RoomValidator) ValidateFields() {
	if r.RoomToValidate.Number != "" {
		r.IsNumberValid = true
		delete(r.ValidationErrors, "number_error")
	}
	if r.RoomToValidate.Capacity > 0 {
		r.IsCapacityValid = true
		delete(r.ValidationErrors, "capacity_error")
	}
}
