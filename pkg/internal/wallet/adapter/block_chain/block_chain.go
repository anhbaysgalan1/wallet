package block_chain

import (
	"context"
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	"math/big"
	"tp_wallet/internal/block_chain/chain/transfer"
	"tp_wallet/internal/block_chain/dto"
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

	// RacingBoatWithdrawalNft 提现非同质化代币(nft)
	RacingBoatWithdrawalNft(
		ctx context.Context,
		_transfer transfer.InputForTransfer,
		_NftID *big.Int, //NFT唯一链上ID
		_key string, //api key
	) (hash string, err error)

	// WithdrawalCurrency 提现同质化代币
	WithdrawalCurrency(ctx context.Context, _transfer transfer.InputForTransfer, _key string, currency string) (hash string, err error)

	// GetAddressNonceForBsc 获取地址nonce值
	GetAddressNonceForBsc(ctx context.Context, _addr string, _key string) (nonce uint64, err error)

	// GetWithdrawalFee 获取链上提现手续费，消耗的是bnb
	GetWithdrawalFee(ctx context.Context) (fee *big.Int, err error)

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

	// MaterialCreate 创建一定数量的材料NFT
	MaterialCreate(
		ctx context.Context,
		_transfer transfer.InputForTransfer,
		_amount *big.Int, //材料数量
		_name string, //材料名称
		_key string, //api key
	) (hash string, err error)
}

type blockChainSrv struct {
	srv dto.BlockChainSrv
}

func (b blockChainSrv) RacingBoatCreateNft(ctx context.Context, _transfer transfer.InputForTransfer, _propsName string, _starRating uint8, _key string) (hash string, err error) {
	return b.srv.RacingBoatCreateNft(ctx, _transfer, _propsName, _starRating, _key)
}

func (b blockChainSrv) RacingBoatWithdrawalNft(ctx context.Context, _transfer transfer.InputForTransfer, _NftID *big.Int, _key string) (hash string, err error) {
	return b.srv.RacingBoatWithdrawalNft(ctx, _transfer, _NftID, _key)
}

func (b blockChainSrv) WithdrawalCurrency(ctx context.Context, _transfer transfer.InputForTransfer, _key string, currency string) (hash string, err error) {
	switch currency {
	case "h2o":
		return b.srv.WithdrawalErc20(ctx, _transfer, _key)
	case "ff":
		return b.srv.WithdrawalFFCoin(ctx, _transfer, _key)
	case "f1":
		return b.srv.WithdrawalF1Coin(ctx, _transfer, _key)
	case "bnb":
		var req = transfer.InputForBNBTransfer{
			FromAddress: _transfer.FromAddress,
			ToAddress:   _transfer.ToAddress,
			Amount:      _transfer.Amount,
			GasLimit:    _transfer.GasLimit,
			GasPrice:    _transfer.GasPrice,
			Nonce:       _transfer.Nonce,
			Private:     _transfer.Private,
		}
		return b.srv.WithdrawalBNB(ctx, req, _key)
	default:
		return "", hcode.ErrCurrencyUnsupported
	}
}

func (b blockChainSrv) GetAddressNonceForBsc(ctx context.Context, _addr string, _key string) (nonce uint64, err error) {
	return b.srv.GetAddressNonceForBsc(ctx, _addr, _key)
}

func (b blockChainSrv) GetWithdrawalFee(ctx context.Context) (fee *big.Int, err error) {
	return b.srv.GetWithdrawalFee(ctx)
}

// MaterialWithdrawal 提现材料NFT
func (b blockChainSrv) MaterialWithdrawal(
	ctx context.Context,
	_transfer transfer.InputForTransfer,
	_NftID *big.Int, //NFT唯一链上ID
	_amount *big.Int, //材料数量
	_key string, //api key
) (hash string, err error) {
	return b.srv.MaterialWithdrawal(ctx, _transfer, _NftID, _amount, _key)
}

// MaterialExpend 消耗材料NFT
func (b blockChainSrv) MaterialExpend(
	ctx context.Context,
	_transfer transfer.InputForTransfer,
	_NftID *big.Int, //NFT唯一链上ID
	_amount *big.Int, //材料数量
	_key string, //api key
) (hash string, err error) {
	return b.srv.MaterialExpend(ctx, _transfer, _NftID, _amount, _key)
}

// MaterialCreate 创建一定数量的材料NFT
func (b blockChainSrv) MaterialCreate(
	ctx context.Context,
	_transfer transfer.InputForTransfer,
	_amount *big.Int, //材料数量
	_name string, //材料名称
	_key string, //api key
) (hash string, err error) {
	return b.srv.MaterialCreate(ctx, _transfer, _amount, _name, _key)
}

func NewBlockChainSrv(srv dto.BlockChainSrv) BlockChainSrv {
	return &blockChainSrv{srv: srv}
}
