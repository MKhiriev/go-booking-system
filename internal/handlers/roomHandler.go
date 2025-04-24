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
	roomIdStr := r.URL.Query().Get("room_id")
	if roomIdStr == "" {
		log.Println("RoomHandler.DeleteRoom(): parameter `room_id` is empty or not passed")
		pkg.ErrorResponse(w, http.StatusBadRequest, "parameter `room_id` is empty or not passed")
		return
	}

	// convert room_id param string to int
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		log.Println("RoomHandler.DeleteRoom(): room_id should be an integer. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "room_id should be an integer")
		return
	}

	// delete room
	_, err = h.service.RoomService.Delete(roomId)
	if err != nil {
		log.Println("RoomHandler.DeleteRoom(): error occured during room deletion. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during room deletion", err.Error())
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
	subjectStr := r.Header.Get("subject")
	r.Header.Del("subject")
	subjectWhoCreatesRoom, conversionError := strconv.Atoi(subjectStr)
	if conversionError != nil {
		log.Println("RoomHandler.CreateRoom(): cannot convert `subject`-header to integer. Details: ", conversionError)
		pkg.ErrorResponse(w, http.StatusBadRequest, "cannot convert `subject`-header to integer", conversionError.Error())
		return
	}

	var roomParams models.Room

	// convert JSON to models.Room type
	err := json.NewDecoder(r.Body).Decode(&roomParams)
	if err != nil {
		log.Println("RoomHandler.CreateRoom(): cannot convert JSON to models.Room struct. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "cannot convert JSON to models.Room struct", err.Error())
		return
	}

	// validate passed room data
	validator := NewRoomValidator(&roomParams)
	if validator.AllRoomFieldsValid != true {
		log.Println("RoomHandler.CreateRoom(): Room data is not valid. Details: ", validator.ValidationErrors)
		pkg.ErrorResponse(w, http.StatusBadRequest, "Room data is not valid", validator.ValidationErrors)
		return
	}

	// to double-check if room id wasn't set
	roomParams.CreatedBy = subjectWhoCreatesRoom

	// create room
	createdRoom, err := h.service.RoomService.Create(roomParams)
	if err != nil {
		log.Println("RoomHandler.CreateRoom(): error occured during Room creation. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during Room creation", err.Error())
		return
	}

	// return created room
	pkg.Response(w, createdRoom)
}

func (h *Handlers) GetRoomById(w http.ResponseWriter, r *http.Request) {
	// get room_id from query path
	roomIdStr := r.URL.Query().Get("room_id")
	if roomIdStr == "" {
		log.Println("RoomHandler.GetRoomById(): parameter `room_id` is empty or not passed")
		pkg.ErrorResponse(w, http.StatusBadRequest, "parameter `room_id` is empty or not passed")
		return
	}

	// convert room_id param string to int
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		log.Println("RoomHandler.GetRoomById(): room_id should be an integer. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "room_id should be an integer", err.Error())
		return
	}

	// get Room from services
	room, err := h.service.RoomService.GetRoomById(roomId)
	if err != nil {
		log.Println("RoomHandler.GetRoomById(): error occured during getting room by id. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during getting room by id", err.Error())
		return
	}

	// return found room
	pkg.Response(w, room)
}

func (h *Handlers) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	// get room_id from query path
	roomIdStr := r.URL.Query().Get("room_id")
	if roomIdStr == "" {
		log.Println("RoomHandler.UpdateRoom(): parameter `room_id` is empty or not passed")
		pkg.ErrorResponse(w, http.StatusBadRequest, "parameter `room_id` is empty or not passed")
		return
	}

	// convert room_id param string to int
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		log.Println("RoomHandler.UpdateRoom(): room_id should be an integer. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "room_id should be an integer", err.Error())
		return
	}

	var roomParamsToUpdate models.Room
	// convert JSON to models.Room type
	err = json.NewDecoder(r.Body).Decode(&roomParamsToUpdate)
	if err != nil {
		log.Println("RoomHandler.UpdateRoom(): cannot convert JSON to models.Room struct. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "cannot convert JSON to models.Room struct", err.Error())
		return
	}

	// validate passed room data
	validator := NewRoomValidator(&roomParamsToUpdate)
	if validator.AllRoomFieldsValid != true {
		log.Println("RoomHandler.UpdateRoom(): Room data is not valid. Details: ", validator.ValidationErrors)
		pkg.ErrorResponse(w, http.StatusBadRequest, "Room data is not valid", validator.ValidationErrors)
		return
	}

	// to double-check if room id wasn't set
	roomParamsToUpdate.RoomId = roomId

	// update room
	updatedRoom, err := h.service.RoomService.Update(roomParamsToUpdate)
	if err != nil {
		log.Println("RoomHandler.UpdateRoom(): error occured during room update. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during room update", err.Error())
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
