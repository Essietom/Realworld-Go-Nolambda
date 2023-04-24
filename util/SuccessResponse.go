package util

import (
	"encoding/json"
	"net/http"
)
type AppResponse struct {
    StatusCode      int               `json:"statusCode"`
    Headers         map[string]string `json:"headers"`
    Body            string            `json:"body"`
    IsBase64Encoded bool              `json:"isBase64Encoded,omitempty"`
}


func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
   	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func NewSuccessResponse(body interface{}, w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	if (*r).Method == "OPTIONS"{
		w.WriteHeader(http.StatusOK)
		return
	}
	// fields["status"] = "success"
	message, err := json.Marshal(body)
	if err != nil {
		//An error occurred processing the json
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occured internally"))
	}


	// if body != nil {
	// 	jsonBody, err := json.Marshal(body)
	// 	if err != nil {
	// 		// return NewErrorResponse(err)
	// 	}
	// 	response.Body = string(jsonBody)
	// }

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(message)
}
