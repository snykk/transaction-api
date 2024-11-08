package requests

import (
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type TransactionDepositOrWithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"` // price lebih besar dari 0
	// PaymentMethod string  `json:"payment_method" binding:"required,oneof=bank_transfer credit_card"`
	// Description string `json:"description" binding:"required"`
}

func (w *TransactionDepositOrWithdrawRequest) ToDomain() *V1Domains.TransactionDomain {
	return &V1Domains.TransactionDomain{
		Amount: w.Amount,
		// Description: w.Description,
	}
}

type TransactionPurchaseRequest struct {
	ProductId int `json:"product_id" binding:"required"`    // price lebih besar dari 0
	Quantity  int `json:"quantity" binding:"required,gt=0"` // price lebih besar dari 0
}

func (w *TransactionPurchaseRequest) ToDomain() *V1Domains.TransactionDomain {
	return &V1Domains.TransactionDomain{
		ProductId: &w.ProductId,
		Quantity:  &w.Quantity,
	}
}
