/*
Author: Adrian Potra
Version 1.0.
*/

package fixtures

import (
	"context"
	"fmt"
	"hotel-room-bookme/db"
	"hotel-room-bookme/types"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})

	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)

	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddHotel(store *db.Store, name string, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDS,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted hotel: ", insertedHotel)

	return insertedHotel
}

func AddRoom(store *db.Store, size string, seaside bool, price float64, hid primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hid,
	}

	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(store *db.Store, usrID, roomID primitive.ObjectID, from, till time.Time, nrPers int) *types.Booking {
	booking := &types.Booking{
		UserID:     usrID,
		RoomID:     roomID,
		FromDate:   from,
		TillDate:   till,
		NumPersons: nrPers,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
