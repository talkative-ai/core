package myerrors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type IMyError interface {
	Parse() (int, string)
	GetDepth() int
	SetDepth(int)
}

type MySimpleError struct {
	Code    int           `json:"code"`
	Message interface{}   `json:"message"`
	Log     string        `json:"-"`
	Req     *http.Request `json:"-"`
	Depth   int           `json:"-"`
}

func (simple *MySimpleError) GetDepth() int {
	return simple.Depth
}

func (simple *MySimpleError) SetDepth(depth int) {
	simple.Depth = depth
}

func (simple *MySimpleError) Parse() (int, string) {

	if simple.Log != "" {
		_, fn, line, _ := runtime.Caller(simple.Depth)
		log.Printf("[ERROR] %s:%d %v", fn, line, simple.Log)
	}

	encoded, err := json.Marshal(simple)
	if err != nil {
		fmt.Println("Problem encoding error", err)
	}

	return simple.Code, string(encoded)
}

func Respond(w http.ResponseWriter, err IMyError) {
	if err.GetDepth() <= 1 {
		err.SetDepth(2)
	}
	code, msg := err.Parse()
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
}

func ServerError(w http.ResponseWriter, rq *http.Request, err error) {
	Respond(w, &MySimpleError{
		Code:  http.StatusInternalServerError,
		Log:   err.Error(),
		Req:   rq,
		Depth: 3,
	})
}
