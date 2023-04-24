package util

import (
	"encoding/json"
	"log"
	"net/http"

	"realworld-go-nolambda/model"
	//"realworld-go-nolambda/model"
)

type InputErrorResponse struct {
	Errors model.InputError `json:"errors"`
}

func NewErrorResponse(statusCode int, err error, w http.ResponseWriter) {
	EnableCors(&w)

	inputError, ok := err.(model.InputError)

	if !ok {
		// Internal server error
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occured internally"))
	}

	body := InputErrorResponse{
		Errors: inputError,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occured internally"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonBody)

}

func NewUnauthorizedResponse(w http.ResponseWriter) {
	EnableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
}
