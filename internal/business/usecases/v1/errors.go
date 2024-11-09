package v1

import "errors"

var (
	ErrAmountMustGreateThanZero    = errors.New("amount must be greater than zero")
	ErrQuantityMustGreaterThanZero = errors.New("quantity must be greater than zero")
)
