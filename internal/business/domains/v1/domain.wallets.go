package v1

import (
	"context"
	"time"
)

type WalletDomain struct {
	Id        string
	UserId    string
	Balance   float64
	User      UserDomain
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type WalletUsecase interface {
	Init(ctx context.Context, userId string) (domain WalletDomain, statusCode int, err error)
	GetByUserId(ctx context.Context, userId string) (domain WalletDomain, statusCode int, err error)
}

type WalletRepository interface {
	CreateByUserId(ctx context.Context, userId string) (WalletDomain, error)
	GetByUserId(ctx context.Context, userId string) (WalletDomain, error)
}