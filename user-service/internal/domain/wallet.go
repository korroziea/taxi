package domain

import (
	"time"

	"github.com/korroziea/taxi/user-service/pkg/utils"
)

type WalletType string

const (
	Personal WalletType = "personal"
	Family   WalletType = "family"
)

type Wallet struct {
	ID        string
	Type      WalletType
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GenWalletID() (string, error) {
	id, err := utils.GenID()
	if err != nil {
		return "", err
	}

	return "wallet_" + id, nil
}

type UserWalletRole string

const (
	Owner  UserWalletRole = "owner"
	Member UserWalletRole = "member"
)

type UserWallet struct {
	UserID    string
	WalletID  string
	Role      UserWalletRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ViewWallet struct {
	ID      string
	Type    WalletType
	Role    UserWalletRole
	Balance int64
}
