package adapter

import (
	"errors"
	"tp_wallet/config"
	"tp_wallet/internal/block_chain/dto"
	"tp_wallet/internal/block_chain/service"
	"tp_wallet/internal/common"
)

func NewBlockChainSrv() dto.BlockChainSrv {
	config.NewBlockBusiness()
	switch config.BlockBusiness.NetType {
	case common.ChainNetWorkTest:
		return &service.BlockChainSrv{
			NetWork: common.ChainNetWorkTest,
		}
	case common.ChainNetWorkRelease:
		return &service.BlockChainSrv{
			NetWork: common.ChainNetWorkRelease,
		}
	default:
		panic(errors.New("NetType unknown"))
	}
}
