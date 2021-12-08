package service

import (
	"github.com/leaf-rain/wallet/internal/account/dto"
	"github.com/leaf-rain/wallet/internal/account/repository"
)

type AccountSrv struct {
	repo repository.AccountRepo
	dto.UnimplementedAccountSrvServer
}

// 获取(创建)帐号 如果已经存在账户，则直接返回
//func (srv AccountSrv) AccountGet(ctx context.Context, req *dto.AccountGetReq) (*dto.AddressResp, error) {
//
//}

// 获取账户转账地址
//func (srv AccountSrv) AddressGet(ctx context.Context, req *dto.AddressGetReq) (*dto.AddressResp, error) {
//
//}
