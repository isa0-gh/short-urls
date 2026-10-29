package response

import (
	"encoding/json/v2"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func JSON(w http.ResponseWriter, status int, data any) error {
	if data == nil {
		w.WriteHeader(status)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.MarshalWrite(w, data); err != nil {
		return err
	}
	return nil
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{
		"error": message,
	})
}

func Unauthorized(w http.ResponseWriter) {
	JSON(w, http.StatusUnauthorized, map[string]string{
		"error": "Unauthorized",
	})
}

func DatabaseErrorStatus(err error) int {
	if errors.Is(err, pgx.ErrNoRows) {
		return http.StatusNotFound
	}
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}
