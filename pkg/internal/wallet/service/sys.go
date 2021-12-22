package service

import (
	"context"
	"tp_wallet/internal/wallet/entity"
)

// SysBalanceAdd 直接加钱
func (srv WalletSrv) SysBalanceAdd(ctx context.Context, uid uint64, amount, currency string, isAdd bool) (entity.Amount, entity.Amount, error) {
	return srv.Repo.BalanceSet(ctx, uid, amount, currency, isAdd)
}
