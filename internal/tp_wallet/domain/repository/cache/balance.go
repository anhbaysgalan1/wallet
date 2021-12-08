package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"tp_wallet/internal/common"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/hcode"
	"tp_wallet/pkg/log"
)

func (c *walletCache) BalanceSave(ctx context.Context, balance []*entity.Balance) error {
	if len(balance) == 0 {
		log.GetLogger().Error("[BalanceSave] length addr zero", zap.Any("balance", balance))
		return nil
	}
	var uid = balance[0].Uid
	var value = make(map[string]interface{})
	for _, item := range balance {
		if uid != item.Uid {
			log.GetLogger().Error("[BalanceSave] uid different", zap.Any("balance", balance))
			return hcode.ErrInternalParameter
		}
		value[item.Currency] = item.Balance
	}
	err := c.cache.HMSet(common.KeyWalletSysBalance(uid), value).Err()
	if err != nil {
		log.GetLogger().Error("[BalanceSave] cache.HMSet failed",
			zap.Any("balance", balance),
			zap.Error(err))
		return hcode.ErrInternalCache
	}
	return nil
}

func (c *walletCache) BalanceDel(ctx context.Context, uid uint64) error {
	err := c.cache.Del(common.KeyWalletSysBalance(uid)).Err()
	if err != nil {
		log.GetLogger().Error("[BalanceDel] cache.Del failed",
			zap.Uint64("uid", uid),
			zap.Error(err))
		return hcode.ErrInternalCache
	}
	return err
}

func (c *walletCache) BalanceGetByUid(ctx context.Context, uid uint64) (map[string]entity.Amount, error) {
	var balance = make(map[string]entity.Amount)
	var result = make(map[string]string)
	var err error
	result, err = c.cache.HGetAll(common.KeyWalletSysBalance(uid)).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.GetLogger().Error("[BalanceGetByUid] cache.HGetAll failed",
				zap.Uint64("uid", uid),
				zap.Error(err))
			return balance, hcode.ErrInternalCache
		}
	}
	for k, v := range result {
		balance[k] = entity.Amount(v)
	}
	return balance, nil
}

// BalanceSysSet 防止系统用户因为锁和数据库的原因并发不足，所以纯redis操作
func (c *walletCache) BalanceSysSet(ctx context.Context, sysUid uint64, amount, currency string, isAdd bool) (entity.Amount, entity.Amount, error) {
	if sysUid > common.UidStart {
		log.GetLogger().Error("[BalanceSysSet] uid not sys",
			zap.Uint64("sysUid", sysUid),
			zap.String("amount", amount),
			zap.String("currency", currency),
			zap.Bool("IsAdd", isAdd))
		return "", "", hcode.ErrInternalParameter
	}
	// 加锁
	var err error
	var lockResult bool
	var lockKey = common.KeyLockAccountBalance(sysUid)
	lockResult, _, err = c.lock.TryLock(lockKey, 0, common.LockAccountBalanceTtl, ctx)
	if err != nil {
		log.GetLogger().Error("[BalanceSysSet] Lock.TryLock failed",
			zap.Uint64("sysUid", sysUid),
			zap.String("amount", amount),
			zap.Error(err))
		return "", "", hcode.ErrServer
	}
	if !lockResult {
		return "", "", hcode.ErrReqLimit
	}
	defer c.lock.UnLock(ctx, lockKey, 0)
	// 获取余额
	var afterBalance entity.Amount
	var sysBalance map[string]entity.Amount
	sysBalance, err = c.BalanceGetByUid(ctx, sysUid)
	if err != nil {
		return "", "", err
	}
	if _, ok := sysBalance[currency]; !ok {
		sysBalance[currency] = "0"
	}
	afterBalance = sysBalance[currency]
	// IsAdd 用来判断是加还是减
	// 这里没有判断系统用户账户余额是否充足
	if isAdd {
		err = afterBalance.Add(amount)
		if err != nil {
			return "", "", err
		}
	} else {
		var result bool
		result, err = afterBalance.Sub(amount)
		if err != nil {
			return "", "", err
		}
		if !result {
			log.GetLogger().Error("[BalanceSysSet] sys balance deficiency", zap.Uint64("sys uid", sysUid), zap.Any("balance", sysBalance[currency]), zap.String("amount", amount))
			return "", "", hcode.ErrSysBalanceInsufficient
		}
	}
	// 修改系统余额
	err = c.cache.HSet(common.KeyWalletSysBalance(sysUid), currency, afterBalance).Err()
	if err != nil {
		log.GetLogger().Error("[BalanceSysSet] cache.Set failed",
			zap.Uint64("uid", sysUid),
			zap.String("amount", amount),
			zap.String("currency", currency),
			zap.String("amount", amount),
			zap.Bool("IsAdd", isAdd),
			zap.Error(err))
		return "", "", err
	}
	log.GetLogger().Info("[BalanceSysSet] success",
		zap.Bool("isAdd", isAdd),
		zap.Uint64("sys uid", sysUid),
		zap.Any("beforeBalance", sysBalance[currency]),
		zap.Any("afterBalance", afterBalance))
	return sysBalance[currency], afterBalance, nil
}
