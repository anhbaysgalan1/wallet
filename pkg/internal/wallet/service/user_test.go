package service

import (
	"fmt"
	"testing"
)

func TestWalletSrv_SysBalanceAdd(t *testing.T) {
	before, after, err := walletSrv.SysBalanceAdd(ctx, 23, "1000000000000000000000000", "f1", true)
	if err != nil {
		fmt.Println("error ======>", err)
	} else {
		fmt.Println(fmt.Printf("success =====> before:%s, after:%s. ", before, after))
	}
}
