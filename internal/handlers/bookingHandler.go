package handlers

import (
	"encoding/json"
	"humoBooking/pkg"
	"log"
	"net/http"
	"strconv"
)

func (h *Handlers) BookRoom(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	users := h.service.BookingService.GetAll()
	data, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		w.Write([]byte("на сервере возникла ошибка"))
		return
	}
	w.Write(data)
}

func (h *Handlers) GetBookingById(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) GetBookingsByRoomId(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) GetBookingsByRoomIdAndBookingTime(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) DeleteBookings(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("booking_id")
	bookingId, err := strconv.Atoi(id)
	if err != nil {
		pkg.ErrorResponse(w, 400, "только цифры")
		return
	}
	_, err = h.service.BookingService.Delete(bookingId)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, 500, "sorry")
		return
	}
	pkg.Response(w, "success")
}

func (h *Handlers) CheckIfRoomAvailable(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) GetOverlappingBookings(w http.ResponseWriter, r *http.Request) {

}
