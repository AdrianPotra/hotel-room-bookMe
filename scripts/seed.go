/*
Author: Adrian Potra
Version 1.0.
*/

package main

import (
	"context"
	"fmt"
	"hotel-room-bookme/api"
	"hotel-room-bookme/db"
	"hotel-room-bookme/db/fixtures"
	"hotel-room-bookme/types"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})

	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := userStore.InsertUser(context.TODO(), user)

	if err != nil {
		log.Fatal(err)
	}
	// this is to have the tokens printed out in the client, so that we don't need to authorize ourselves via api calls - can be removed later
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}

func seedRoom(size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

// helper func
func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted hotel: ", insertedHotel)

	return insertedHotel

}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time, nrPers int) {
	booking := &types.Booking{
		UserID:     userID,
		RoomID:     roomID,
		FromDate:   from,
		TillDate:   till,
		NumPersons: nrPers,
	}
	resp, err := bookingStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("booking ID :", resp.ID)

}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		err           error
		ctx           = context.Background()
		mongoEndpoint = os.Getenv("MONGO_DB_URL")
		mongoDBName   = os.Getenv("MONGO_DB_NAME")
	)
	//database connection
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   db.NewMongoHotelStore(client),
	}

	user := fixtures.AddUser(store, "james", "thesecond", false)
	fmt.Println("user from fixtures: ", user)
	//admin :=  fixtures.AddUser(store, "admin", "admin", true)
	//fmt.Println("admin user from fixtures: ", admin)
	hotel := fixtures.AddHotel(store, "Some Hotel", "Alaska", 5, nil)
	fmt.Println("hotel from fixtures: ", hotel)
	room := fixtures.AddRoom(store, "large", true, 107.15, hotel.ID)
	fmt.Println("room from fixtures: ", room)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 1), 2)
	fmt.Println("booking from fixtures: ", booking)
	fmt.Println("booking ID from fixtures: ", booking.ID)

	for i := 0; i < 51; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}

	hames := seedUser(false, "hames", "foo", "hames.foo@goo.com", "supersecurepassword")
	seedUser(true, "admin", "admin", "admin@admin.com", "admin123")
	h1 := seedHotel("Bellucia", "France", 3)
	h2 := seedHotel("Cozy Hotel", "The Netherlands", 4)
	seedRoom("small", false, 89.99, h1.ID)
	room2H1 := seedRoom("normal", false, 119.99, h1.ID)
	seedRoom("kingsize", true, 179.99, h1.ID)
	seedRoom("small", false, 89.99, h2.ID)
	seedRoom("normal", false, 119.99, h2.ID)
	seedRoom("kingsize", true, 179.99, h2.ID)
	seedBooking(hames.ID, room2H1.ID, time.Now(), time.Now().AddDate(0, 0, 1), 2)

}

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		err           error
		mongoEndpoint = os.Getenv("MONGO_DB_URL")
		mongoDBName   = os.Getenv("MONGO_DB_NAME")
	)
	//database connection
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)

}
