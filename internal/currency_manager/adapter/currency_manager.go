package adapter

import (
	"github.com/leaf-rain/wallet/internal/currency_manager/dto"
	"github.com/leaf-rain/wallet/internal/currency_manager/service"
)

func NewCurrencyMangerSrv() dto.CurrencyManagerSrvServer {
	return &service.CurrencyMangerSrv{
		CurrencyList: service.NewCurrencyCfg(),
		NetList:      service.NewNetCfg(),
	}
}
