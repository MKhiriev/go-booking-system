package handlers

import (
	"encoding/json"
	"humoBooking/internal/models"
	"humoBooking/pkg"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (h *Handlers) BookRoom(w http.ResponseWriter, r *http.Request) {
	var bookingParamsToCreate models.Booking
	// convert JSON to models.Booking type
	err := json.NewDecoder(r.Body).Decode(&bookingParamsToCreate)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "BookingHandler.BookRoom(): cannot convert JSON to models.Booking struct", err.Error())
		return
	}

	// validate passed booking data
	validator := NewBookingValidator(&bookingParamsToCreate)
	if validator.AllBookingFieldsValid != true {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "BookingHandler.BookRoom(): Booking data is not valid!", validator)
		return
	}

	// update booking
	updatedBooking, err := h.service.BookingService.BookRoom(
		bookingParamsToCreate.UserId,
		bookingParamsToCreate.RoomId,
		bookingParamsToCreate.DateTimeStart,
		bookingParamsToCreate.DateTimeEnd,
	)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "BookingService.BookRoom(): error occured", err.Error())
		return
	}

	// return updated booking
	pkg.Response(w, updatedBooking)
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
	// convert booking_id param string to int
	bookingId, err := strconv.Atoi(bookingIdStr)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "BookingHandler.GetBookingById(): booking_id should be an integer!", err.Error())
		return
	}

	// get Booking from services
	booking, err := h.service.BookingService.GetBookingById(bookingId)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "BookingService.GetBookingById(): error occured", err.Error())
		return
	}

	// return found booking
	pkg.Response(w, booking)
}

func (h *Handlers) GetBookingsByRoomId(w http.ResponseWriter, r *http.Request) {
	// get room_id from query path
	bookingIdStr := r.URL.Query().Get("room_id")
	// convert room_id param string to int
	bookingId, err := strconv.Atoi(bookingIdStr)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "BookingHandler.GetBookingsByRoomId(): room_id should be an integer!", err.Error())
		return
	}

	// get Booking slice from services
	bookings, err := h.service.BookingService.GetBookingsByRoomId(bookingId)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "BookingService.GetBookingsByRoomId(): error occured", err.Error())
		return
	}

	// return found booking
	pkg.Response(w, bookings)
}

func (h *Handlers) GetBookingsByRoomIdAndBookingTime(w http.ResponseWriter, r *http.Request) {
	// get room_id, datetime_start, datetime_end from query path
	roomIdStr := r.URL.Query().Get("room_id")
	dateTimeStartStr := r.URL.Query().Get("datetime_start")
	dateTimeEndStr := r.URL.Query().Get("datetime_end")

	// validate params
	validator := NewBookingQueryValidator(map[string]string{
		"room_id":        roomIdStr,
		"datetime_start": dateTimeStartStr,
		"datetime_end":   dateTimeEndStr,
	})

	if validator.AllBookingParamsValid == false {
		pkg.ErrorResponse(w, http.StatusBadRequest, "BookingHandler.GetBookingsByRoomIdAndBookingTime(): Booking data is not valid!", validator)
		return
	}

	// convert VALIDATED params: room_id to int, datetime_start and datetime_end to time.Time
	roomId, dateTimeStart, dateTimeEnd := convertRoomIdAndDateTimeRange(roomIdStr, dateTimeStartStr, dateTimeEndStr)

	// get Booking slice from services
	bookings, err := h.service.BookingService.GetBookingsByRoomIdAndBookingTime(roomId, dateTimeStart, dateTimeEnd)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "BookingService.GetBookingsByRoomId(): error occured", err.Error())
		return
	}

	// return found booking
	pkg.Response(w, bookings)
}

func (h *Handlers) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	// get booking_id from query path
	_ = r.URL.Query().Get("booking_id")

	var bookingParamsToUpdate models.Booking
	// convert JSON to models.Booking type
	err := json.NewDecoder(r.Body).Decode(&bookingParamsToUpdate)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "BookingHandler.UpdateBooking(): cannot convert JSON to models.Booking struct", err.Error())
		return
	}

	// validate passed booking data
	validator := NewBookingValidator(&bookingParamsToUpdate)
	if validator.AllBookingFieldsValid != true {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "BookingHandler.UpdateBooking(): Booking data is not valid!", validator)
		return
	}

	// update booking
	updatedBooking, err := h.service.BookingService.Update(bookingParamsToUpdate)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "BookingService.Update(): error occured", err.Error())
		return
	}

	// return updated booking
	pkg.Response(w, updatedBooking)
}

func (h *Handlers) DeleteBookings(w http.ResponseWriter, r *http.Request) {
	// get booking_id string from query
	id := r.URL.Query().Get("booking_id")
	// convert booking_id string to int
	bookingId, err := strconv.Atoi(id)
	if err != nil {
		pkg.ErrorResponse(w, http.StatusBadRequest, "BookingHandler.DeleteBookings(): Booking data is not valid!", err.Error())
		return
	}

	// delete Booking
	_, err = h.service.BookingService.Delete(bookingId)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "BookingService.Delete(): error occured", err.Error())
		return
	}

	pkg.Response(w, "success")
}

func (h *Handlers) CheckIfRoomAvailable(w http.ResponseWriter, r *http.Request) {
	// get room_id, datetime_start, datetime_end from query path
	roomIdStr := r.URL.Query().Get("room_id")
	dateTimeStartStr := r.URL.Query().Get("datetime_start")
	dateTimeEndStr := r.URL.Query().Get("datetime_end")

	// validate params
	validator := NewBookingQueryValidator(map[string]string{
		"room_id":        roomIdStr,
		"datetime_start": dateTimeStartStr,
		"datetime_end":   dateTimeEndStr,
	})

	if validator.AllBookingParamsValid == false {
		pkg.ErrorResponse(w, http.StatusBadRequest, "BookingHandler.GetOverlappingBookings(): Booking data is not valid!", validator)
		return
	}

	// convert VALIDATED params: room_id to int, datetime_start and datetime_end to time.Time
	roomId, dateTimeStart, dateTimeEnd := convertRoomIdAndDateTimeRange(roomIdStr, dateTimeStartStr, dateTimeEndStr)

	// get result is room available during given time frame
	available, err := h.service.BookingService.CheckIfRoomAvailable(roomId, dateTimeStart, dateTimeEnd)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "BookingService.GetBookingsByRoomIdAndBookingTime(): error occured", err.Error())
		return
	}

	// return found booking
	pkg.Response(w, available)
}

func (h *Handlers) GetOverlappingBookings(w http.ResponseWriter, r *http.Request) {
	// get room_id, datetime_start, datetime_end from query path
	roomIdStr := r.URL.Query().Get("room_id")
	dateTimeStartStr := r.URL.Query().Get("datetime_start")
	dateTimeEndStr := r.URL.Query().Get("datetime_end")

	// validate params
	validator := NewBookingQueryValidator(map[string]string{
		"room_id":        roomIdStr,
		"datetime_start": dateTimeStartStr,
		"datetime_end":   dateTimeEndStr,
	})

	if validator.AllBookingParamsValid == false {
		pkg.ErrorResponse(w, http.StatusBadRequest, "BookingHandler.GetOverlappingBookings(): Booking data is not valid!", validator)
		return
	}

	// convert VALIDATED params: room_id to int, datetime_start and datetime_end to time.Time
	roomId, dateTimeStart, dateTimeEnd := convertRoomIdAndDateTimeRange(roomIdStr, dateTimeStartStr, dateTimeEndStr)

	// get Booking slice from services
	bookings, err := h.service.BookingService.GetBookingsByRoomIdAndBookingTime(roomId, dateTimeStart, dateTimeEnd)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "BookingService.GetBookingsByRoomIdAndBookingTime(): error occured", err.Error())
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

type BookingQueryValidator struct {
	BookingQueryParams    map[string]string `json:"booking_query"`
	ValidationErrors      map[string]string `json:"validation_errors"`
	ParamsToValidate      map[string]bool   `json:"params_to_validate"`
	IsUserIdValid         bool              `json:"is_userid_valid"`
	IsRoomIdValid         bool              `json:"is_roleid_valid"`
	IsDateTimeStartValid  bool              `json:"is_datetime_start_valid"`
	IsDateTimeEndValid    bool              `json:"is_datetime_end_valid"`
	AllBookingParamsValid bool              `json:"all_booking_params_valid"`
}

func NewBookingQueryValidator(bookingQueryParams map[string]string, fieldsToValidate ...string) *BookingQueryValidator {
	validationErrors := map[string]string{
		"user_id":        "user_id should be integer greater than zero",
		"room_id":        "room_id should be integer greater than zero",
		"datetime_start": "datetime_start: wrong format or empty string",
		"datetime_end":   "datetime_end: wrong format or empty string",
	}
	allParamsToValidate := map[string]bool{
		"user_id":        false,
		"room_id":        false,
		"datetime_start": false,
		"datetime_end":   false,
	}

	validator := &BookingQueryValidator{
		BookingQueryParams:    bookingQueryParams,
		ValidationErrors:      validationErrors,
		ParamsToValidate:      allParamsToValidate,
		AllBookingParamsValid: false,
	}

	if len(fieldsToValidate) != 0 {
		newParamsToValidate := validator.deleteUnneccessaryParams(fieldsToValidate)
		validator.ParamsToValidate = newParamsToValidate
	}

	validator.IsBookingQueryValid()

	return validator
}

func (b *BookingQueryValidator) deleteUnneccessaryParams(fieldsToCheck []string) map[string]bool {
	result := make(map[string]bool, 4)
	allNeccessaryFields := b.ParamsToValidate
	// check if in LIST OF INPUT FIELDS only neccessary fields are present
	// iterate over fields to check
	for _, fieldToCheck := range fieldsToCheck {
		// for each field to check
		// iterate over all neccessary fields
		for neccessaryField := range allNeccessaryFields {
			// if fieldToCheck == necessaryField => append fieldToCheck to result and go to the next iteration of inner for loop
			// else do nothing
			if fieldToCheck == neccessaryField {
				result[fieldToCheck] = false
				delete(b.ValidationErrors, fieldToCheck)
				continue
			}
		}
	}

	return result
}

func (b *BookingQueryValidator) IsBookingQueryValid() {
	areValid := true
	b.ValidateParams()

	for _, isParamValid := range b.ParamsToValidate {
		areValid = areValid && isParamValid
	}

	b.AllBookingParamsValid = areValid
}

func (b *BookingQueryValidator) ValidateParams() {
	layout := "2006-01-02 15:04:05 Z0700"
	for field := range b.ParamsToValidate {
		if field == "user_id" {
			userId, userIdErr := strconv.Atoi(b.BookingQueryParams["user_id"])

			if userIdErr == nil && userId > 0 {
				b.IsUserIdValid = true
				b.ParamsToValidate[field] = true
				delete(b.ValidationErrors, "roomid_error")
			}
		} else if field == "room_id" {
			roomId, roomIdErr := strconv.Atoi(b.BookingQueryParams["room_id"])

			if roomIdErr == nil && roomId > 0 {
				b.IsRoomIdValid = true
				b.ParamsToValidate[field] = true
				delete(b.ValidationErrors, "room_id")
			}
		} else if field == "datetime_start" {
			dateTimeStart, dateTimeStartErr := time.Parse(layout, b.BookingQueryParams["datetime_start"])

			if dateTimeStartErr == nil && !dateTimeStart.IsZero() {
				b.IsUserIdValid = true
				b.ParamsToValidate[field] = true
				delete(b.ValidationErrors, "datetime_start_error")
			}
		} else if field == "datetime_end" {
			dateTimeEnd, dateTimeEndErr := time.Parse(layout, b.BookingQueryParams["datetime_end"])

			if dateTimeEndErr == nil && !dateTimeEnd.IsZero() {
				b.IsDateTimeEndValid = true
				b.ParamsToValidate[field] = true
				delete(b.ValidationErrors, "datetime_end_error")
			}
		}
	}
}

func convertRoomIdAndDateTimeRange(roomIdStr string, dateTimeStartStr string, dateTimeEndStr string) (int, time.Time, time.Time) {
	layout := "2006-01-02 15:04:05 Z0700"
	roomId, _ := strconv.Atoi(roomIdStr)
	dateTimeStart, _ := time.Parse(layout, dateTimeStartStr)
	dateTimeEnd, _ := time.Parse(layout, dateTimeEndStr)

	return roomId, dateTimeStart, dateTimeEnd
}
