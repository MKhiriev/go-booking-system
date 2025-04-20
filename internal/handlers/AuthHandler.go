package handlers

import (
	"encoding/json"
	"humoBooking/internal/models"
	"humoBooking/internal/services"
	"humoBooking/pkg"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistrationParams struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
	RoleId    int    `json:"role_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type EncodedRefreshJWTToken struct {
	EncodedRefreshToken string `json:"refresh_token"`
}

type JWTTokens struct {
	AccessToken  pkg.JWTToken `json:"token"`
	RefreshToken pkg.JWTToken `json:"refresh_token"`
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	loginParams := LoginParams{}
	decodingJSONError := json.NewDecoder(r.Body).Decode(&loginParams)
	if decodingJSONError != nil {
		log.Println("AuthHandler.Login(): error occured during decoding JSON. Details: ", decodingJSONError.Error())
		pkg.ErrorResponse(w, http.StatusBadRequest, "error occured during decoding JSON", decodingJSONError.Error())
		return
	}

	validator := NewLogingParamsValidator(&loginParams)
	if validator.AllLoginParamsFieldsValid != true {
		log.Println("AuthHandler.Login(): login data is not valid!", validator.ValidationErrors)
		pkg.ErrorResponse(w, http.StatusBadRequest, "login data is not valid!", validator)
		return
	}

	// Identification & Authentication
	foundUser, loginError := h.service.AuthService.CheckIfUserExistsAndPasswordIsCorrect(loginParams.Username, loginParams.Password)
	if loginError != nil {
		log.Println("AuthHandler.Login(): error occured during login")
		pkg.ErrorResponse(w, http.StatusBadRequest, "error occured during login", loginError.Error())
		return
	}

	// get the identity of who sent the token
	identity := pkg.IPAddressIdentity{
		IP: strings.Split(r.RemoteAddr, ":")[0],
	}

	token, refreshToken := h.service.AuthService.GenerateTokens(foundUser, identity)
	JWTtokens := JWTTokens{AccessToken: token, RefreshToken: refreshToken}

	pkg.Response(w, JWTtokens)
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	registrationParams := RegistrationParams{}
	decodingJSONError := json.NewDecoder(r.Body).Decode(&registrationParams)
	if decodingJSONError != nil {
		log.Println("AuthHandler.Login(): error occured during decoding JSON. Details: ", decodingJSONError.Error())
		pkg.ErrorResponse(w, http.StatusBadRequest, "error occured during decoding JSON", decodingJSONError.Error())
		return
	}

	userData := models.User{
		Name:      registrationParams.Name,
		Email:     registrationParams.Email,
		Telephone: registrationParams.Telephone,
		RoleId:    registrationParams.RoleId,
		UserName:  registrationParams.Username,
		Password:  registrationParams.Password,
		Active:    true,
	}
	userValidator := NewUserValidator(&userData)
	if userValidator.AllUserFieldsValid != true {
		log.Println("AuthHandler.Register(): User data is not valid!", userValidator.ValidationErrors)
		pkg.ErrorResponse(w, http.StatusBadRequest, "User data is not valid!", userValidator.ValidationErrors)
		return
	}

	usernameAndPasswordParams := LoginParams{
		Username: registrationParams.Username,
		Password: registrationParams.Password,
	}
	usernameAndPasswordValidator := NewLogingParamsValidator(&usernameAndPasswordParams)
	if usernameAndPasswordValidator.AllLoginParamsFieldsValid != true {
		log.Println("AuthHandler.Register(): Username or Password data is not valid!", usernameAndPasswordValidator.ValidationErrors)
		pkg.ErrorResponse(w, http.StatusBadRequest, "Username or Password data is not valid", usernameAndPasswordValidator.ValidationErrors)
		return
	}

	user, err := h.service.AuthService.Create(userData)
	if err != nil {
		log.Println("AuthHandler.Register(): error occured during User creation", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during User creation", err.Error())
		return
	}

	pkg.Response(w, user)
}

// RefreshToken Client should send JWTTokens!
func (h *Handlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		log.Println("AuthHandler.RefreshToken(): empty authorization header.")
		pkg.ErrorResponse(w, http.StatusBadRequest, "empty authorization header")
		return
	}

	authorizationHeaderParts := strings.Split(authorizationHeader, " ")
	if !strings.Contains(authorizationHeader, "Bearer") || len(authorizationHeaderParts) != 2 {
		log.Println("AuthHandler.RefreshToken(): wrong authorization header format. Reason: no 'Bearer' or no access token. Authorization header: ", authorizationHeader)
		pkg.ErrorResponse(w, http.StatusBadRequest, "wrong authorization header format. Reason: no 'Bearer' or no access token", authorizationHeader)
		return
	}

	accessToken := authorizationHeaderParts[1]
	refreshTokenJSON := EncodedRefreshJWTToken{}
	decodingJSONError := json.NewDecoder(r.Body).Decode(&refreshTokenJSON)
	if decodingJSONError != nil {
		log.Println("AuthHandler.RefreshToken(): error occured during decoding JSON. Details: ", decodingJSONError.Error())
		pkg.ErrorResponse(w, http.StatusBadRequest, "error occured during decoding JSON", decodingJSONError.Error())
		return
	}

	// prepare data for token validation
	ipAddress := strings.Split(r.RemoteAddr, ":")[0]
	refreshToken := refreshTokenJSON.EncodedRefreshToken

	// Validate access token
	accessTokenValidator := h.service.AuthService.ValidateAccessToken(accessToken, ipAddress)
	if accessTokenValidator.ValidationError != nil && accessTokenValidator.ValidationError.Error() != services.IsExpired {
		log.Println("AuthHandler.RefreshToken(): error occured during validation of JWT Access Token. Details: ", accessTokenValidator.ValidationError)
		pkg.ErrorResponse(w, http.StatusBadRequest, "error occured during validation of JWT Access Token", accessTokenValidator.ValidationError.Error())
		return
	}

	// Validate refresh token
	refreshTokenValidator := h.service.AuthService.ValidateRefreshToken(refreshToken, ipAddress)
	if refreshTokenValidator.ValidationError != nil {
		log.Println("AuthHandler.RefreshToken(): error occured during validation of JWT Refresh Token. Details: ", refreshTokenValidator.ValidationError)
		pkg.ErrorResponse(w, http.StatusBadRequest, "error occured during validation of JWT Refresh Token", refreshTokenValidator.ValidationError.Error())
		return
	}

	// check if refresh token is expired
	if refreshTokenValidator.IsExpired == true {
		log.Println("AuthHandler.RefreshToken(): JWT Refresh Token is expired. Details: ", refreshTokenValidator.ValidationError)
		pkg.ErrorResponse(w, http.StatusBadRequest, "JWT Refresh Token is expired", refreshTokenValidator.ValidationError.Error())
		return
	}

	// if tokens are assigned to different users - deny
	if accessTokenValidator.AccessTokenClaims.Subject != refreshTokenValidator.RefreshTokenClaims.Subject {
		log.Println("AuthHandler.RefreshToken(): tokens are assigned to different users")
		pkg.ErrorResponse(w, http.StatusForbidden, "tokens are assigned to different users", refreshTokenValidator.ValidationError.Error())
		return
	}

	userId, stringToIntConversionError := strconv.Atoi(refreshTokenValidator.RefreshTokenClaims.Subject)
	if stringToIntConversionError != nil {
		log.Println("AuthHandler.RefreshToken(): cannot convert Subject claim to Integer. Details: ", stringToIntConversionError)
		pkg.ErrorResponse(w, http.StatusForbidden, "cannot convert Subject claim to Integer", stringToIntConversionError.Error())
		return
	}

	// get User to generate new set of tokens
	user, userByIdError := h.service.UserService.GetUserById(userId)
	if userByIdError != nil {
		log.Println("AuthHandler.RefreshToken(): error occured during user search by id. Details: ", userByIdError)
		pkg.ErrorResponse(w, http.StatusForbidden, "error occured during user search by id", userByIdError.Error())
		return
	}

	// if all tokens are valid - generate a new pair of tokens
	identity := pkg.IPAddressIdentity{IP: ipAddress}

	newAccessToken, newRefreshToken := h.service.AuthService.GenerateTokens(user, identity)
	JWTtokens := JWTTokens{AccessToken: newAccessToken, RefreshToken: newRefreshToken}

	pkg.Response(w, JWTtokens)
}

type LoginParamsValidator struct {
	LoginParamsToValidate     *LoginParams      `json:"passed_login_params"`
	ValidationErrors          map[string]string `json:"validation_errors"`
	IsUsernameValid           bool              `json:"is_username_valid"`
	IsPasswordValid           bool              `json:"is_password_valid"`
	AllLoginParamsFieldsValid bool              `json:"all_login_fields_valid"`
}

func NewLogingParamsValidator(loginParams *LoginParams) *LoginParamsValidator {
	validationErrors := map[string]string{
		"username_error": "LogingParams.Username: username should not be an empty string",
		"password_error": "LogingParams.Password: password should not be an empty string",
	}

	validator := &LoginParamsValidator{LoginParamsToValidate: loginParams, ValidationErrors: validationErrors, AllLoginParamsFieldsValid: false}
	validator.IsLoginValid()

	return validator
}

func (l *LoginParamsValidator) IsLoginValid() {
	l.ValidateFields()

	if l.IsUsernameValid && l.IsPasswordValid {
		l.AllLoginParamsFieldsValid = true
	}
}

func (l *LoginParamsValidator) ValidateFields() {
	if l.LoginParamsToValidate.Username != "" {
		l.IsUsernameValid = true
		delete(l.ValidationErrors, "username_error")
	}
	if l.LoginParamsToValidate.Password != "" {
		l.IsPasswordValid = true
		delete(l.ValidationErrors, "password_error")
	}
}
