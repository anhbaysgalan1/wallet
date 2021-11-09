package adapter

import (
	"context"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/leaf-rain/wallet/internal/account/entity"
	"github.com/leaf-rain/wallet/internal/block_chain/eth"
	"github.com/leaf-rain/wallet/pkg/log"
	"go.uber.org/zap"
	"math/big"
)

type EthService interface {
	// GetEtHAddrAndPrivateByNum 获取新创建帐号及私钥
	GetEtHAddrAndPrivateByNum(ctx context.Context, num int) (resp []*entity.EntityAddressPrivate)
	// GetNewBlockHeight 获取最新块高
	GetNewBlockHeight(ctx context.Context) (uint64, error)
	// GetDataByBlockHeight 获取块高数据通过块高
	GetDataByBlockHeight(ctx context.Context, height *big.Int) (*types.Block, error)
	// GetDataByBlockHash 获取区块数据通过区块hash
	GetDataByBlockHash(ctx context.Context, hash ethCommon.Hash) (*types.Block, error)
	// GetBalanceByHeight 获取用户余额
	GetBalanceByHeight(ctx context.Context, account ethCommon.Address, height *big.Int) (*big.Int, error)
	// Close 关闭连接
	Close()
}

func NewEth(url, appid, currencyId string) (EthService, error) {
	client, err := ethclient.Dial(url + appid)
	if err != nil {
		log.GetLogger().Error("[NewEth] failed", zap.Error(err))
		return &eth.EthClient{}, err
	}
	return &eth.EthClient{EthClient: client, CurrencyId: currencyId}, err
}
