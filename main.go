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
	"hotel-room-bookme/api"
	"hotel-room-bookme/db"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Configuration
// 1. MongoDB endpoint
// 2. listenAddress of our HTTP server
// 3. JWT secret
// 4. MongoDB Name

// fiber error handler
var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {

	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	//database connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	// handler initialization
	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}

		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		appv1          = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin          = appv1.Group("/admin", api.AdminAuth)
	)

	//auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

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

	//room handlers
	appv1.Get("room/", roomHandler.HandleGetRooms)
	appv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	//booking handlers

	appv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	// something to do - cancel a booking
	appv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
