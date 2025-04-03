package handlers

import (
	"encoding/json"
	"humoBooking/pkg"
	"log"
	"net/http"
	"strconv"
)

func (h *Handlers) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("room_id")
	roomId, err := strconv.Atoi(id)
	if err != nil {
		pkg.ErrorResponse(w, 400, "только цифры")
		return
	}
	_, err = h.service.RoomService.Delete(roomId)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, 500, "sorry")
		return
	}
	pkg.Response(w, "success")
}

func (h *Handlers) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	users := h.service.RoomService.GetAll()
	data, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		w.Write([]byte("на сервере возникла ошибка"))
		return
	}
	w.Write(data)
}

func (h *Handlers) CreateRoom(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) GetRoomById(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
