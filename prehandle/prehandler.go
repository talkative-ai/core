package prehandle

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"time"

	"github.com/artificial-universe-maker/go-utilities"
	"github.com/artificial-universe-maker/go-utilities/myerrors"
	jwt "github.com/dgrijalva/jwt-go"
)

// Prehandler type is exactly the same as http.HandlerFunc except that a return bool is expected to indicate success/failure
type Prehandler func(http.ResponseWriter, *http.Request) bool

// PreHandle accepts an http.HandlerFunc and preprocesses it with n-prehandlers.
// If any prehandler returns false, the process will be aborted and the handler will never be reached
func PreHandle(handle http.HandlerFunc, prehandlers ...Prehandler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		for _, pre := range prehandlers {
			if !pre(w, r) {
				// The prehandler signals an abort
				return
			}
		}
		handle(w, r)
	}

}

// SetJSON sets the Content-Type to application/json
func SetJSON(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Content-Type", "application/json")
	return true
}

func roundD(val float64) int64 {
	if val < 0 {
		return int64(val - 1.0)
	}
	return int64(val)
}

// JWT ensures that the x-token JWT does exist and is valid
func JWT(w http.ResponseWriter, r *http.Request) bool {
	tkn := r.Header.Get("x-token")

	token, err := utilities.ParseJTWClaims(tkn)
	if err != nil {
		myerrors.Respond(w, &myerrors.MySimpleError{
			Req:     r,
			Code:    http.StatusUnauthorized,
			Message: "JWT_INVALID",
		})
		return false
	}

	if err != nil {
		myerrors.Respond(w, &myerrors.MySimpleError{
			Req:     r,
			Code:    http.StatusUnauthorized,
			Message: "JWT_INVALID",
		})
		return false
	}

	if roundD(token["exp"].(float64)) < time.Now().Unix() {
		myerrors.Respond(w, &myerrors.MySimpleError{
			Req:     r,
			Code:    http.StatusUnauthorized,
			Message: "JWT_EXPIRED",
		})
		return false
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(time.Minute * 6000).Unix(),
		"data": token["data"],
	})

	tokenString, err := newToken.SignedString([]byte(os.Getenv("JWT_KEY")))

	w.Header().Set("x-token", tokenString)
	return true
}

// RequireBody forces a body to exist with a maximum length lim
// If the body does not exist, an http.StatusBadRequest is returned. This is required for POST requests
// This prehandler protects against overflows and null-pointer exceptions
func RequireBody(lim int64) Prehandler {
	return func(w http.ResponseWriter, r *http.Request) bool {
		if r.Body == nil {
			myerrors.Respond(w, &myerrors.MySimpleError{
				Req:     r,
				Code:    http.StatusBadRequest,
				Message: "EMPTY_BODY",
			})
			return false
		}

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, lim))
		if err != nil {
			myerrors.Respond(w, &myerrors.MySimpleError{
				Req:     r,
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return false
		}

		if len(body) <= 0 {
			myerrors.Respond(w, &myerrors.MySimpleError{
				Req:     r,
				Code:    http.StatusBadRequest,
				Message: "EMPTY_BODY",
			})
			return false
		}

		r.Header.Set("X-Body", string(body))

		return true
	}
}
