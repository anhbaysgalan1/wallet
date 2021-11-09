package eth

import (
	"context"
	"crypto/ecdsa"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/leaf-rain/wallet/internal/account/dto"
	"github.com/leaf-rain/wallet/internal/account/entity"
	"github.com/leaf-rain/wallet/pkg/log"
	"github.com/leaf-rain/wallet/pkg/tool"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"math/big"
)

type EthClient struct {
	EthClient  *ethclient.Client
	CurrencyId string
}

func (e *EthClient) Close() {
	e.EthClient.Close()
}

func (e EthClient) ethGetAddrAndPrivate() (resp *entity.EntityAddressPrivate) {
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
	return &entity.EntityAddressPrivate{
		Id:         primitive.NewObjectID(),
		Address:    address,
		Private:    private,
		Currency:   e.CurrencyId,
		Status:     dto.AccountType_AccountType_UnUsed,
		CreateTime: tool.GetTimeUnixMilli(),
		UpdateTime: tool.GetTimeUnixMilli(),
		Remarks:    "",
	}
}

// GetEtHAddrAndPrivateByNum 获取新创建帐号及私钥
func (e *EthClient) GetEtHAddrAndPrivateByNum(ctx context.Context, num int) (resp []*entity.EntityAddressPrivate) {
	resp = make([]*entity.EntityAddressPrivate, num)
	for i := 0; i < num; i++ {
		resp[i] = e.ethGetAddrAndPrivate()
		if len(resp[i].Address) == 0 || len(resp[i].Private) == 0 {
			i--
			break
		}
	}
	return resp
}

// GetNewBlockHeight 获取最新块高
func (a EthClient) GetNewBlockHeight(ctx context.Context) (uint64, error) {
	return a.EthClient.BlockNumber(ctx)
}

// GetDataByBlockHeight 获取块高数据通过块高
func (a EthClient) GetDataByBlockHeight(ctx context.Context, height *big.Int) (*types.Block, error) {
	return a.EthClient.BlockByNumber(ctx, height)
}

// GetDataByBlockHash 获取区块数据通过区块hash
func (a EthClient) GetDataByBlockHash(ctx context.Context, hash ethCommon.Hash) (*types.Block, error) {
	return a.EthClient.BlockByHash(ctx, hash)
}

// GetBalanceByHeight 获取用户余额
func (a EthClient) GetBalanceByHeight(ctx context.Context, account ethCommon.Address, height *big.Int) (*big.Int, error) {
	return a.EthClient.BalanceAt(ctx, account, height)
}
