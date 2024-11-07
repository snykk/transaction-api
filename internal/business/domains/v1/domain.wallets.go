package v1

import (
	"time"
)

type WalletDomain struct {
	Id        int
	UserId    int
	Balance   float64
	User      UserDomain
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type WalletUsecase interface {
}

type WalletRepository interface {
}
