/*
Author: Adrian Potra
Version 1.0.

we are going to use jtw parse sample code from https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac

*/

package middleware

import (
	"fmt"
	"hotel-room-bookme/db"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("---JWT authentication---")

		token, ok := c.GetReqHeaders()["X-Api-Token"] // most of the time x-aoi-token is something custom, so we call it like that

		if !ok {
			return fmt.Errorf("Unauthorized")
		}

		tokenS := strings.Join(token, " ")
		fmt.Println("token to String is ", tokenS)

		claims, err := validateToken(tokenS)
		if err != nil {
			return err
		}
		fmt.Println("claims: ", claims)
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		// check token expiration
		if time.Now().Unix() > expires {
			return fmt.Errorf("token expired")
		}

		fmt.Println("expires: ", expires)
		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			fmt.Println("getuserbyid error")
			return fmt.Errorf("Unauthorized")
		}
		fmt.Println("user before it adds to context: ", user)
		//Set the current authenticated user to the context
		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method ", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		// getting environment var from go
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("NEVER PRINT SECRET ", secret)
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil // casting the string to bytes
	})

	if err != nil {
		fmt.Println("failed to parse JWT token: ", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil

}
