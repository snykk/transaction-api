package utils

import (
	"database/sql"
	"errors"
	"net/http"

	postgresRepo "github.com/snykk/transaction-api/internal/datasources/repositories/postgres/v1"

	"github.com/jackc/pgconn"
)

func MapDBError(err error) (int, error) {
	// Periksa apakah error berasal dari database PostgreSQL
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "42P01":
			return http.StatusInternalServerError, errors.New("database error: table does not exist")
		case "42703":
			return http.StatusInternalServerError, errors.New("database error: column does not exist")
		case "23505":
			return http.StatusConflict, errors.New("database error: unique constraint violation")
		case "23503":
			return http.StatusBadRequest, errors.New("database error: foreign key violation")
		case "40P01":
			return http.StatusConflict, errors.New("database error: deadlock detected")
		case "57014":
			return http.StatusRequestTimeout, errors.New("database error: query timeout")
		default:
			return http.StatusInternalServerError, errors.New("unexpected database error")
		}
	}

	// Periksa error custom untuk insufficient balance dan insufficient stock
	if errors.Is(err, postgresRepo.ErrInsufficientBalance) {
		return http.StatusUnprocessableEntity, postgresRepo.ErrInsufficientBalance
	}
	if errors.Is(err, postgresRepo.ErrInsufficientProductStock) {
		return http.StatusUnprocessableEntity, postgresRepo.ErrInsufficientBalance
	}

	if errors.Is(err, postgresRepo.ErrProductNotFound) {
		return http.StatusNotFound, postgresRepo.ErrProductNotFound
	}

	// Periksa apakah error adalah sql.ErrNoRows
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, errors.New("data not found")
	}

	// Error umum lainnya
	return http.StatusInternalServerError, errors.New("failed to process the request")
}
