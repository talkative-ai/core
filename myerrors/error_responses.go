package myerrors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IMyError interface {
	Parse() (int, string)
}

type MySimpleError struct {
	Code    int           `json:"code"`
	Message interface{}   `json:"message"`
	Log     string        `json:"-"`
	Req     *http.Request `json:"-"`
}

func (simple *MySimpleError) Parse() (int, string) {

	if simple.Log != "" {
		fmt.Println(simple.Log)
	}

	encoded, err := json.Marshal(simple)
	if err != nil {
		fmt.Println("Problem encoding error", err)
	}

	return simple.Code, string(encoded)
}

func Respond(w http.ResponseWriter, err IMyError) {
	code, msg := err.Parse()
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
}

func ServerError(w http.ResponseWriter, rq *http.Request, err error) {
	Respond(w, &MySimpleError{
		Code: http.StatusInternalServerError,
		Log:  err.Error(),
		Req:  rq,
	})
}
