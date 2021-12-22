package service

import (
	"context"
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.uber.org/zap"
	"tp_wallet/config"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/log"
)

// AccountGet 获取(创建)帐号 如果已经存在账户，则直接返回
func (srv WalletSrv) AccountGet(ctx context.Context, req *walletPb.UidReq) (*walletPb.AccountGetResp, error) {
	return srv.Repo.AccountGetAndCreate(ctx, req)
}

// BalanceGet 获取用户余额
func (srv WalletSrv) BalanceGet(ctx context.Context, req *walletPb.UidReq) (*walletPb.AccountGetResp, error) {
	var balance map[string]entity.Amount
	var err error
	balance, err = srv.Repo.BalanceGetByUid(ctx, req.GetUid())
	if err != nil {
		return nil, err
	}
	return &walletPb.AccountGetResp{Uid: req.Uid, Balance: entity.AmountMapToString(balance), IsExist: true}, nil
}

// TransferCurrencyForOffline 离线转账
func (srv WalletSrv) TransferCurrencyForOffline(ctx context.Context, req *walletPb.TransferForOfflineReq) (*walletPb.Empty, error) {
	return srv.Repo.TransferCurrencyForOffline(ctx, req)
}

// TransferCurrencyCash 提现
func (srv WalletSrv) TransferCurrencyCash(ctx context.Context, req *walletPb.TransferCashReq) (*walletPb.Empty, error) {
	var err error
	var gas entity.Gas
	// 校验币种信息
	if _, ok := config.CurrencyMap[req.Currency]; !ok {
		return nil, hcode.ErrCurrencyUnsupported
	} else {
		var currency = config.CurrencyMap[req.Currency]
		if !currency.GetCashStatus() {
			return nil, hcode.ErrCurrencyUnsupported
		}
		// 是否大于最小提现金额
		var resultInt int
		resultInt, err = entity.Amount(req.GetAmount()).Cmp(currency.MinCash)
		if err != nil {
			return nil, err
		}
		if resultInt <= 0 {
			return nil, hcode.ErrLessFee
		}
		gas = entity.Gas{
			GasPrice:    currency.Net[currency.DefaultNet].GasPrice,
			GasLimit:    currency.Net[currency.DefaultNet].GasLimit,
			GasCurrency: currency.Net[currency.DefaultNet].GasCurrency,
		}
	}
	// 校验系统账号
	if _, ok := config.SysAccountMap[req.Currency]; !ok {
		return nil, hcode.ErrCurrencyUnsupported
	}
	bill := &entity.Bill{
		FromUid: config.SysAccountMap[req.Currency].SysUid,
		ToUid:   req.Uid,
		ToAddr:  req.GetToAddr(),
		Uid:     req.GetUid(),
		Gas:     gas,
	}
	bill.BalanceRecord = entity.BalanceRecord{
		Currency:      req.Currency,
		Amount:        req.Amount,
		ReceiveAmount: "",
		BeforeBalance: "",
		AfterBalance:  "",
	}
	return srv.Repo.TransferCurrencyCash(ctx, bill)
}

// GetSysTransferAddr 获取系统账户地址
func (srv WalletSrv) GetSysTransferAddr(ctx context.Context, in *walletPb.CurrencyReq) (*walletPb.AddrResp, error) {
	if currency, ok := config.CurrencyMap[in.GetCurrency()]; !ok {
		log.GetLogger().Error("[GetSysTransferAddr] config.CurrencyMap failed", zap.Any("req", in))
		return nil, hcode.ErrCurrencyUnsupported
	} else {
		if !currency.GetChargeStatus() {
			return nil, hcode.ErrCurrencyUnsupported
		}
	}
	if account, ok := config.SysAccountMap[in.GetCurrency()]; !ok {
		log.GetLogger().Error("[GetSysTransferAddr] config.SysAccountMap failed", zap.Any("req", in))
		return nil, hcode.ErrCurrencyUnsupported
	} else {
		return &walletPb.AddrResp{Addr: account.AddrIncome}, nil
	}
}
