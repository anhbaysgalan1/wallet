package repository

import (
	"context"
	"errors"
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"tp_wallet/config"
	common2 "tp_wallet/internal/common"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/log"
	"tp_wallet/pkg/tool"
)

// TransferCurrencyForOffline 离线转账(同步)
func (repo RepositoryStruct) TransferCurrencyForOffline(ctx context.Context, req *walletPb.TransferForOfflineReq) (*walletPb.Empty, error) {
	// 生成订单
	var err error
	// 事务
	if err = repo.Mongo.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			log.GetLogger().Error("[TransferCurrencyForOffline] collection.InsertOne failed",
				zap.Error(err))
			return hcode.ErrInternalDb
		}
		var orderId = tool.GetSnowFlake().GetId()
		var tpRecord entity.TpRecord
		if req.GetTpRecord() != nil {
			tpRecord = entity.TpRecord{
				OrderId: req.GetTpRecord().GetOrderId(),
				Type:    req.GetTpRecord().GetType(),
				Remarks: req.GetTpRecord().GetRemarks(),
				Data:    req.GetTpRecord().GetData(),
			}
		}
		var fromBill = &entity.Bill{
			Id:             primitive.NewObjectID(),
			NumericalOrder: uint64(orderId),
			Uid:            req.From,
			Cid:            req.Cid,
			FromUid:        req.GetFrom(),
			ToUid:          req.GetTo(),
			TransferType:   req.GetTransferType(),
			BillType:       entity.BillType_Eip20,
			BillStatus:     walletPb.BillStatus_Success, // 本地只有一次性成功
			BalanceRecord: entity.BalanceRecord{
				Currency:      req.Currency,
				Amount:        req.Amount,
				ReceiveAmount: req.Amount,
				BeforeBalance: "",
				AfterBalance:  "",
			},
			TpRecord: tpRecord,
		}
		var toBill = &entity.Bill{
			Id:             primitive.NewObjectID(),
			NumericalOrder: uint64(orderId),
			Uid:            req.To,
			Cid:            req.Cid,
			FromUid:        req.GetFrom(),
			ToUid:          req.GetTo(),
			TransferType:   req.GetTransferType(),
			BillType:       entity.BillType_Eip20,
			BillStatus:     walletPb.BillStatus_Success, // 本地只有一次性成功
			BalanceRecord: entity.BalanceRecord{
				Currency:      req.Currency,
				Amount:        req.Amount,
				ReceiveAmount: req.Amount,
				BeforeBalance: "",
				AfterBalance:  "",
			},
			TpRecord: tpRecord,
		}
		// 修改from余额
		fromBill.BalanceRecord.BeforeBalance, fromBill.BalanceRecord.AfterBalance, err = repo.BalanceSet(sessionContext, fromBill.Uid, fromBill.BalanceRecord.Amount, fromBill.BalanceRecord.Currency, false)
		if err != nil {
			return err
		}
		// 创建from订单
		if err = repo.Db.BillCreate(sessionContext, fromBill); err != nil {
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				log.GetLogger().Error("[TransferCurrencyForOffline] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return hcode.ErrInternalDb
			}
			return err
		}
		// 修改to余额
		toBill.BalanceRecord.BeforeBalance, toBill.BalanceRecord.AfterBalance, err = repo.BalanceSet(sessionContext, toBill.Uid, toBill.BalanceRecord.Amount, toBill.BalanceRecord.Currency, true)
		if err != nil {
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				log.GetLogger().Error("[TransferCurrencyForOffline] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return hcode.ErrInternalDb
			}
			return err
		}
		// 创建to订单
		if err = repo.Db.BillCreate(sessionContext, toBill); err != nil {
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				log.GetLogger().Error("[TransferCurrencyForOffline] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return hcode.ErrInternalDb
			}
			return err
		}
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.CommitTransaction failed",
				zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.CommitTransaction failed",
			zap.Error(err))
		return &walletPb.Empty{}, err
	}
	return &walletPb.Empty{}, err
}

// TransferCurrencyCash 提现 (异步)
func (repo RepositoryStruct) TransferCurrencyCash(ctx context.Context, req *entity.Bill) (*walletPb.Empty, error) {
	var err error
	// 加锁
	var lockResult bool
	var lockKey = common2.KeyLockCurrencyCash(req.ToUid)
	lockResult, _, err = repo.Lock.TryLock(lockKey, 0, common2.LockCurrencyTtl, ctx)
	if err != nil {
		return nil, hcode.ErrInternalCache
	}
	if !lockResult {
		return nil, hcode.ErrReqLimit
	}
	defer repo.Lock.UnLock(ctx, lockKey, 0)
	// 生成订单
	req.Id = primitive.NewObjectID()
	req.NumericalOrder = uint64(tool.GetSnowFlake().GetId())
	req.TransferType = walletPb.TransferType_CurrencyCASH
	req.BillType = entity.BillType_Eip20
	req.BillStatus = walletPb.BillStatus_Queuing
	// 事务
	err = repo.Mongo.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			log.GetLogger().Error("[TransferCurrencyCash] collection.InsertOne failed",
				zap.Error(err))
			return hcode.ErrInternalDb
		}
		// 修改From余额
		req.BalanceRecord.BeforeBalance, req.BalanceRecord.AfterBalance, err = repo.BalanceSet(ctx, req.ToUid, req.BalanceRecord.Amount, req.BalanceRecord.Currency, false)
		if err != nil {
			return err
		}
		// 生成订单
		err = repo.Db.BillCreate(ctx, req)
		if err != nil {
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				log.GetLogger().Error("[TransferCurrencyCash] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return hcode.ErrInternalDb
			}
			return hcode.ErrServer
		}
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.GetLogger().Error("[TransferCurrencyCash] sessionContext.CommitTransaction failed",
				zap.Error(err))
			return err
		}
		return nil
	})
	return &walletPb.Empty{}, err
}

func (repo RepositoryStruct) TransferNftCash(ctx context.Context, req *entity.Bill) (*walletPb.Empty, error) {
	var err error
	// 加锁
	var lockResult bool
	var lockKey string = common2.KeyLockNftCash(req.ContractRecord.NftToken)
	lockResult, _, err = repo.Lock.TryLock(lockKey, 0, common2.LockNftTtl, ctx)
	if err != nil {
		return nil, hcode.ErrInternalCache
	}
	if !lockResult {
		return nil, hcode.ErrReqLimit
	}
	defer repo.Lock.UnLock(ctx, lockKey, 0)
	// 判断合约条件
	if _, ok := config.SysAccountMap[req.ContractRecord.ContractType]; !ok {
		return nil, hcode.ErrContractNftAttributionForAddr
	}
	req.FromUid = config.SysAccountMap[req.ContractRecord.ContractType].SysUid
	req.FromAddr = config.SysAccountMap[req.ContractRecord.ContractType].AddrExpenditure[0]
	// 生成订单
	req.Id = primitive.NewObjectID()
	req.NumericalOrder = uint64(tool.GetSnowFlake().GetId())
	req.TransferType = walletPb.TransferType_NftCASH
	req.BillType, err = config.GetBillTypeBySysUid(config.SysAccountMap[req.ContractRecord.ContractType].SysUid)
	if err != nil {
		return nil, hcode.ErrBillType
	}
	// 如果是1155代币，需要提现数量大于1
	if req.BillType == entity.BillType_Eip1155 {
		if req.ContractRecord.Num < 1 {
			log.GetLogger().Error("[TransferNftCash] cash num error", zap.Any("bill", req))
			return nil, hcode.ErrParameter
		}
	}
	req.BillStatus = walletPb.BillStatus_Queuing
	// 事务
	err = repo.Mongo.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			log.GetLogger().Error("[TransferNftCash] collection.InsertOne failed",
				zap.Error(err))
			return hcode.ErrInternalDb
		}
		// 创建订单
		err = repo.Db.BillCreate(sessionContext, req)
		if err != nil {
			return err
		}
		if err = repo.Db.NftOwnerSetByToken(ctx, &entity.NftOwner{
			NftToken:      req.ContractRecord.NftToken,
			ContractToken: req.ContractRecord.ContractType,
			Status:        entity.NftOwnerStatus_Cash,
			OwnerAddress:  req.ToAddr,
		}); err != nil {
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				log.GetLogger().Error("[TransferNftCash] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return hcode.ErrInternalDb
			}
			return err
		}
		// 扣减nft材料库存
		if req.BillType == entity.BillType_Eip1155 {
			if err = repo.Db.NftOwnerInventory(ctx, req.ContractRecord.ContractType, req.ContractRecord.NftToken, req.ToAddr, req.ContractRecord.Num, false); err != nil {
				if err := sessionContext.AbortTransaction(sessionContext); err != nil {
					log.GetLogger().Error("[TransferNftCash] sessionContext.AbortTransaction failed",
						zap.Error(err))
					return hcode.ErrInternalDb
				}
				return err
			}
		}
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.GetLogger().Error("[TransferNftCash] sessionContext.CommitTransaction failed",
				zap.Error(err))
			return err
		}
		return nil
	})
	return &walletPb.Empty{}, err
}

// BalanceSet 修改用户余额，isAdd是用来判断加减钱，true: 加钱，false:减钱
func (repo RepositoryStruct) BalanceSet(ctx context.Context, uid uint64, amount, currency string, isAdd bool) (entity.Amount, entity.Amount, error) {
	var err error
	if _, ok := config.SysAccountMap[currency]; !ok {
		log.GetLogger().Error("[BalanceSet] currency unsupported",
			zap.Uint64("uid", uid),
			zap.String("amount", amount),
			zap.String("string", currency),
			zap.Bool("idAdd", isAdd))
		return "", "", hcode.ErrCurrencyUnsupported
	}
	// 系统用户单独调用修改余额方法
	if uid < common2.UidStart {
		return repo.Cache.BalanceSysSet(ctx, uid, amount, currency, isAdd)
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
	var balanceMap map[string]entity.Amount
	var balance entity.Amount
	balanceMap, err = repo.BalanceGetByUid(ctx, uid)
	if err != nil && !errors.Is(err, hcode.ErrBalanceGet) {
		return "", "", err
	}
	if _, ok := balanceMap[currency]; !ok { // 玩家没有该币种余额，数据库添加
		balance = "0"
		balanceMap[currency] = "0"
		if err = repo.Db.BalanceCreate(ctx, uid, currency); err != nil {
			return "", "", err
		}
	} else {
		balance = balanceMap[currency]
	}

	if isAdd {
		balance, err = balance.Add(amount)
		if err != nil {
			return "", "", err
		}
	} else {
		var calculateResult bool
		balance, calculateResult, err = balance.Sub(amount)
		if err != nil {
			return "", "", err
		}
		if !calculateResult {
			log.GetLogger().Error("[BalanceSet] balance insufficient")
			return "", "", hcode.ErrBalanceInsufficient
		}
	}
	var updates = &entity.Balance{
		Uid:      uid,
		Currency: currency,
		Balance:  balance,
	}
	err = repo.Db.BalanceSet(ctx, updates)
	if err != nil {
		log.GetLogger().Error("[BalanceSet] Db.BalanceSet failed",
			zap.Uint64("uid", uid),
			zap.String("amount", amount),
			zap.Any("beforeBalance", balanceMap[currency]),
			zap.Any("afterBalance", balance),
			zap.Any("newBalance", updates),
			zap.Error(err))
		return "", "", hcode.ErrServer
	}
	// 缓存双删
	_ = repo.Cache.BalanceDel(ctx, uid)
	// 打印成功log
	log.GetLogger().Info("[BalanceSet] Db.BalanceSet failed",
		zap.Uint64("uid", uid),
		zap.String("amount", amount),
		zap.Any("beforeBalance", balanceMap[currency]),
		zap.Any("afterBalance", balance),
		zap.Any("newBalance", updates),
		zap.Error(err))
	return balanceMap[currency], balance, nil
}
