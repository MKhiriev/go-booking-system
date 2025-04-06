package handlers

import (
	"encoding/json"
	"humoBooking/internal/models"
	"humoBooking/pkg"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userParams models.User
	userDecoder := json.NewDecoder(r.Body)
	err := userDecoder.Decode(&userParams)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// validate passed user data
	validator := NewUserValidator(&userParams)
	if validator.AllUserFieldsValid != true {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "Wrong user data provided", validator)
		return
	}

	createdUser, err := h.service.UserService.Create(userParams)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "Error occured during User creation")
		return
	}
	pkg.Response(w, createdUser)
}

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.service.UserService.GetAll()

	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Response(w, usersJSON)
}

func (h *Handlers) GetUserById(w http.ResponseWriter, r *http.Request) {
	// get user_id from path
	userIdStr := r.URL.Query().Get("user_id")
	// convert string to int
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, err.Error(), "user_id is not an integer!")
	}

	// get User from services
	user, err := h.service.UserService.GetUserById(userId)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, err.Error(), "UserService.GetUserById(): error occured")
	}

	// return found user
	pkg.Response(w, user)
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var userParamsToUpdate models.User
	// convert JSON to models.User type
	err := json.NewDecoder(r.Body).Decode(&userParamsToUpdate)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// validate passed user data
	validator := NewUserValidator(&userParamsToUpdate)
	if validator.AllUserFieldsValid != true {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "Wrong user data provided", validator)
		return
	}

	// update user
	updatedUser, err := h.service.UserService.Update(userParamsToUpdate)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, http.StatusBadRequest, err.Error(), "UserService.GetUserById(): error occured")
	}

	// return updated user
	pkg.Response(w, updatedUser)
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user_id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		pkg.ErrorResponse(w, 400, "только цифры")
		return
	}
	_, err = h.service.UserService.Delete(userId)
	if err != nil {
		log.Println(err)
		pkg.ErrorResponse(w, 500, "sorry")
		return
	}
	pkg.Response(w, "success")
}

type UserValidator struct {
	UserToValidate     *models.User      `json:"passed_user"`
	ValidationErrors   map[string]string `json:"validation_errors"`
	IsNameValid        bool              `json:"is_name_valid"`
	IsRoleIdValid      bool              `json:"is_roleid_valid"`
	IsEmailValid       bool              `json:"is_email_valid"`
	IsTelephoneValid   bool              `json:"is_telephone_valid"`
	AllUserFieldsValid bool              `json:"all_user_fields_valid"`
}

func NewUserValidator(user *models.User) *UserValidator {
	validationErrors := map[string]string{
		"name_error":      "User.Name: should not be empty string",
		"roleid_error":    "User.Role: should not be negative integer or zero",
		"email_error":     "User.Email: wrong email format",
		"telephone_error": "User.Name: wrong telephone number",
	}

	validator := &UserValidator{UserToValidate: user, ValidationErrors: validationErrors, AllUserFieldsValid: false}
	validator.IsUserValid()

	return validator
}

func (u *UserValidator) IsUserValid() {
	u.ValidateFields()

	if u.IsNameValid && u.IsRoleIdValid && u.IsEmailValid && u.IsTelephoneValid {
		u.AllUserFieldsValid = true
	}
}

func (u *UserValidator) ValidateFields() {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	telephoneRegex := regexp.MustCompile(`^(?:\+992|992|8)?9\d{8}$`)

	if u.UserToValidate.Name != "" {
		u.IsNameValid = true
		delete(u.ValidationErrors, "nameError")
	}
	if u.UserToValidate.RoleId > 0 {
		u.IsRoleIdValid = true
		delete(u.ValidationErrors, "roleIdError")
	}
	if emailRegex.MatchString(u.UserToValidate.Email) {
		u.IsEmailValid = true
		delete(u.ValidationErrors, "emailError")
	}
	if telephoneRegex.MatchString(u.UserToValidate.Telephone) {
		u.IsTelephoneValid = true
		delete(u.ValidationErrors, "telephoneError")
	}
}
