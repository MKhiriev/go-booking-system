package handlers

import (
	"encoding/json"
	"humoBooking/internal/models"
	"humoBooking/pkg"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *Handlers) BookRoom(w http.ResponseWriter, r *http.Request) {
	subjectStr := r.Header.Get("subject")
	r.Header.Del("subject")
	subjectWhoCreatesBooking, conversionError := strconv.Atoi(subjectStr)
	if conversionError != nil {
		log.Println("BookingHandler.BookRoom(): cannot convert `subject`-header to integer. Details: ", conversionError)
		pkg.ErrorResponse(w, http.StatusBadRequest, "cannot convert `subject`-header to integer", conversionError.Error())
		return
	}

	var bookingParamsToCreate models.Booking
	// convert JSON to models.Booking type
	err := json.NewDecoder(r.Body).Decode(&bookingParamsToCreate)
	if err != nil {
		log.Println("BookingHandler.BookRoom(): cannot convert JSON to models.Booking struct. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "cannot convert JSON to models.Booking struct", err.Error())
		return
	}

	// validate passed booking data
	validator := NewBookingValidator(&bookingParamsToCreate)
	if validator.AllBookingFieldsValid != true {
		log.Println("BookingHandler.BookRoom(): booking data is not valid. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "booking data is not valid", validator)
		return
	}

	// create booking
	createdBooking, err := h.service.BookingService.BookRoom(
		bookingParamsToCreate.UserId,
		bookingParamsToCreate.RoomId,
		bookingParamsToCreate.DateTimeStart,
		bookingParamsToCreate.DateTimeEnd,
		subjectWhoCreatesBooking,
	)
	if err != nil {
		log.Println("BookingHandler.BookRoom(): error occured during Room Booking. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during Room Booking", err.Error())
		return
	}

	// return created booking
	pkg.Response(w, createdBooking)
}

func (h *Handlers) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	// get all bookings
	bookings := h.service.BookingService.GetAll()

	// return all bookings
	pkg.Response(w, bookings)
}

func (h *Handlers) GetBookingById(w http.ResponseWriter, r *http.Request) {
	// get booking_id from query path
	bookingIdStr := r.URL.Query().Get("booking_id")
	if bookingIdStr == "" {
		log.Println("BookingHandler.GetBookingById(): parameter `booking_id` is empty or not passed")
		pkg.ErrorResponse(w, http.StatusBadRequest, "parameter `booking_id` is empty or not passed")
		return
	}
	// convert booking_id param string to int
	bookingId, err := strconv.Atoi(bookingIdStr)
	if err != nil {
		log.Println("BookingHandler.GetBookingById(): booking_id should be an integer. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "booking_id should be an integer", err.Error())
		return
	}

	// get Booking from services
	booking, err := h.service.BookingService.GetBookingById(bookingId)
	if err != nil {
		log.Println("BookingHandler.GetBookingById(): error occured during getting booking by id. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during getting booking by id", err.Error())
		return
	}

	// return found booking
	pkg.Response(w, booking)
}

func (h *Handlers) GetBookingsByRoomId(w http.ResponseWriter, r *http.Request) {
	// get room_id from query path
	roomIdStr := r.URL.Query().Get("room_id")
	if roomIdStr == "" {
		log.Println("BookingHandler.GetBookingsByRoomId(): parameter `room_id` is empty or not passed")
		pkg.ErrorResponse(w, http.StatusBadRequest, "parameter `room_id` is empty or not passed")
		return
	}

	// convert room_id param string to int
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		log.Println("BookingHandler.GetBookingsByRoomId(): room_id should be an integer. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "room_id should be an integer", err.Error())
		return
	}

	// get Booking slice from services
	foundBooking, err := h.service.BookingService.GetBookingsByRoomId(roomId)
	if err != nil {
		log.Println("BookingHandler.GetBookingsByRoomId(): error occured during getting booking by room id. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during getting booking by room id", err.Error())
		return
	}

	// return found booking
	pkg.Response(w, foundBooking)
}

func (h *Handlers) GetBookingsByRoomIdAndBookingTime(w http.ResponseWriter, r *http.Request) {
	// validate params
	validator := NewBookingQueryParamsValidator(r.URL.RawQuery)

	if validator.AllQueryParamsValid == false {
		log.Println("BookingHandler.GetBookingsByRoomIdAndBookingTime(): booking query is not valid. Details: ", validator.ValidationErrors)
		pkg.ErrorResponse(w, http.StatusBadRequest, "booking query is not valid", validator.ValidationErrors)
		return
	}

	roomId, dateTimeStart, dateTimeEnd := validator.RoomId, validator.DateTimeStart, validator.DateTimeEnd
	// get Booking slice from services
	bookings, err := h.service.BookingService.GetBookingsByRoomIdAndBookingTime(roomId, dateTimeStart, dateTimeEnd)
	if err != nil {
		log.Println("BookingHandler.GetBookingsByRoomIdAndBookingTime(): error occured during getting bookings by room id and booking time. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during getting bookings by room id and booking time", err.Error())
		return
	}

	// return found booking
	pkg.Response(w, bookings)
}

func (h *Handlers) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	// get booking_id from query path
	bookingIdStr := r.URL.Query().Get("booking_id")
	if bookingIdStr == "" {
		log.Println("BookingHandler.UpdateBooking(): parameter `booking_id` is empty or not passed")
		pkg.ErrorResponse(w, http.StatusBadRequest, "parameter `booking_id` is empty or not passed")
		return
	}

	// convert booking_id param string to int
	bookingId, err := strconv.Atoi(bookingIdStr)
	if err != nil {
		log.Println("BookingHandler.UpdateBooking(): booking_id should be an integer. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "booking_id should be an integer", err.Error())
		return
	}

	var bookingParamsToUpdate models.Booking
	// convert JSON to models.Booking type
	jsonConversionErr := json.NewDecoder(r.Body).Decode(&bookingParamsToUpdate)
	if jsonConversionErr != nil {
		log.Println("BookingHandler.UpdateBooking(): cannot convert JSON to models.Booking struct. Details: ", jsonConversionErr)
		pkg.ErrorResponse(w, http.StatusBadRequest, "cannot convert JSON to models.Booking struct", jsonConversionErr.Error())
		return
	}

	// validate passed booking data
	validator := NewBookingValidator(&bookingParamsToUpdate)
	if validator.AllBookingFieldsValid != true {
		log.Println("BookingHandler.UpdateBooking(): Booking data is not valid. Details: ", validator.ValidationErrors)
		pkg.ErrorResponse(w, http.StatusBadRequest, "Booking data is not valid", validator.ValidationErrors)
		return
	}

	// to double-check if booking id wasn't set
	bookingParamsToUpdate.BookingId = bookingId

	// update booking
	updatedBooking, err := h.service.BookingService.Update(bookingParamsToUpdate)
	if err != nil {
		log.Println("BookingHandler.UpdateBooking(): error occured during booking update. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during booking update", err.Error())
		return
	}

	// return updated booking
	pkg.Response(w, updatedBooking)
}

func (h *Handlers) DeleteBookings(w http.ResponseWriter, r *http.Request) {
	// get booking_id string from query
	bookingIdStr := r.URL.Query().Get("booking_id")
	if bookingIdStr == "" {
		log.Println("BookingHandler.DeleteBookings(): parameter `booking_id` is empty or not passed")
		pkg.ErrorResponse(w, http.StatusBadRequest, "parameter `booking_id` is empty or not passed")
		return
	}

	// convert booking_id string to int
	bookingId, err := strconv.Atoi(bookingIdStr)
	if err != nil {
		log.Println("BookingHandler.DeleteBookings(): booking data is not valid. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "booking data is not valid", err.Error())
		return
	}

	// delete Booking
	_, err = h.service.BookingService.Delete(bookingId)
	if err != nil {
		log.Println("BookingHandler.DeleteBookings(): error occured during booking deletion. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during booking deletion", err.Error())
		return
	}

	pkg.Response(w, "success")
}

func (h *Handlers) CheckIfRoomAvailable(w http.ResponseWriter, r *http.Request) {
	// validate params
	validator := NewBookingQueryParamsValidator(r.URL.RawQuery)

	if validator.AllQueryParamsValid == false {
		log.Println("BookingHandler.CheckIfRoomAvailable(): booking query is not valid. Details: ", validator.ValidationErrors)
		pkg.ErrorResponse(w, http.StatusBadRequest, "booking query is not valid", validator.ValidationErrors)
		return
	}

	// convert VALIDATED params: room_id to int, datetime_start and datetime_end to time.Time
	roomId, dateTimeStart, dateTimeEnd := validator.RoomId, validator.DateTimeStart, validator.DateTimeEnd

	// get result is room available during given time frame
	available, err := h.service.BookingService.CheckIfRoomAvailable(roomId, dateTimeStart, dateTimeEnd)
	if err != nil {
		log.Println("BookingHandler.CheckIfRoomAvailable(): error occured during room availability check. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during room availability check", err.Error())
		return
	}

	// return found booking
	pkg.Response(w, available)
}

func (h *Handlers) GetOverlappingBookings(w http.ResponseWriter, r *http.Request) {
	// validate params
	validator := NewBookingQueryParamsValidator(r.URL.RawQuery)

	if validator.AllQueryParamsValid == false {
		log.Println("BookingHandler.GetOverlappingBookings(): booking query is not valid. Details: ", validator.ValidationErrors)
		pkg.ErrorResponse(w, http.StatusBadRequest, "booking query is not valid", validator.ValidationErrors)
		return
	}

	// convert VALIDATED params: room_id to int, datetime_start and datetime_end to time.Time
	roomId, dateTimeStart, dateTimeEnd := validator.RoomId, validator.DateTimeStart, validator.DateTimeEnd

	// get Booking slice from services
	bookings, err := h.service.BookingService.GetBookingsByRoomIdAndBookingTime(roomId, dateTimeStart, dateTimeEnd)
	if err != nil {
		log.Println("BookingHandler.GetOverlappingBookings(): error occured during getting overlapping bookings. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during getting overlapping bookings", err.Error())
		return
	}

	// return found booking
	pkg.Response(w, bookings)
}

type BookingValidator struct {
	BookingToValidate     *models.Booking   `json:"passed_booking"`
	ValidationErrors      map[string]string `json:"validation_errors"`
	IsUserIdValid         bool              `json:"is_name_valid"`
	IsRoomIdValid         bool              `json:"is_roleid_valid"`
	IsDateTimeStartValid  bool              `json:"is_datetime_start_valid"`
	IsDateTimeEndValid    bool              `json:"is_datetime_end_valid"`
	AllBookingFieldsValid bool              `json:"all_booking_fields_valid"`
}

func NewBookingValidator(booking *models.Booking) *BookingValidator {
	validationErrors := map[string]string{
		"userid_error":         "Booking.UserId: should not be negative integer or zero",
		"roomid_error":         "Booking.RoomId: should not be negative integer or zero",
		"datetime_start_error": "Booking.DateTimeStart: should not be null time value",
		"datetime_end_error":   "Booking.DateTimeEnd: should not be null time value",
	}

	validator := &BookingValidator{BookingToValidate: booking, ValidationErrors: validationErrors, AllBookingFieldsValid: false}
	validator.IsBookingValid()

	return validator
}

func (b *BookingValidator) IsBookingValid() {
	b.ValidateFields()

	if b.IsUserIdValid && b.IsRoomIdValid && b.IsDateTimeStartValid && b.IsDateTimeEndValid {
		b.AllBookingFieldsValid = true
	}
}

func (b *BookingValidator) ValidateFields() {
	if b.BookingToValidate.RoomId > 0 {
		b.IsRoomIdValid = true
		delete(b.ValidationErrors, "roomid_error")
	}
	if b.BookingToValidate.UserId > 0 {
		b.IsUserIdValid = true
		delete(b.ValidationErrors, "userid_error")
	}
	if !b.BookingToValidate.DateTimeStart.IsZero() {
		b.IsDateTimeStartValid = true
		delete(b.ValidationErrors, "datetime_start_error")
	}
	if !b.BookingToValidate.DateTimeEnd.IsZero() {
		b.IsDateTimeEndValid = true
		delete(b.ValidationErrors, "datetime_end_error")
	}
}

type BookingQueryParamsValidator struct {
	RawQueryParamsString string            `json:"raw_query_params_string"`
	QueryParams          map[string]string `json:"query_params"`

	RoomId        int       `json:"room_id"`
	DateTimeStart time.Time `json:"datetime_start"`
	DateTimeEnd   time.Time `json:"datetime_end"`

	ValidationErrors    map[string]string `json:"validation_errors"`
	AllQueryParamsValid bool              `json:"all_login_fields_valid"`
}

func NewBookingQueryParamsValidator(rawQueryParamsString string) *BookingQueryParamsValidator {
	validator := &BookingQueryParamsValidator{RawQueryParamsString: rawQueryParamsString, QueryParams: make(map[string]string), AllQueryParamsValid: false}
	validator.AreParamsValid()

	return validator
}

func (b *BookingQueryParamsValidator) AreParamsValid() {
	b.ValidateParams()

	if len(b.ValidationErrors) == 0 {
		b.AllQueryParamsValid = true
	}
}

func (b *BookingQueryParamsValidator) ValidateParams() {
	if b.RawQueryParamsString == "" {
		b.ValidationErrors["empty_params_query_error"] = "no parameters have been passed"
		return
	}

	queryParams := strings.Split(b.RawQueryParamsString, "&")
	for _, queryParam := range queryParams {
		paramValuePair := strings.Split(queryParam, "=")
		param := paramValuePair[0]
		value := paramValuePair[1]
		b.QueryParams[param] = value
	}

	if roomIdStr, ok := b.QueryParams["room_id"]; ok == true {
		roomId, conversionError := strconv.Atoi(roomIdStr)
		if conversionError != nil || roomId < 1 {
			b.ValidationErrors["roomId_error"] = "should not be negative integer or zero"
		}
		b.RoomId = roomId
	}

	emptyTime := time.Time{}
	if dateTimeStartStr, ok := b.QueryParams["datetime_start"]; ok == true {
		dateTimeStart, timeConversionError := time.Parse(time.RFC3339, dateTimeStartStr)
		if timeConversionError != nil || dateTimeStart == emptyTime {
			b.ValidationErrors["datetime_start"] = "wrong DateTime format"
		}
		b.DateTimeStart = dateTimeStart
	}

	if dateTimeEndStr, ok := b.QueryParams["datetime_end"]; ok == true {
		dateTimeEnd, timeConversionError := time.Parse(time.RFC3339, dateTimeEndStr)
		if timeConversionError != nil || dateTimeEnd == emptyTime {
			b.ValidationErrors["datetime_end"] = "wrong DateTime format"
		}
		b.DateTimeEnd = dateTimeEnd
	}
}
