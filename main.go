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
	"hotel-room-bookme/middleware"
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
		appv1          = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin          = appv1.Group("/admin", middleware.AdminAuth)
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

	app.Listen(*listenAddr)
}
