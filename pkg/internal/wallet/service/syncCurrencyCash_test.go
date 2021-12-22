package service

import (
	"fmt"
	"testing"
)

func TestWalletSrv_JobCurrencyWalletTransferToBlock(t *testing.T) {
	_, err := walletSrv.JobCurrencyWalletTransferToBlock(ctx, nil)
	if err != nil {
		fmt.Println("failed =====>", err)
	}
}
