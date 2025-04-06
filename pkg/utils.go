package pkg

import (
	"encoding/json"
	"log"
	"net/http"
)

func ErrorResponse(writer http.ResponseWriter, statusCode int, message string, details ...interface{}) {
	writer.WriteHeader(statusCode)
	if len(details) == 0 {
		writer.Write([]byte(message))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"error":   message,
		"details": details[0],
	}

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Println(err)
		return
	}
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
