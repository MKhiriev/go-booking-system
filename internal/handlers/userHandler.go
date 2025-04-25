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

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// get all users
	users := h.service.UserService.GetAll()

	// return all users
	pkg.Response(w, users)
}

func (h *Handlers) GetUserById(w http.ResponseWriter, r *http.Request) {
	// get user_id from query path
	userIdStr := r.URL.Query().Get("user_id")
	if userIdStr == "" {
		log.Println("UserHandler.DeleteRoom(): parameter `user_id` is empty or not passed")
		pkg.ErrorResponse(w, http.StatusBadRequest, "parameter `user_id` is empty or not passed")
		return
	}

	// convert user_id param string to int
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Println("UserHandler.GetUserById(): user_id should be an integer. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "user_id should be an integer", err.Error())
		return
	}

	// get User from services
	user, err := h.service.UserService.GetUserById(userId)
	if err != nil {
		log.Println("UserHandler.GetUserById(): error occured during getting user by id. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during getting user by id", err.Error())
		return
	}

	// return found user
	pkg.Response(w, user)
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// get user_id from query path
	userIdStr := r.URL.Query().Get("user_id")
	if userIdStr == "" {
		log.Println("UserHandler.UpdateUser(): parameter `user_id` is empty or not passed")
		pkg.ErrorResponse(w, http.StatusBadRequest, "parameter `user_id` is empty or not passed")
		return
	}

	// convert user_id param string to int
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Println("UserHandler.UpdateUser(): user_id should be an integer. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "user_id should be an integer", err.Error())
		return
	}

	var userParamsToUpdate models.User
	// convert JSON to models.User type
	err = json.NewDecoder(r.Body).Decode(&userParamsToUpdate)
	if err != nil {
		log.Println("UserHandler.UpdateUser(): cannot convert JSON to models.User struct. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "cannot convert JSON to models.User struct", err.Error())
		return
	}

	// validate passed user data
	validator := NewUserValidator(&userParamsToUpdate)
	if validator.AllUserFieldsValid != true {
		log.Println("UserHandler.UpdateUser(): User data is not valid. Details: ", validator.ValidationErrors)
		pkg.ErrorResponse(w, http.StatusBadRequest, "User data is not valid", validator.ValidationErrors)
		return
	}

	// to double-check if user_id wasn't set
	userParamsToUpdate.UserId = userId

	// update user
	updatedUser, err := h.service.UserService.Update(userParamsToUpdate)
	if err != nil {
		log.Println("UserHandler.UpdateUser(): error occured during user update. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during user update", err.Error())
		return
	}

	// return updated user
	pkg.Response(w, updatedUser)
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// get user_id from query path
	userIdStr := r.URL.Query().Get("user_id")
	if userIdStr == "" {
		log.Println("UserHandler.DeleteUser(): parameter `user_id` is empty or not passed")
		pkg.ErrorResponse(w, http.StatusBadRequest, "parameter `user_id` is empty or not passed")
		return
	}

	// convert user_id param string to int
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Println("UserHandler.DeleteUser(): user_id should be an integer. Details: ", err)
		pkg.ErrorResponse(w, http.StatusBadRequest, "user_id should be an integer", err.Error())
		return
	}

	// delete user
	_, err = h.service.UserService.Delete(userId)
	if err != nil {
		log.Println("UserHandler.DeleteUser(): error occured during user deletion. Details: ", err)
		pkg.ErrorResponse(w, http.StatusInternalServerError, "error occured during user deletion", err.Error())
		return
	}

	// return success message
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
		"roleid_error":    "User.RoleId: should not be negative integer or zero",
		"email_error":     "User.Email: wrong email format",
		"telephone_error": "User.Telephone: wrong telephone number",
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
		delete(u.ValidationErrors, "name_error")
	}
	if u.UserToValidate.RoleId > 0 {
		u.IsRoleIdValid = true
		delete(u.ValidationErrors, "roleid_error")
	}
	if emailRegex.MatchString(u.UserToValidate.Email) {
		u.IsEmailValid = true
		delete(u.ValidationErrors, "email_error")
	}
	if telephoneRegex.MatchString(u.UserToValidate.Telephone) {
		u.IsTelephoneValid = true
		delete(u.ValidationErrors, "telephone_error")
	}
}
