package store

import "web3-account-abstraction-api/internal/model"

type Store interface {
	CountWallet() (int, error)
	CreateWallet(model.UserWallet) error
	GetAllWallet() ([]model.UserWallet, error)
}
