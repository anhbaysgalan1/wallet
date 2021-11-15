package binance

import (
	"fmt"
	"github.com/leaf-rain/wallet/internal/binance/dto"
	"testing"
)

func Test_binance_AddressGet(t *testing.T) {
	result, _ := binanceSrv.AddressGet(ctx, &dto.AddressGetReq{
		Coin:       "BTC",
		RecvWindow: 60000,
	})
	fmt.Println("result ======>", result)
}

func Test_binance_AccountGet(t *testing.T) {
	result, _ := binanceSrv.AccountGet(ctx, &dto.AccountGetReq{})
	fmt.Println("result ======>", result)
}
