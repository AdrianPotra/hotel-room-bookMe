/*
Author: Adrian Potra
Version 1.0.
*/

package api

import (
	"errors"
	"fmt"
	"hotel-room-bookme/db"
	"hotel-room-bookme/types"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

// A handler should only do:
// - serialization of the incoming request
// - do some data fetching from db
// - call some business logic, not do business logic itself
// - return the data back to the user

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var Params AuthParams
	if err := c.BodyParser(&Params); err != nil {
		return err
	}
	fmt.Println("-----", Params)
	// find user by email
	user, err := h.userStore.GetUserByEmail(c.Context(), Params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, Params.Password) {
		return invalidCredentials(c)
	}
	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}

	return c.JSON(resp)
}

// make claims and then tokens
func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 48).Unix()

	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}

	return tokenStr
}
