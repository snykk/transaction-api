package utils

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/jackc/pgconn"
)

// MapDBError mengembalikan status code dan pesan error yang sesuai
func MapDBError(err error) (int, error) {
	// Periksa apakah error berasal dari database PostgreSQL
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "42P01": // Tabel tidak ditemukan
			return http.StatusInternalServerError, errors.New("database error: table does not exist")
		case "42703": // Kolom tidak ditemukan
			return http.StatusInternalServerError, errors.New("database error: column does not exist")
		case "23505": // Pelanggaran constraint UNIQUE (duplikat nilai pada kolom unik)
			return http.StatusConflict, errors.New("database error: unique constraint violation")
		case "23503": // Pelanggaran constraint FOREIGN KEY
			return http.StatusBadRequest, errors.New("database error: foreign key violation")
		case "40P01": // Deadlock detected
			return http.StatusConflict, errors.New("database error: deadlock detected")
		case "57014": // Query timeout
			return http.StatusRequestTimeout, errors.New("database error: query timeout")
		default:
			return http.StatusInternalServerError, errors.New("unexpected database error")
		}
	}

	// Periksa apakah error adalah sql.ErrNoRows
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, errors.New("data not found")
	}

	// Error umum lainnya
	return http.StatusInternalServerError, errors.New("failed to process the request")
}
