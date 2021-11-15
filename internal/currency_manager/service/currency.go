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
	dto.UnimplementedCurrencyManagerSrvServer
}

// CurrencyGetForList 获取所有币种信息
func (cms CurrencyMangerSrv) CurrencyGetForList(ctx context.Context, req *dto.Empty) (*dto.CurrencyList, error) {
	var result = new(dto.CurrencyList)
	result.CurrencyList = make([]*dto.Currency, len(*cms.CurrencyList))
	for index, item := range *cms.CurrencyList {
		result.CurrencyList[index] = item.ToPb()
	}
	if len(result.GetCurrencyList()) <= 0 {
		return nil, hcode.ErrCurrencyNotFound
	}
	return result, nil
}
func (cms CurrencyMangerSrv) CurrencyGetForMap(ctx context.Context, req *dto.Empty) (*dto.CurrencyMap, error) {
	var result = new(dto.CurrencyMap)
	result.CurrencyMap = make(map[string]*dto.Currency)
	for _, item := range *cms.CurrencyList {
		result.CurrencyMap[item.Name] = item.ToPb()
	}
	if len(result.GetCurrencyMap()) <= 0 {
		return nil, hcode.ErrCurrencyNotFound
	}
	return result, nil
}

// NetGetByCy 获取转账网络信息
func (cms CurrencyMangerSrv) NetGetByCy(ctx context.Context, req *dto.NameReq) (*dto.NetList, error) {
	var result = new(dto.NetList)
	result.Nets = make([]*dto.Net, 0)
	for _, item := range *cms.NetList {
		if item.Name == req.Name {
			result.Nets = append(result.Nets, item.ToPb())
		}
	}
	if len(result.GetNets()) <= 0 {
		return nil, hcode.ErrCurrencyNotFound
	}
	return result, nil
}
func (cms CurrencyMangerSrv) NetGetByName(ctx context.Context, req *dto.NameReq) (*dto.Net, error) {
	for _, item := range *cms.NetList {
		if item.Name == req.Name {
			return item.ToPb(), nil
		}
	}
	return nil, hcode.ErrNetNotFound
}
