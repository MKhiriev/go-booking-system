package handlers

import (
	"github.com/gorilla/mux"
	"humoBooking/internal/services"
	"net/http"
)

type Handlers struct {
	service *services.Service
}

func NewHandler(s *services.Service) *Handlers {
	return &Handlers{service: s}
}

func (h *Handlers) Init() *mux.Router {
	router := mux.NewRouter()
	router.Use(CORS, RecoverAllPanic)

	// User Handler
	router.HandleFunc("/user/create", h.CreateUser).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/user/all", h.GetAllUsers).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/user", h.GetUserById).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/user/update", h.UpdateUser).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/user/drop", h.DeleteUser).Methods(http.MethodDelete, http.MethodOptions)

	// Room Handler
	router.HandleFunc("/room/create", h.CreateRoom).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/room/all", h.GetAllRooms).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/room", h.GetRoomById).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/room/update", h.UpdateRoom).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/room/drop", h.DeleteRoom).Methods(http.MethodDelete, http.MethodOptions)

	// Booking Handler
	router.HandleFunc("/booking/all", h.GetAllBookings).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/booking", h.GetBookingById).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/booking/room", h.GetBookingsByRoomId).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/booking/room_time", h.GetBookingsByRoomIdAndBookingTime).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/booking/drop", h.DeleteBookings).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/booking/available/room", h.CheckIfRoomAvailable).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/booking/create", h.BookRoom).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/booking/overlapping", h.GetOverlappingBookings).Methods(http.MethodPost, http.MethodOptions)

	return router
}
