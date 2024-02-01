/*
 Author: Adrian Potra
 Version: 1.0

 to import gofiber - github.com/gofiber/fiber/v3
 - also we will use mongodb, so will have to install driver dependency - go get go.mongodb.org/mongo-driver/mongo

 mongodb+srv://APotra:<password>@adrianscluster.xfsrbdy.mongodb.net/?retryWrites=true&w=majority

*/

package main

import (
	"context"
	"flag"
	"hotel-room-bookme/api"
	"hotel-room-bookme/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// fiber error handler
var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":8080", "The listen Address of the API server")
	flag.Parse()

	//database connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// handler initialization
	var (
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)
		store      = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		app          = fiber.New(config)
		appv1        = app.Group("/api/v1")
	)

	// user handlers
	appv1.Put("/user/:id", userHandler.HandlePutUser)
	appv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	appv1.Post("/user", userHandler.HandlePostUser)
	appv1.Get("/user", userHandler.HandleGetUsers)
	appv1.Get("/user/:id", userHandler.HandleGetUser)
	// hotel handlers
	appv1.Get("/hotel", hotelHandler.HandleGetHotels)
	appv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	appv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	app.Listen(*listenAddr)
}
