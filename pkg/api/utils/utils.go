package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string
}

func SendAsJSON(w http.ResponseWriter, object interface{}) error {
	json, err := json.Marshal(object)
	if err != nil {
		HttpJSONError(w, err.Error(), http.StatusBadRequest)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
	return nil
}

func RespondWithJSON(w http.ResponseWriter, message string) error {
	var jsonResponse Response
	jsonResponse.Message = message
	err := SendAsJSON(w, jsonResponse)
	if err != nil {
		return err
	}
	return nil
}

func DecodeJSONObject(w http.ResponseWriter, r *http.Request, object interface{}) error {
	err := json.NewDecoder(r.Body).Decode(object)
	if err != nil {
		HttpJSONError(w, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func HttpJSONError(w http.ResponseWriter, error string, http_status_code int) error {
	response := ErrorResponse{
		ErrorMessage: error,
	}
	err := SendAsJSON(w, response)
	if err != nil {
		return err
	}
	return nil
}