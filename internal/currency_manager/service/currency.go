package service

import (
	"context"
	"github.com/leaf-rain/wallet/internal/currency_manager/dto"
	"github.com/leaf-rain/wallet/internal/currency_manager/entity"
	"github.com/leaf-rain/wallet/pkg/hcode"
)

type CurrencyMangerSrv struct {
	CurrencyList *entity.CurrencyList
	NetList      *entity.CurrencyNetList
}

func NewCurrencyMangerSrv() dto.CurrencyManagerSrvServer {
	return &CurrencyMangerSrv{
		CurrencyList: NewCurrencyCfg(),
		NetList:      NewNetCfg(),
	}
}

// CurrencyGetForList 获取所有币种信息
func (cms CurrencyMangerSrv) CurrencyGetForList(ctx context.Context, req *dto.Empty) (*dto.CurrencyList, error) {
	var result = new(dto.CurrencyList)
	result.CurrencyList = make([]*dto.Currency, len(*cms.CurrencyList))
	for index, item := range *cms.CurrencyList {
		result.CurrencyList[index] = item.ToPb()
	}
	return result, nil
}
func (cms CurrencyMangerSrv) CurrencyGetForMap(ctx context.Context, req *dto.Empty) (*dto.CurrencyMap, error) {
	var result = new(dto.CurrencyMap)
	result.CurrencyMap = make(map[string]*dto.Currency)
	for _, item := range *cms.CurrencyList {
		result.CurrencyMap[item.Name] = item.ToPb()
	}
	return result, nil
}

// NetGetByCy 获取转账网络信息
func (cms CurrencyMangerSrv) NetGetByCy(ctx context.Context, req *dto.NameReq) (*dto.NetList, error) {
	return nil, hcode.ErrServer
}
func (cms CurrencyMangerSrv) NetGetByName(ctx context.Context, req *dto.NameReq) (*dto.Net, error) {
	return nil, hcode.ErrServer
}
