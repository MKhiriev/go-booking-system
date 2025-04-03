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
	router.HandleFunc("/users", h.GetAllUsers).Methods(http.MethodGet, http.MethodOptions)
	//router.HandleFunc("/user-by-id", h.GetUserById).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/drop", h.DeleteUser).Methods("DELETE", "OPTION")
	router.HandleFunc("/create", h.CreateUser).Methods(http.MethodPost, http.MethodOptions)
	return router
}
