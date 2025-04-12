package wallet

import "github.com/korroziea/taxi/user-service/internal/domain"

type walletResp struct {
	ID      string                `json:"id"`
	Type    domain.WalletType     `json:"type"`
	Role    domain.UserWalletRole `json:"role"`
	Balance int64                 `json:"balance"`
}

func toWalletView(wallet domain.ViewWallet) walletResp {
	resp := walletResp{
		ID:      wallet.ID,
		Type:    wallet.Type,
		Role:    wallet.Role,
		Balance: wallet.Balance,
	}

	return resp
}

type walletListResp struct {
	Wallets []walletResp `json:"wallets"`
}

func toWalletListView(wallets []domain.ViewWallet) walletListResp {
	viewWallets := make([]walletResp, len(wallets))
	for i := range viewWallets {
		viewWallets[i] = toWalletView(wallets[i])
	}

	resp := walletListResp{
		Wallets: viewWallets,
	}

	return resp
}
