package errors

import (
	"encoding/json"
	"net/http"
)

func JsonError(w http.ResponseWriter, status int, msg string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"message": msg,
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}
