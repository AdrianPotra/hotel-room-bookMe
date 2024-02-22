/*
Author: Adrian Potra
Version 1.0.
*/

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hotel-room-bookme/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.store.User)
	app.Post("/", userHandler.HandlePostUser)

	// the bytes send as JSON to our handler
	params := types.CreateUserParams{
		Email:     "somefoo@foo.com",
		FirstName: "Baba",
		LastName:  "Yaga",
		Password:  "fgflgfdlkgjdfjgdfgjdfgj",
	}
	//marshal the params to bytes
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	fmt.Println("user: ", user)
	fmt.Println("Response status: ", resp.Status)
	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the encrypted password not to be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}

}
