package repository

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"tp_wallet/config"
	"tp_wallet/internal/common"
	"tp_wallet/pkg/hcode"
	"tp_wallet/pkg/log"
)

func lockSysAddrForTransfer(str string) string {
	return "tp_wallet:transfer_" + str
}

func (repo RepositoryStruct) NonceGetByAddr(ctx context.Context, addr string) (uint64, error) {
	var nonce uint64
	var err error
	nonce, err = repo.Cache.NonceGetByAddr(ctx, addr)
	if err != nil {
		if errors.Is(err, hcode.ErrGetAddressNonce) { //
			nonce, err = repo.BlockChainSrv.GetAddressNonceForBsc(ctx, addr, config.BlockBusiness.KeyTransfer)
			if err != nil {
				return 0, err
			} else {
				_ = repo.NonceIncr(ctx, addr, int64(nonce))
				return nonce, nil
			}
		}
		return 0, err
	}
	return nonce, nil
}

func (repo RepositoryStruct) NonceIncr(ctx context.Context, addr string, incr int64) error {
	return repo.Cache.NonceIncr(ctx, addr, incr)
}

func (repo RepositoryStruct) GetAndLockAddr(ctx context.Context, addr string) (uint64, error) {
	var lockKey string
	var lockResult bool
	var err error
	var nonce uint64
	lockKey = lockSysAddrForTransfer(addr)
	lockResult, _, err = repo.Lock.TryLock(lockKey, 0, common.NonceAddrTtl, ctx)
	if err != nil {
		log.GetLogger().Error("[GetF1H2OSysAddr] lock.TryLock failed", zap.String("lockKey", lockKey), zap.Any("ttl", common.NonceAddrTtl), zap.Error(err))
		return 0, hcode.ErrInternalCache
	}
	if !lockResult {
		return 0, hcode.ErrReqLimit
	}
	// 获取nonce值
	nonce, err = repo.NonceGetByAddr(ctx, addr)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

// GetH2OSysExpendAddr 获取h20系统支出交易地址  返回地址，锁
func (repo RepositoryStruct) GetH2OSysExpendAddr(ctx context.Context) (string, uint64, error) {
	var lockKey string
	var lockResult bool
	var err error
	var nonce uint64
	for _, item := range config.BlockBusiness.H2OUidAddrExpenditure {
		lockKey = lockSysAddrForTransfer(item)
		lockResult, _, err = repo.Lock.TryLock(lockKey, 0, common.NonceAddrTtl, ctx)
		if err != nil {
			log.GetLogger().Error("[GetF1H2OSysAddr] lock.TryLock failed", zap.String("lockKey", lockKey), zap.Any("ttl", common.NonceAddrTtl), zap.Error(err))
			return "", 0, hcode.ErrInternalCache
		}
		if lockResult { // 加锁成功，返回地址
			// 获取nonce值
			nonce, err = repo.NonceGetByAddr(ctx, item)
			if err != nil {
				return "", 0, err
			}
			return item, nonce, nil
		} else {
			continue
		}
	}
	return "", 0, hcode.ErrReqLimit
}

// GetNftSysExpendAddr 获取Nft系统支出交易地址  返回地址，锁
func (repo RepositoryStruct) GetNftSysExpendAddr(ctx context.Context) (string, uint64, error) {
	var lockKey string
	var lockResult bool
	var err error
	var nonce uint64
	for _, item := range config.BlockBusiness.NftAddrExpenditure {
		lockKey = lockSysAddrForTransfer(item)
		lockResult, _, err = repo.Lock.TryLock(lockKey, 0, common.NonceAddrTtl, ctx)
		if err != nil {
			log.GetLogger().Error("[GetNftSysExpendAddr] lock.TryLock failed", zap.String("lockKey", lockKey), zap.Any("ttl", common.NonceAddrTtl), zap.Error(err))
			return "", 0, hcode.ErrInternalCache
		}
		if lockResult { // 加锁成功，返回地址
			// 获取nonce值
			nonce, err = repo.NonceGetByAddr(ctx, item)
			if err != nil {
				return "", 0, err
			}
			return item, nonce, nil
		} else {
			continue
		}
	}
	return "", 0, hcode.ErrReqLimit
}

// GetTpSysExpendAddr 获取第三方合约系统系统支出交易地址  返回地址，锁
func (repo RepositoryStruct) GetTpSysExpendAddr(ctx context.Context, tp string) (string, uint64, error) {
	var lockKey string
	var lockResult bool
	var err error
	var nonce uint64
	if _, ok := config.SysAccountMap[tp]; !ok {
		return "", 0, hcode.ErrSysAccountNotFound
	}
	for _, item := range config.SysAccountMap[tp].AddrExpenditure {
		lockKey = lockSysAddrForTransfer(item)
		lockResult, _, err = repo.Lock.TryLock(lockKey, 0, common.NonceAddrTtl, ctx)
		if err != nil {
			log.GetLogger().Error("[GetTpSysExpendAddr] lock.TryLock failed", zap.String("lockKey", lockKey), zap.Any("ttl", common.NonceAddrTtl), zap.Error(err))
			return "", 0, hcode.ErrInternalCache
		}
		if lockResult { // 加锁成功，返回地址
			// 获取nonce值
			nonce, err = repo.NonceGetByAddr(ctx, item)
			if err != nil {
				return "", 0, err
			}
			return item, nonce, nil
		} else {
			continue
		}
	}
	return "", 0, hcode.ErrReqLimit
}

func (repo RepositoryStruct) UnlockNonceAddr(ctx context.Context, addr string) {
	var lockKey = lockSysAddrForTransfer(addr)
	_ = repo.Lock.UnLock(ctx, lockKey, 0)
}
