package models

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

// User TODO add getters/setters
// User TODO add username and password
type User struct {
	UserId    int    `json:"user_id" gorm:"primarykey"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
	RoleId    int    `json:"role_id"`

	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (u *User) SetName(name string) error {
	if name == "" {
		return errors.New("cannot use empty string as User.Name")
	}
	u.Name = name
	return nil
}

// SetEmail TODO add regexp check
func (u *User) SetEmail(email string) error {
	if email == "" {
		return errors.New("cannot use empty string as User.Email")
	}
	u.Email = email
	return nil
}

func (u *User) SetTelephone(telephone string) error {
	if telephone == "" {
		return errors.New("cannot use empty string as User.Telephone")
	}
	u.Telephone = telephone
	return nil
}

func (u *User) SetRole(roleId int) error {
	u.RoleId = roleId
	return nil
}

func (u *User) SetActive(active bool) error {
	u.Active = active
	return nil
}

func (u *User) SetCreatedAt(createdAt time.Time) error {
	if createdAt.Unix() > time.Now().Unix() {
		return errors.New("cannot assign creation datetime greater than current time as User.CreatedAt")
	}
	u.CreatedAt = createdAt
	return nil
}

func ParseUser(data []byte) (*User, error) {
	user := &User{}
	err := json.Unmarshal(data, user)
	if err != nil {
		log.Println("Models.ParseUserData error: ", err)
		return user, err
	}
	return user, nil
}
