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

	// Auth Handler
	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", h.Register).Methods(http.MethodPost, http.MethodOptions)
	auth.HandleFunc("/login", h.Login).Methods(http.MethodPost, http.MethodOptions)
	auth.HandleFunc("/refresh", h.RefreshToken).Methods(http.MethodPost, http.MethodOptions)

	// User Handler
	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/all", h.GetAllUsers).Methods(http.MethodGet, http.MethodOptions)
	user.HandleFunc("/", h.GetUserById).Methods(http.MethodGet, http.MethodOptions)
	user.HandleFunc("/update", h.UpdateUser).Methods(http.MethodPost, http.MethodOptions)
	user.HandleFunc("/drop", h.DeleteUser).Methods(http.MethodDelete, http.MethodOptions)

	// Room Handler
	room := router.PathPrefix("/room").Subrouter()
	room.HandleFunc("/create", h.CreateRoom).Methods(http.MethodPost, http.MethodOptions)
	room.HandleFunc("/all", h.GetAllRooms).Methods(http.MethodGet, http.MethodOptions)
	room.HandleFunc("/", h.GetRoomById).Methods(http.MethodGet, http.MethodOptions)
	room.HandleFunc("/update", h.UpdateRoom).Methods(http.MethodPost, http.MethodOptions)
	room.HandleFunc("/drop", h.DeleteRoom).Methods(http.MethodDelete, http.MethodOptions)

	// Booking Handler
	booking := router.PathPrefix("/booking").Subrouter()
	booking.HandleFunc("/all", h.GetAllBookings).Methods(http.MethodGet, http.MethodOptions)
	booking.HandleFunc("/", h.GetBookingById).Methods(http.MethodGet, http.MethodOptions)
	booking.HandleFunc("/room", h.GetBookingsByRoomId).Methods(http.MethodGet, http.MethodOptions)
	booking.HandleFunc("/room_time", h.GetBookingsByRoomIdAndBookingTime).Methods(http.MethodGet, http.MethodOptions)
	booking.HandleFunc("/drop", h.DeleteBookings).Methods(http.MethodDelete, http.MethodOptions)
	booking.HandleFunc("/available/room", h.CheckIfRoomAvailable).Methods(http.MethodGet, http.MethodOptions)
	booking.HandleFunc("/create", h.BookRoom).Methods(http.MethodPost, http.MethodOptions)
	booking.HandleFunc("/overlapping", h.GetOverlappingBookings).Methods(http.MethodGet, http.MethodOptions)
	booking.HandleFunc("/update", h.UpdateBooking).Methods(http.MethodPatch, http.MethodOptions)

	return router
}
