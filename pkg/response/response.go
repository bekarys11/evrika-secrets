package resp

import (
	"encoding/json"
	"errors"
	"github.com/google/jsonapi"
	"io"
	"log/slog"
	"net/http"
)

type Err struct {
	Message string `json:"message" example:"internal server error"`
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) {
	out, err := json.Marshal(data)

	if err != nil {
		slog.Error(err.Error())
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(out); err != nil {
		slog.Error(err.Error())
	}
}

func WriteApiJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) {
	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := jsonapi.MarshalPayload(w, data); err != nil {
		slog.Error(err.Error())
	}
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must have only single JSON value")
	}

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload Err
	payload.Message = err.Error()

	WriteJSON(w, statusCode, payload)
}
