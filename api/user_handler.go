/*
Author: Adrian Potra
Version: 1.0
*/
package api

import (
	"errors"
	"hotel-room-bookme/db"
	"hotel-room-bookme/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// we need a way to access the db
type UserHandler struct {
	userStore db.UserStore
}

// we make a constructor func
func NewUserHandler(userStor db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStor,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		//values bson.M
		params types.UpdateUserParams
		userID = c.Params("id")
	)

	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	// filter - a way to choose to update specific mongo documents
	filter := db.Map{"_id": userID}

	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": userID})

}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {

	userID := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"delete": userID})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	// encode JSON request into params
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error:": "not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return ErrResourceNotFound("user")
	}

	return c.JSON(users)
}
