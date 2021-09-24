package eth

import (
	"context"
	"crypto/ecdsa"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/leaf-rain/wallet/common"
	"github.com/leaf-rain/wallet/internal/block_chain/model"
	"github.com/leaf-rain/wallet/pkg/log"
	"github.com/leaf-rain/wallet/pkg/tool"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"math/big"
)

func ethGetAddrAndPrivate() (resp model.AddressPrivate) {
	// 生成私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.GetLogger().Error("[ethGetAddrAndPrivate] private failed,", zap.Error(err))
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	private := hexutil.Encode(privateKeyBytes)
	// 私钥导出地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.GetLogger().Error("[ethGetAddrAndPrivate] addr failed,", zap.Error(err))
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return model.AddressPrivate{
		Id:         primitive.NewObjectID(),
		Address:    address,
		Private:    private,
		Currency:   common.Currency_ETH.String(),
		Net:        common.Net_ETH_Miner.String(),
		Status:     model.AddressStatus_UnUsed,
		CreateTime: tool.GetTimeUnixMilli(),
		UpdateTime: tool.GetTimeUnixMilli(),
		Remarks:    "",
	}
}

func (a *eth) GetEtHAddrAndPrivateByNum(ctx context.Context, num int) (resp []model.AddressPrivate) {
	resp = make([]model.AddressPrivate, num)
	for i := 0; i < num; i++ {
		resp[i] = ethGetAddrAndPrivate()
		if len(resp[i].Address) == 0 || len(resp[i].Private) == 0 {
			i--
			break
		}

	}
	return resp
}

// GetNewBlockHeight 获取最新块高
func (a eth) GetNewBlockHeight(ctx context.Context) (uint64, error) {
	return a.ethClient.BlockNumber(ctx)
}

// GetDataByBlockHeight 获取块高数据通过块高
func (a eth) GetDataByBlockHeight(ctx context.Context, height *big.Int) (*types.Block, error) {
	return a.ethClient.BlockByNumber(ctx, height)
}

// GetDataByBlockHash 获取区块数据通过区块hash
func (a eth) GetDataByBlockHash(ctx context.Context, hash ethCommon.Hash) (*types.Block, error) {
	return a.ethClient.BlockByHash(ctx, hash)
}

// GetBalanceByHeight 获取用户余额
func (a eth) GetBalanceByHeight(ctx context.Context, account ethCommon.Address, height *big.Int) (*big.Int, error) {
	return a.ethClient.BalanceAt(ctx, account, height)
}
