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
	GetAllWallets(ctx context.Context) (domains []WalletDomain, statusCode int, err error)
	Init(ctx context.Context, userId string) (domain WalletDomain, statusCode int, err error)
	GetWalletByUserId(ctx context.Context, userId string) (domain WalletDomain, statusCode int, err error)
}

type WalletRepository interface {
	GetAllWallets(ctx context.Context) ([]WalletDomain, error)
	CreateWalletByUserId(ctx context.Context, userId string) (WalletDomain, error)
	GetWalletByUserId(ctx context.Context, userId string) (WalletDomain, error)
}
