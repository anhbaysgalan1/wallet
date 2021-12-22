package service

import (
	"context"
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"strings"
	"tp_wallet/internal/wallet/adapter/block_chain"
	"tp_wallet/internal/wallet/adapter/props"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/internal/wallet/repository"
	redisCommon "tp_wallet/pkg/redisCache/common"
)

type WalletSrv struct {
	Repo          repository.Repository
	BlockChainSrv block_chain.BlockChainSrv
	Lock          redisCommon.RedisLock
	walletPb.UnimplementedWalletSrvServer
	PropsSrv props.PropsSrv
}

// DealWithBill 处理订单
func (srv WalletSrv) DealWithBill(ctx context.Context, in *walletPb.BillInfo) (*walletPb.Empty, error) {
	in.FromAddr = strings.ToLower(in.FromAddr)
	in.ToAddr = strings.ToLower(in.ToAddr)
	if in.ContractRecord != nil {
		in.ContractRecord.ContractAddr = strings.ToLower(in.ContractRecord.ContractAddr)
	}
	var bill = entity.PbToBill(in)
	if bill == nil {
		return nil, hcode.ErrInternalParameter
	}
	var err = srv.Repo.BillDealWith(ctx, bill)
	return nil, err
}
