package server

import (
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func (s *Server) ServerRun(handler http.Handler, port string) error {
	s.server = &http.Server{
		Addr:         "localhost:" + port,
		Handler:      nil,
		ReadTimeout:  time.Second * 10, // 10 секунд максимум на чтение
		WriteTimeout: time.Second * 10, // 10 секунда на отправку запроса
	}

	return s.server.ListenAndServe()
}
