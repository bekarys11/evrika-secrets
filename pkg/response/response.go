package resp

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type New struct {
	Data any `json:"data"`
}

type Err struct {
	Message string `json:"message"`
}

//type Data struct {
//	ID         int    `json:"id,omitempty"`
//	Type       string `json:"type,omitempty"`
//	Attributes any    `json:"attributes"`
//	Included   any    `json:"included"`
//}

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

func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload Err
	payload.Message = err.Error()

	WriteJSON(w, statusCode, payload)
}
