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

func (srv *WalletSrv) JobNftWalletTransferToBlock(ctx context.Context, in *walletPb.Empty) (*walletPb.Empty, error) {
	var bills []*entity.Bill
	var err error
	// 获取订单
	bills, err = srv.Repo.BillNftGetForAsync(ctx)
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
			log.GetLogger().Error("[JobNftWalletTransferToBlock] Lock.TryLock failed",
				zap.String("key", lockKey),
				zap.Any("ttl", common.KeyQueueBillToPendingTtl),
				zap.Error(err))
			continue
		}
		if !lockResult {
			log.GetLogger().Warn("[JobNftWalletTransferToBlock] Lock.TryLock lockResult is false",
				zap.String("key", lockKey),
				zap.Any("ttl", common.KeyQueueBillToPendingTtl))
			continue
		}
		// 获取地址&nonce
		sysAddr, nonce, err = srv.Repo.GetTpSysExpendAddr(ctx, item.ContractRecord.ContractType)
		if err != nil {
			_ = srv.Lock.UnLock(ctx, lockKey, 0)
			return nil, err
		}
		var RacingBoatWithdrawalNftReq = transfer.InputForTransfer{
			FromAddress:     sysAddr,
			ContractAddress: config.SysAccountMap[item.ContractRecord.ContractType].ContractAddress,
			ToAddress:       item.ToAddr,
			GasLimit:        config.Fee.GasLimit,
			GasPrice:        config.Fee.GasPrice,
			Nonce:           nonce,
			Private:         "61a9ab8490e64b70b3c9bc59c8bfba47eb6bd208eb15e47807f76966fd160959",
		}
		var nonceIncr int64 = 1
		var nftToken = new(big.Int)
		nftToken.SetString(item.ContractRecord.NftToken, 0)
		for times := 1; times < Retry; times++ { // 重试机制
			hash = ""
			switch item.BillType {
			case entity.BillType_Eip721:
				hash, err = srv.BlockChainSrv.RacingBoatWithdrawalNft(ctx, RacingBoatWithdrawalNftReq, nftToken, config.BlockBusiness.KeyTransfer)
			case entity.BillType_Eip1155:
				var num = new(big.Int)
				num.SetUint64(item.ContractRecord.Num)
				hash, err = srv.BlockChainSrv.MaterialWithdrawal(ctx, RacingBoatWithdrawalNftReq, nftToken, num, config.BlockBusiness.KeyTransfer)
			default:
				log.GetLogger().Error("[JobWalletTransferToBlock] bill type error", zap.Any("bill", item))
				_ = srv.Lock.UnLock(ctx, lockKey, 0)
				break
			}
			if err != nil { // 错误继续重试
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
						log.GetLogger().Error("[JobWalletTransferToBlock] WalletSrv.BlockChainSrv.F1TransferCoin failed",
							zap.String("sys addr", sysAddr),
							zap.Any("bill", item),
							zap.Error(err))
					}
					time.Sleep(time.Second * 3) // 3秒后重试
					runtime.Gosched()
					continue
				}
			} else { // 成功直接返回
				item.FromAddr = sysAddr
				item.Hash = hash
				log.GetLogger().Info("[JobWalletTransferToBlock] success",
					zap.String("sys addr", sysAddr),
					zap.Any("bill", item))
				break
			}
		}
		if len(hash) > 0 { // 订单上链中
			item.BillStatus = walletPb.BillStatus_Pending
			_ = srv.Repo.BillDealWithPending(ctx, item)
			// nonce 值 增长
			_ = srv.Repo.NonceIncr(ctx, item.FromAddr, nonceIncr)
		} else { // 订单失败
			item.BillStatus = walletPb.BillStatus_Failed
			item.Remark = err.Error()
			_ = srv.Repo.BillDealWithFailed(ctx, item)
		}
		// 给转账地址解锁
		srv.Repo.UnlockNonceAddr(ctx, item.FromAddr)
		_ = srv.Lock.UnLock(ctx, lockKey, 0)
	}
	return nil, nil
}
