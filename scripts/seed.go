/*
Author: Adrian Potra
Version 1.0.
*/

package main

import (
	"context"
	"fmt"
	"hotel-room-bookme/db"
	"hotel-room-bookme/types"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

// helper func
func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99,
		},
		{
			Size:  "normal",
			Price: 129,
		},
		{
			Size:  "kingsize",
			Price: 229,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted hotel: ", insertedHotel)

	for _, room := range rooms {

		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("inserted room: ", insertedRoom)
	}

}

func main() {

	seedHotel("Bellucia", "France", 3)
	seedHotel("Cozy Hotel", "The Netherlands", 4)

}

func init() {
	var err error
	//database connection
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)

}
