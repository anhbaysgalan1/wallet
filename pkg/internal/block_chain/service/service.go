package service

import (
	"context"
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	"go.uber.org/zap"
	"math/big"
	"tp_wallet/internal/block_chain/chain/transfer"
	"tp_wallet/pkg/log"
)

type BlockChainSrv struct {
	NetWork string
}

// RacingBoatCreateNft 生产并下发nft到用户第三方钱包地址
func (srv BlockChainSrv) RacingBoatCreateNft(
	ctx context.Context,
	_transfer transfer.InputForTransfer,
	_propsName string, //赛艇名称+编号
	_starRating uint8, //赛艇星级
	_key string,
) (hash string, err error) {
	hash, _, err = transfer.CreateAssetRacingBoatForBsc(_transfer, srv.NetWork, _propsName, _starRating)
	if err != nil {
		log.GetLogger().Error("CreateAssetRacingBoatForBsc", zap.Error(err))
		if err.Error() == transfer.ErrNonceTooLow {
			return "", hcode.ErrNonceTooLow
		}

		if err.Error() == transfer.ErrInsufficientFunds {
			return "", hcode.ErrInsufficientFunds
		}
		return "", hcode.ErrCreateAssetRacingBoat
	}
	return hash, nil
}

// RacingBoatWithdrawalNft nft提现
func (srv BlockChainSrv) RacingBoatWithdrawalNft(
	ctx context.Context,
	_transfer transfer.InputForTransfer,
	_NftID *big.Int, //NFT唯一链上ID
	_key string,
) (hash string, err error) {
	hash, _, err = transfer.WithdrawalRacingBoatForBsc(_transfer, srv.NetWork, _NftID)
	if err != nil {
		log.GetLogger().Error("WithdrawalRacingBoatForBsc", zap.Error(err))
		if err.Error() == transfer.ErrNonceTooLow {
			return "", hcode.ErrNonceTooLow
		}

		if err.Error() == transfer.ErrInsufficientFunds {
			return "", hcode.ErrInsufficientFunds
		}
		return "", hcode.ErrWithdrawal721NFTRacingBoat
	}
	return hash, nil
}

// MaterialCreate 创建一定数量的材料NFT
func (srv BlockChainSrv) MaterialCreate(
	ctx context.Context,
	_transfer transfer.InputForTransfer,
	_amount *big.Int, //材料数量
	_name string, //材料名称
	_key string,
) (hash string, err error) {
	hash, _, err = transfer.CreateAssetMaterialForBsc(_transfer, srv.NetWork, _amount, _name)
	if err != nil {
		log.GetLogger().Error("CreateAssetMaterialForBsc", zap.Error(err))
		if err.Error() == transfer.ErrNonceTooLow {
			return "", hcode.ErrNonceTooLow
		}

		if err.Error() == transfer.ErrInsufficientFunds {
			return "", hcode.ErrInsufficientFunds
		}
		return "", hcode.ErrCreateMaterial
	}
	return hash, nil
}

// MaterialWithdrawal 提现材料NFT
func (srv BlockChainSrv) MaterialWithdrawal(
	ctx context.Context,
	_transfer transfer.InputForTransfer,
	_NftID *big.Int, //NFT唯一链上ID
	_amount *big.Int, //材料数量
	_key string,
) (hash string, err error) {
	hash, _, err = transfer.WithdrawalMaterialForBsc(_transfer, srv.NetWork, _NftID, _amount)
	if err != nil {
		log.GetLogger().Error("WithdrawalMaterialForBsc", zap.Error(err))
		if err.Error() == transfer.ErrNonceTooLow {
			return "", hcode.ErrNonceTooLow
		}

		if err.Error() == transfer.ErrInsufficientFunds {
			return "", hcode.ErrInsufficientFunds
		}
		return "", hcode.ErrWithdrawalMaterial
	}
	return hash, nil
}

// MaterialExpend 消耗材料NFT
func (srv BlockChainSrv) MaterialExpend(
	ctx context.Context,
	_transfer transfer.InputForTransfer,
	_NftID *big.Int, //NFT唯一链上ID
	_amount *big.Int, //材料数量
	_key string,
) (hash string, err error) {
	hash, _, err = transfer.ExpandMaterialForBsc(_transfer, srv.NetWork, _NftID, _amount)
	if err != nil {
		log.GetLogger().Error("ExpandMaterialForBsc", zap.Error(err))
		if err.Error() == transfer.ErrNonceTooLow {
			return "", hcode.ErrNonceTooLow
		}

		if err.Error() == transfer.ErrInsufficientFunds {
			return "", hcode.ErrInsufficientFunds
		}
		return "", hcode.ErrExpandMaterial
	}
	return hash, nil
}

// WithdrawalErc20 提现H2O代币
func (srv BlockChainSrv) WithdrawalErc20(ctx context.Context, _transfer transfer.InputForTransfer, _key string) (hash string, err error) {
	hash, _, err = transfer.Withdrawal20TokenForBsc(_transfer, srv.NetWork)
	if err != nil {
		log.GetLogger().Error("WithdrawalH2O", zap.Error(err))
		if err.Error() == transfer.ErrNonceTooLow {
			return "", hcode.ErrNonceTooLow
		}

		if err.Error() == transfer.ErrInsufficientFunds {
			return "", hcode.ErrInsufficientFunds
		}
		return "", hcode.ErrWithdrawalH2O
	}
	return hash, nil
}

// WithdrawalFFCoin 提现FF代币
func (srv BlockChainSrv) WithdrawalFFCoin(ctx context.Context, _transfer transfer.InputForTransfer, _key string) (hash string, err error) {
	hash, _, err = transfer.Withdrawal20TokenForBsc(_transfer, srv.NetWork)
	if err != nil {
		log.GetLogger().Error("WithdrawalFFCoin", zap.Error(err))
		if err.Error() == transfer.ErrNonceTooLow {
			return "", hcode.ErrNonceTooLow
		}

		if err.Error() == transfer.ErrInsufficientFunds {
			return "", hcode.ErrInsufficientFunds
		}
		return "", hcode.ErrWithdrawalFFCoin
	}
	return hash, nil
}

// WithdrawalF1Coin 提现F1代币
func (srv BlockChainSrv) WithdrawalF1Coin(ctx context.Context, _transfer transfer.InputForTransfer, _key string) (hash string, err error) {
	hash, _, err = transfer.Withdrawal20TokenForBsc(_transfer, srv.NetWork)
	if err != nil {
		log.GetLogger().Error("WithdrawalF1Coin", zap.Error(err))
		if err.Error() == transfer.ErrNonceTooLow {
			return "", hcode.ErrNonceTooLow
		}

		if err.Error() == transfer.ErrInsufficientFunds {
			return "", hcode.ErrInsufficientFunds
		}
		return "", hcode.ErrWithdrawalF1Coin
	}
	return hash, nil
}

// WithdrawalBNB 提现BNB
func (srv BlockChainSrv) WithdrawalBNB(ctx context.Context, _transfer transfer.InputForBNBTransfer, _key string) (hash string, err error) {
	hash, _, err = transfer.WithdrawalBNBForBsc(_transfer, srv.NetWork)
	if err != nil {
		log.GetLogger().Error("WithdrawalBNB", zap.Error(err))
		if err.Error() == transfer.ErrNonceTooLow {
			return "", hcode.ErrNonceTooLow
		}

		if err.Error() == transfer.ErrInsufficientFunds {
			return "", hcode.ErrInsufficientFunds
		}
		return "", hcode.ErrWithdrawalBNB
	}
	return hash, nil
}

// GetAddressNonceForBsc 获取地址nonce值
func (srv BlockChainSrv) GetAddressNonceForBsc(ctx context.Context, _addr string, _key string) (nonce uint64, err error) {
	nonce, err = transfer.GetAccountNonce(_addr, srv.NetWork)
	if err != nil {
		log.GetLogger().Error("GetAccountNonce", zap.Error(err))
		return 0, hcode.ErrGetAddressNonce
	}
	return nonce, nil
}

// GetWithdrawalFee 获取链上提现手续费，消耗的是bnb
func (srv BlockChainSrv) GetWithdrawalFee(ctx context.Context) (fee *big.Int, err error) {
	return new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetInt64(200000*5)), nil
}
