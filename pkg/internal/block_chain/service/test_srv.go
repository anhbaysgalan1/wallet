package service

import (
	"context"
	"math/big"
	"tp_wallet/internal/block_chain/chain/transfer"
	"tp_wallet/pkg/tool"
)

type BlockTestSrv struct{}

func (BlockTestSrv) RacingBoatCreateNft(ctx context.Context, _transfer transfer.InputForTransfer, _propsName string, _starRating uint8, _key string) (hash string, err error) {
	return "test_hash" + tool.Int64ToStr(tool.GetTimeUnixMilli()), nil
}

func (BlockTestSrv) RacingBoatWithdrawalNft(ctx context.Context, _transfer transfer.InputForTransfer, _NftID *big.Int, _key string) (hash string, err error) {
	return "test_hash" + tool.Int64ToStr(tool.GetTimeUnixMilli()), nil
}

func (BlockTestSrv) WithdrawalErc20(ctx context.Context, _transfer transfer.InputForTransfer, _key string) (hash string, err error) {
	return "test_hash" + tool.Int64ToStr(tool.GetTimeUnixMilli()), nil
}

func (BlockTestSrv) WithdrawalFFCoin(ctx context.Context, _transfer transfer.InputForTransfer, _key string) (hash string, err error) {
	return "test_hash" + tool.Int64ToStr(tool.GetTimeUnixMilli()), nil
}

func (BlockTestSrv) WithdrawalF1Coin(ctx context.Context, _transfer transfer.InputForTransfer, _key string) (hash string, err error) {
	return "test_hash" + tool.Int64ToStr(tool.GetTimeUnixMilli()), nil
}

func (BlockTestSrv) WithdrawalBNB(ctx context.Context, _transfer transfer.InputForBNBTransfer, _key string) (hash string, err error) {
	return "test_hash" + tool.Int64ToStr(tool.GetTimeUnixMilli()), nil
}

func (BlockTestSrv) GetAddressNonceForBsc(ctx context.Context, _addr string, _key string) (nonce uint64, err error) {
	return 1, nil
}

func (BlockTestSrv) GetWithdrawalFee(ctx context.Context) (fee *big.Int, err error) {
	return new(big.Int).SetInt64(21000), nil
}

// MaterialWithdrawal 提现材料NFT
func (BlockTestSrv) MaterialWithdrawal(
	ctx context.Context,
	_transfer transfer.InputForTransfer,
	_NftID *big.Int, //NFT唯一链上ID
	_amount *big.Int, //材料数量
	_key string, //api key
) (hash string, err error) {
	return "test_hash" + tool.Int64ToStr(tool.GetTimeUnixMilli()), nil
}

// MaterialExpend 消耗材料NFT
func (BlockTestSrv) MaterialExpend(
	ctx context.Context,
	_transfer transfer.InputForTransfer,
	_NftID *big.Int, //NFT唯一链上ID
	_amount *big.Int, //材料数量
	_key string, //api key
) (hash string, err error) {
	return "test_hash" + tool.Int64ToStr(tool.GetTimeUnixMilli()), nil
}
