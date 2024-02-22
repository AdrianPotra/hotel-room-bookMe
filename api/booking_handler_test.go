package api

import (
	"encoding/json"
	"fmt"
	"hotel-room-bookme/db/fixtures"
	"hotel-room-bookme/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		nonAuthUser    = fixtures.AddUser(&db.store, "Jimmy", "Jones", false)
		user           = fixtures.AddUser(&db.store, "james", "foo", false)
		hotel          = fixtures.AddHotel(&db.store, "bar hotel", "Some Location", 4, nil)
		room           = fixtures.AddRoom(&db.store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 3)
		nrPers         = 2
		booking        = fixtures.AddBooking(&db.store, user.ID, room.ID, from, till, nrPers)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		route          = app.Group("/", JWTAuthentication(db.store.User))
		bookingHandler = NewBookingHandler(&db.store)
	)

	_ = app

	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil) // we don't have body, so it's nil
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 code, got %d", resp.StatusCode)
	}

	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

	if bookingResp.ID != booking.ID {
		t.Fatalf("expected %s , got %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected %s , got %s", booking.UserID, bookingResp.UserID)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected non 200 status code, got %d", resp.StatusCode)
	}

}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		adminUser      = fixtures.AddUser(&db.store, "admin", "admin", true)
		user           = fixtures.AddUser(&db.store, "james", "foo", false)
		hotel          = fixtures.AddHotel(&db.store, "bar hotel", "Some Location", 4, nil)
		room           = fixtures.AddRoom(&db.store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 3)
		nrPers         = 2
		booking        = fixtures.AddBooking(&db.store, user.ID, room.ID, from, till, nrPers)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		admin          = app.Group("/", JWTAuthentication(db.store.User), AdminAuth)
		bookingHandler = NewBookingHandler(&db.store)
	)
	//fmt.Println(booking)
	_ = booking

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil) // we don't have body, so it's nil
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response, got %d", resp.StatusCode)
	}

	var bookings []*types.Booking

	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got none %d", len(bookings))
	}

	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s , got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s , got %s", booking.UserID, have.UserID)
	}

	fmt.Println("Booking Test Made is: ", bookings)

	// test non-admin cannot access the bookings
	admin.Get("/", bookingHandler.HandleGetBookings)
	req = httptest.NewRequest("GET", "/", nil) // we don't have body, so it's nil
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized, got %d", resp.StatusCode)
	}

}
