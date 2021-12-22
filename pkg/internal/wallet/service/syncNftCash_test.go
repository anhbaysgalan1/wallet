package service

import (
	"fmt"
	"testing"
)

func TestWalletSrv_JobNftWalletTransferToBlock(t *testing.T) {
	_, err := walletSrv.JobNftWalletTransferToBlock(ctx, nil)
	if err != nil {
		fmt.Println("failed =====>", err)
	}
}
