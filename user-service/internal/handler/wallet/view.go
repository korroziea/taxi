package wallet

import "github.com/korroziea/taxi/user-service/internal/domain"

type createWalletResp struct {
	ID      string            `json:"id"`
	Type    domain.WalletType `json:"type"`
	Balance int64             `json:"balance"`
}

func toView(wallet domain.Wallet) createWalletResp {
	resp := createWalletResp{
		ID:      wallet.ID,
		Type:    wallet.Type,
		Balance: wallet.Balance,
	}

	return resp
}
