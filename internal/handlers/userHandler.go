package handlers

import (
	"encoding/json"
	"humoBooking/pkg"
	"log"
	"net/http"
	"strconv"
)

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user_id")
	user_id, err := strconv.Atoi(id)
	if err != nil {
		pkg.ErrorResponse(w, 400, "только цифры")
		return
	}
	_, err = h.service.UserService.Delete(user_id)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, 500, "sorry")
		return
	}
	pkg.Response(w, "success")
}

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.service.UserService.GetAll()
	data, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		w.Write([]byte("на сервере возникла ошибка"))
		return
	}
	w.Write(data)
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) GetUserById(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
