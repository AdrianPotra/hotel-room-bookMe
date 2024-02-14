package api

import (
	"encoding/json"
	"fmt"
	"hotel-room-bookme/db/fixtures"
	"hotel-room-bookme/types"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	user := fixtures.AddUser(&db.store, "james", "foo", false)
	hotel := fixtures.AddHotel(&db.store, "bar hotel", "Some Location", 4, nil)
	room := fixtures.AddRoom(&db.store, "small", true, 4.4, hotel.ID)
	from := time.Now()
	till := from.AddDate(0, 0, 3)
	nrPers := 2
	booking := fixtures.AddBooking(&db.store, user.ID, room.ID, from, till, nrPers)
	//fmt.Println(booking)
	_ = booking
	app := fiber.New()
	bookingHandler := NewBookingHandler(&db.store)
	app.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil) // we don't have body, so it's nil
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var bookings []*types.Booking

	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	fmt.Println(bookings)

}
