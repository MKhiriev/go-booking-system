package handlers

import (
	"humoBooking/pkg"
	"log"
	"net/http"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS, PUT, HEAD, TRACE, CONNECT")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Expose-Headers", "*")
		if r.Method == "OPTIONS" {
			w.Write([]byte("OPTIONS"))
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RecoverAllPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				log.Println("Panic is processed!")
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) AuthorizationCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check destination address
		destination := r.URL
		// check if `Authorization` header exists
		authorizationHeader := r.Header.Get("Authorization")

		destinationPathIsAuthLogin := destination.Path == "/auth/login"
		destinationPathIsAuthRegister := destination.Path == "/auth/register"
		destinationPathIsAuthRefresh := destination.Path == "/auth/refresh"

		// if user want to refresh tokens => check if authorization header is not empty - success
		if destinationPathIsAuthRefresh {
			// check if autorization header is not empty => procceed to next.ServeHTTP(w, r)
			if authorizationHeader != "" {
				next.ServeHTTP(w, r)
			} else {
				log.Println("MiddleWare.AuthorizationCheck(`/auth/refresh`): error occured before token refreshment: `Authorization` header is empty")
				pkg.ErrorResponse(w, http.StatusBadRequest, "error occured before token refreshment", "`Authorization` header is empty")
				return
			}
		}

		// if user wants to log in or register and header `authorization` should be empty => procceed to next.ServeHTTP(w, r)
		if destinationPathIsAuthLogin || destinationPathIsAuthRegister {
			next.ServeHTTP(w, r)
			return
		}

		ipAddress := strings.Split(r.RemoteAddr, ":")[0]
		encodedAccessToken := strings.Split(authorizationHeader, " ")[1]

		// else check Access Token JWT in `Authorization` header is valid
		validator := h.service.AuthService.ValidateAccessToken(encodedAccessToken, ipAddress)
		if validator.ValidationError != nil {
			log.Println("AuthHandler.AuthorizationCheck(): validation of Access JWT token failed. Details: ", validator.ValidationError)
			pkg.ErrorResponse(w, http.StatusBadRequest, "validation of Access JWT token failed", validator.ValidationError.Error())
			return
		}

		subjectString := validator.AccessTokenClaims.Subject
		roleString := validator.AccessTokenClaims.Role

		var recordIdString string
		if r.URL.Query().Has("booking_id") {
			recordIdString = r.URL.Query().Get("booking_id")
		}
		if r.URL.Query().Has("room_id") {
			recordIdString = r.URL.Query().Get("room_id")
		}
		if r.URL.Query().Has("user_id") {
			recordIdString = r.URL.Query().Get("user_id")
		}

		var recordType string
		recordType = strings.Split(destination.Path, "/")[1]

		// check for permission to
		isAccessGranted, permissionCheckError := h.service.AuthService.CheckPermissions(destination.Path, recordType, recordIdString, subjectString, roleString)
		if permissionCheckError != nil {
			log.Println("AuthHandler.AuthorizationCheck(): error occurred during permission check. Details: ", permissionCheckError)
			pkg.ErrorResponse(w, http.StatusBadRequest, "error occurred during permission check", permissionCheckError.Error())
			return
		}

		if isAccessGranted != true {
			log.Println("AuthHandler.AuthorizationCheck(): access denied")
			pkg.ErrorResponse(w, http.StatusUnauthorized, "access denied")
			return
		}

		next.ServeHTTP(w, r)
	})
}
