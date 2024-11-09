package v1

import "errors"

// Error custom untuk kondisi bisnis
var (
	ErrInsufficientBalance      = errors.New("insufficient balance")
	ErrInsufficientProductStock = errors.New("insufficient product stock")
	ErrProductNotFound          = errors.New("product not found")
)
