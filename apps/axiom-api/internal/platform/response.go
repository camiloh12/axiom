package platform

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to write JSON response", "error", err)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*AppError); ok {
		WriteJSON(w, appErr.Code, appErr)
		return
	}
	slog.Error("unexpected error", "error", err)
	WriteJSON(w, http.StatusInternalServerError, &AppError{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
	})
}
