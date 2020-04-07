package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// StatusResponse is a status response returned by an API endpoint
type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Validable is an interface to implement to validate a struct while parsing it
type Validable interface {
	Validate() error
}

// ReadJSON reads a json payload and unmarshal it to an interface
func ReadJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "unable to read body")
		return err
	}

	err = json.Unmarshal(requestBody, v)
	if err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "unable to parse body")
		return err
	}
	return nil
}

// ReadValidateJSON reads a json payload and unmarshal it to an interface and validates it
func ReadValidateJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	err := ReadJSON(w, r, v)
	if err != nil {
		return err
	}

	validable, ok := v.(Validable)
	if !ok {
		return fmt.Errorf("value doesn't implement Validable interface")
	}

	err = validable.Validate()
	if err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

// WriteJSON writes json value
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	j, _ := json.Marshal(v)

	w.Write(j)
}

// WriteErrorJSON writes a json error
func WriteErrorJSON(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, StatusResponse{
		Status:  "error",
		Message: message,
	})
}
