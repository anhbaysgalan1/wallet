package job

import (
	"context"
)

// JobCurrencyWalletTransferToBlock 定时处理异步上链任务
func (s *WalletJob) JobCurrencyWalletTransferToBlock(ctx context.Context) {
	_, _ = s.WalletSrv.JobCurrencyWalletTransferToBlock(ctx, nil)
}

// JobNftWalletTransferToBlock 定时处理异步上链任务
func (s *WalletJob) JobNftWalletTransferToBlock(ctx context.Context) {
	_, _ = s.WalletSrv.JobNftWalletTransferToBlock(ctx, nil)
}
