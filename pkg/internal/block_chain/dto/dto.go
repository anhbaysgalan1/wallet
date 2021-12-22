package dto

import (
	"context"
	"math/big"
	"tp_wallet/internal/block_chain/chain/transfer"
)

type BlockChainSrv interface {
	// RacingBoatCreateNft 生产并下发nft到用户第三方钱包地址
	RacingBoatCreateNft(
		ctx context.Context,
		_transfer transfer.InputForTransfer,
		_propsName string, //赛艇名称+编号
		_starRating uint8, //赛艇星级
		_key string, //api key
	) (hash string, err error)

	// RacingBoatWithdrawalNft nft提现
	RacingBoatWithdrawalNft(
		ctx context.Context,
		_transfer transfer.InputForTransfer,
		_NftID *big.Int, //NFT唯一链上ID
		_key string, //api key
	) (hash string, err error)

	// MaterialCreate 创建一定数量的材料NFT
	MaterialCreate(
		ctx context.Context,
		_transfer transfer.InputForTransfer,
		_amount *big.Int, //材料数量
		_name string, //材料名称
		_key string, //api key
	) (hash string, err error)

	// MaterialWithdrawal 提现材料NFT
	MaterialWithdrawal(
		ctx context.Context,
		_transfer transfer.InputForTransfer,
		_NftID *big.Int, //NFT唯一链上ID
		_amount *big.Int, //材料数量
		_key string, //api key
	) (hash string, err error)

	// MaterialExpend 消耗材料NFT
	MaterialExpend(
		ctx context.Context,
		_transfer transfer.InputForTransfer,
		_NftID *big.Int, //NFT唯一链上ID
		_amount *big.Int, //材料数量
		_key string, //api key
	) (hash string, err error)

	// WithdrawalErc20 提现h20代币
	WithdrawalErc20(ctx context.Context, _transfer transfer.InputForTransfer, _key string) (hash string, err error)

	// WithdrawalFFCoin 提现ff代币
	WithdrawalFFCoin(ctx context.Context, _transfer transfer.InputForTransfer, _key string) (hash string, err error)

	// WithdrawalF1Coin 提现f1代币
	WithdrawalF1Coin(ctx context.Context, _transfer transfer.InputForTransfer, _key string) (hash string, err error)

	// WithdrawalBNB 提现BNB代币
	WithdrawalBNB(ctx context.Context, _transfer transfer.InputForBNBTransfer, _key string) (hash string, err error)

	// GetAddressNonceForBsc 获取地址nonce值
	GetAddressNonceForBsc(ctx context.Context, _addr string, _key string) (nonce uint64, err error)

	// GetWithdrawalFee 获取链上提现手续费，消耗的是bnb
	GetWithdrawalFee(ctx context.Context) (fee *big.Int, err error)
}
