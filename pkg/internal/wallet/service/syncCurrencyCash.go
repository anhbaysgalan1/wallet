package service

import (
	"context"
	"errors"
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.uber.org/zap"
	"math/big"
	"runtime"
	"time"
	"tp_wallet/config"
	"tp_wallet/internal/block_chain/chain/transfer"
	"tp_wallet/internal/common"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/log"
)

func (srv *WalletSrv) JobCurrencyWalletTransferToBlock(ctx context.Context, in *walletPb.Empty) (*walletPb.Empty, error) {
	var bills []*entity.Bill
	var err error
	// 获取订单
	bills, err = srv.Repo.BillCurrencyGetForAsync(ctx)
	if err != nil {
		return nil, err
	}
	// 处理订单
	var Retry = 5 // 重试次数
	for _, item := range bills {
		var lockKey = common.KeyQueueBillToPending(item.Id.Hex())
		var hash string
		var lockResult bool
		var sysAddr string
		var nonce uint64
		lockResult, _, err = srv.Lock.TryLock(lockKey, 0, common.KeyQueueBillToPendingTtl, ctx)
		if err != nil {
			log.GetLogger().Error("[JobH2OWalletTransferToBlock] Lock.TryLock failed",
				zap.String("key", lockKey),
				zap.Any("ttl", common.KeyQueueBillToPendingTtl),
				zap.Error(err))
			continue
		}
		if !lockResult {
			log.GetLogger().Warn("[JobH2OWalletTransferToBlock] Lock.TryLock lockResult is false",
				zap.String("key", lockKey),
				zap.Any("ttl", common.KeyQueueBillToPendingTtl))
			continue
		}
		// 获取地址&nonce
		sysAddr, nonce, err = srv.Repo.GetCurrencySysExpendAddr(ctx, item.BalanceRecord.Currency)
		if err != nil {
			return nil, err
		}
		var amount = new(big.Int)
		amount.SetString(item.BalanceRecord.Amount, 0)
		var WithdrawalErc20Req = transfer.InputForTransfer{
			FromAddress:     sysAddr,
			ToAddress:       item.ToAddr,
			ContractAddress: config.SysAccountMap[item.BalanceRecord.Currency].ContractAddress,
			Amount:          amount,
			GasLimit:        config.Fee.GasLimit,
			GasPrice:        config.Fee.GasPrice,
			Nonce:           nonce,
			Private:         "ba35c41591cf77ba661b756401716ba1dde51ef5349f65d09b8957481a14d0d0",
		}
		var nonceIncr int64 = 1
		for times := 1; times < Retry; times++ { // 重试机制
			hash = ""
			hash, err = srv.BlockChainSrv.WithdrawalCurrency(ctx, WithdrawalErc20Req, config.BlockBusiness.KeyTransfer, item.BalanceRecord.Currency) // 这里是给from地址打钱
			if err != nil {                                                                                                                          // 错误继续重试
				switch {
				case errors.Is(err, hcode.ErrNonceTooHigh):
					nonceIncr += 1
					times -= 1
					continue
				case errors.Is(err, hcode.ErrNonceTooLow):
					nonceIncr -= 1
					times -= 1
					continue
				default:
					if times == Retry {
						log.GetLogger().Error("[JobH2OWalletTransferToBlock] WalletSrv.BlockChainSrv.F1TransferCoin failed",
							zap.String("sys addr", sysAddr),
							zap.Any("bill", item),
							zap.Error(err))
					}
					time.Sleep(time.Second)
					runtime.Gosched()
					continue
				}
			} else { // 成功直接返回
				log.GetLogger().Info("[JobH2OWalletTransferToBlock] success",
					zap.String("sys addr", sysAddr),
					zap.Any("bill", item))
				break
			}
		}
		if len(hash) > 0 { // 订单上链中
			item.FromAddr = sysAddr
			item.Hash = hash
			item.BillStatus = walletPb.BillStatus_Pending
			_ = srv.Repo.BillDealWithPending(ctx, item)
			// 给转账地址解锁
			srv.Repo.UnlockNonceAddr(ctx, item.FromAddr)
			// nonce 值 增长
			_ = srv.Repo.NonceIncr(ctx, item.FromAddr, nonceIncr)
		} else { // 订单失败
			item.BillStatus = walletPb.BillStatus_Failed
			item.Remark = err.Error()
			_ = srv.Repo.BillDealWithFailed(ctx, item)
		}
		_ = srv.Lock.UnLock(ctx, lockKey, 0)
	}
	return nil, nil
}
