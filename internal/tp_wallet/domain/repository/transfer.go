package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"tp_wallet/config"
	common2 "tp_wallet/internal/common"
	"tp_wallet/internal/wallet/dto"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/hcode"
	"tp_wallet/pkg/log"
	"tp_wallet/pkg/tool"
)

// TransferH2OForOffline 离线转账(同步)
func (repo RepositoryStruct) TransferH2OForOffline(ctx context.Context, req *dto.TransferForOfflineReq) (*dto.Empty, error) {
	// 生成订单
	var err error
	var now = tool.GetTimeUnixMilli()
	var bill = &entity.Bill{
		Id:             primitive.NewObjectID(),
		Amount:         req.GetAmount(),
		BillType:       req.TransferType,
		BillStatus:     dto.BillStatus_Success, // 本地只有一次性成功
		From:           req.From,
		To:             req.To,
		CreateTime:     now,
		UpdateTime:     now,
		IsBalanceTrade: true,
	}
	// 修改from余额
	bill.FromBeforeBalance, bill.FromAfterBalance, err = repo.BalanceSet(ctx, bill.From, bill.Amount, false)
	if err != nil {
		return nil, err
	}
	// 修改to余额
	bill.ToBeforeBalance, bill.ToAfterBalance, err = repo.BalanceSet(ctx, bill.To, bill.Amount, true)
	if err != nil {
		return nil, err
	}
	// 生成订单
	_ = repo.Db.BillCreate(ctx, bill)
	return &dto.Empty{}, nil
}

// TransferH2OCash 提现 (异步)
func (repo RepositoryStruct) TransferH2OCash(ctx context.Context, req *entity.Bill) (*dto.Empty, error) {
	var err error
	// 加锁
	var lockResult bool
	var lockKey string = common2.KeyLockH2OCash(req.To)
	lockResult, _, err = repo.Lock.TryLock(lockKey, 0, common2.LockH2OTtl, ctx)
	if err != nil {
		return nil, hcode.ErrInternalCache
	}
	if !lockResult {
		return nil, hcode.ErrReqLimit
	}
	defer repo.Lock.UnLock(ctx, lockKey, 0)
	// 生成订单
	var now = tool.GetTimeUnixMilli()
	req.Id = primitive.NewObjectID()
	req.BillType = dto.TransferType_H2OCASH
	req.BillStatus = dto.BillStatus_Queuing
	req.From = config.BlockBusiness.H2OSysUid
	req.CreateTime = now
	req.UpdateTime = now
	req.IsBalanceTrade = true
	// 修改From余额
	req.ToBeforeBalance, req.ToAfterBalance, err = repo.BalanceSet(ctx, req.To, req.Amount, false)
	if err != nil {
		return nil, err
	}
	// 生成订单
	err = repo.Db.BillCreate(ctx, req)
	if err != nil {
		return nil, hcode.ErrServer
	}
	return &dto.Empty{}, nil
}

func (repo RepositoryStruct) TransferNftCash(ctx context.Context, req *entity.Bill) (*dto.Empty, error) {
	var err error
	// 加锁
	var lockResult bool
	var lockKey string = common2.KeyLockNftCash(req.NftToken)
	lockResult, _, err = repo.Lock.TryLock(lockKey, 0, common2.LockNftTtl, ctx)
	if err != nil {
		return nil, hcode.ErrInternalCache
	}
	if !lockResult {
		return nil, hcode.ErrReqLimit
	}
	defer repo.Lock.UnLock(ctx, lockKey, 0)
	// 判断合约条件
	if _, ok := config.SysAccountMap[req.ContractType]; !ok {
		return nil, hcode.ErrContractNftAttributionForAddr
	}
	req.From = config.SysAccountMap[req.ContractType].SysUid
	var owner *entity.NftContract
	// 判断nft归属
	owner, err = repo.NftContractGetByNftToken(ctx, req.ContractType, req.NftToken)
	if err != nil {
		return nil, err
	}
	if owner.OwnerAddress != req.ToAddr || owner.DeleteTime > 0 {
		return nil, hcode.ErrContractNftAttributionForAddr
	}
	// 生成订单
	var now = tool.GetTimeUnixMilli()
	req.Id = primitive.NewObjectID()
	req.BillType = dto.TransferType_NftCASH
	req.BillStatus = dto.BillStatus_Queuing
	req.CreateTime = now
	req.UpdateTime = now
	// 软删除nft合约归属
	err = repo.NftContractSetByToken(ctx, &entity.NftContract{NftToken: owner.NftToken, DeleteTime: now})
	if err != nil {
		return nil, err
	}
	// 生成订单
	err = repo.Db.BillCreate(ctx, req)
	if err != nil {
		return nil, hcode.ErrServer
	}
	return &dto.Empty{}, nil
}

// BalanceSet 修改用户余额，isAdd是用来判断加减钱，true: 加钱，false:减钱
func (repo RepositoryStruct) BalanceSet(ctx context.Context, uid uint64, amount string, isAdd bool) (string, string, error) {
	var err error
	// 系统用户单独调用修改余额方法
	if uid < common2.UidStart {
		return repo.Cache.BalanceSysSet(ctx, amount, isAdd)
	}
	// 用户加锁
	var lockResult bool
	var lockKey = common2.KeyLockAccountBalance(uid)
	lockResult, _, err = repo.Lock.TryLock(lockKey, 0, common2.LockAccountBalanceTtl, ctx)
	if err != nil {
		log.GetLogger().Error("[BalanceSet] Lock.TryLock failed",
			zap.Uint64("uid", uid),
			zap.String("amount", amount),
			zap.Error(err))
		return "", "", hcode.ErrServer
	}
	if !lockResult {
		return "", "", hcode.ErrReqLimit
	}
	defer repo.Lock.UnLock(ctx, lockKey, 0)
	// 删除缓存
	err = repo.Cache.BalanceDel(ctx, uid)
	if err != nil {
		return "", "", hcode.ErrServer
	}
	// 计算
	var beforeBalance, afterBalance string
	var calculateResult bool
	beforeBalance, err = repo.BalanceGetByUid(ctx, uid)
	if err != nil {
		return "", "", err
	}
	afterBalance, calculateResult, err = tool.BigCalculate(beforeBalance, amount, isAdd)
	if err != nil {
		return "", "", hcode.ErrInternalParameter
	}
	if !calculateResult { // 账户余额不足
		return "", "", hcode.ErrBalanceInsufficient
	}
	// 修改余额
	var newBalance = &entity.Balance{
		Uid:     uid,
		Balance: afterBalance,
	}
	err = repo.Db.BalanceSet(ctx, newBalance)
	if err != nil {
		log.GetLogger().Error("[BalanceSet] Db.BalanceSet failed",
			zap.Uint64("uid", uid),
			zap.String("amount", amount),
			zap.String("beforeBalance", beforeBalance),
			zap.String("afterBalance", afterBalance),
			zap.Any("newBalance", newBalance),
			zap.Error(err))
		return "", "", hcode.ErrServer
	}
	// 缓存双删
	_ = repo.Cache.BalanceDel(ctx, uid)
	// 打印成功log
	log.GetLogger().Info("[BalanceSet] success",
		zap.Uint64("uid", uid),
		zap.Bool("is add", isAdd),
		zap.String("before balance", beforeBalance),
		zap.String("after balance", afterBalance))
	return beforeBalance, afterBalance, nil
}
