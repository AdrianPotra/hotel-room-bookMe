/*
Author: Adrian Potra
Version 1.0.

for password encryption we will use golang.org/x/crypto/bcrypt package
*/

package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12 // some random nr for now
	minFirstNameLen = 2
	minLastNameLen  = 3
	minPasswordLen  = 6
)

type UpdateUserParams struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{} // initializing map

	if len(p.FirstName) > 0 {
		m["firstname"] = p.FirstName
	}

	if len(p.LastName) > 0 {
		m["lastname"] = p.LastName
	}

	return m
}

// for the request scope
type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json: "email"`
	Password  string `json: "password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email is invalid")
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

// for the domain scope
type User struct {
	ID                primitive.ObjectID `bson: "_id,omitempty"  json:"id,omitempty"`
	FirstName         string             `bson: "firstName"  json:"firstName"`
	LastName          string             `bson: "lastName"  json:"lastName"`
	Email             string             `bson: "email" json: "email"`
	EncryptedPassword string             `bson: "EncryptedPassword" json: "-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw), // encpw is going to be bytes, so we cast into string
	}, nil

}