package adapter

import (
	"context"
	"github.com/leaf-rain/wallet/internal/common"
)

type AccountAdapter interface {
	// AddressGet 获取帐号
	AddressGet(ctx context.Context, currency common.Currency) (address string)
}
