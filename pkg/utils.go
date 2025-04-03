package pkg

import (
	"encoding/json"
	"log"
	"net/http"
)

func ErrorResponse(writer http.ResponseWriter, statusCode int, text string) {
	writer.WriteHeader(statusCode)
	writer.Write([]byte(text))
}

func Response(writer http.ResponseWriter, data any) {
	dataToWrite, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(dataToWrite)
}
