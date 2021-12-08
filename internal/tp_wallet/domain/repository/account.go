package repository

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"tp_wallet/config"
	"tp_wallet/internal/common"
	"tp_wallet/internal/wallet/dto"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/hcode"
	"tp_wallet/pkg/log"
)

func (repo RepositoryStruct) AccountGetAndCreate(ctx context.Context, req *dto.UidReq) (*dto.AccountGetResp, error) {
	// 地址是否已经注册
	if req.GetUid() <= common.UidStart { // 预留10000个系统账号
		return nil, hcode.ErrUidIsTooSort
	}
	var uid uint64
	var err error
	uid, err = repo.UidGetByAddress(ctx, req.GetAddr())
	if err != nil && !errors.Is(err, hcode.ErrUserNotFound) {
		return nil, hcode.ErrServer
	}
	if uid == 0 { // 注册用户
		err = repo.AccountRegister(ctx, req.GetUid(), req.GetCurrency(), req.GetAddr())
		if err != nil {
			return nil, err
		} else {
			return &dto.AccountGetResp{
				Uid: req.GetUid(),
				//Address: []string{req.GetAddr()},
				Balance: nil,
				IsExist: false,
			}, nil
		}
	} else { // 查询用户信息
		var balance map[string]entity.Amount
		balance, err = repo.BalanceGetByUid(ctx, req.GetUid())
		if err != nil {
			return nil, err
		}
		var result = &dto.AccountGetResp{
			Uid:     req.GetUid(),
			Balance: entity.AmountMapToString(balance),
			IsExist: true,
		}
		return result, nil
	}
}

func (repo RepositoryStruct) AddressGetByUid(ctx context.Context, uid uint64) (map[string]struct{}, error) {
	if uid == config.BlockBusiness.H2OSysUid {
		var result = make(map[string]struct{})
		for _, item := range config.BlockBusiness.H2OUidAddrExpenditure {
			result[item] = struct{}{}
		}
		result[config.BlockBusiness.H2OUidAddrIncome] = struct{}{}
		return result, nil
	}
	if uid == config.BlockBusiness.NftSysUid {
		var result = make(map[string]struct{})
		for _, item := range config.BlockBusiness.NftAddrExpenditure {
			result[item] = struct{}{}
		}
		result[config.BlockBusiness.NftAddrIncome] = struct{}{}
		return result, nil
	}
	for _, items := range config.BlockBusiness.AddressForGame {
		var result = make(map[string]struct{})
		if uid == items.SysUid {
			for _, item := range items.AddrExpenditure {
				result[item] = struct{}{}
			}
			result[items.AddrIncome] = struct{}{}
		}
	}
	var addr map[string]struct{}
	var err error
	addr, err = repo.Cache.AddressGetByUid(ctx, uid)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return addr, err
		} else {
			var addrS []*entity.Address
			addrS, err = repo.Db.AddressGetByUid(ctx, uid)
			if err != nil || len(addrS) <= 0 {
				return nil, hcode.ErrAddressGet
			}
			addr = make(map[string]struct{})
			for _, item := range addrS {
				addr[item.Address] = struct{}{}
			}
			_ = repo.Cache.AddressSaveByUid(ctx, uid, addr)
		}
	}
	return addr, nil
}

func (repo RepositoryStruct) BalanceGetByUid(ctx context.Context, uid uint64) (map[string]entity.Amount, error) {
	var balance map[string]entity.Amount
	var err error
	balance, err = repo.Cache.BalanceGetByUid(ctx, uid)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return balance, err
		} else {
			var balanceResp []*entity.Balance
			balanceResp, err = repo.Db.BalanceGetByUid(ctx, uid)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return balance, hcode.ErrUserNotFound
				}
				return balance, hcode.ErrBalanceGet
			}
			balance = make(map[string]entity.Amount)
			for _, item := range balanceResp {
				balance[item.Currency] = item.Balance
			}
			_ = repo.Cache.BalanceSave(ctx, balanceResp)
		}
	}
	return balance, nil
}

func (repo RepositoryStruct) UidGetByAddress(ctx context.Context, addr string) (uint64, error) {
	if _, ok := config.BlockSysAddrToUid[addr]; ok {
		return config.BlockSysAddrToUid[addr], nil
	}
	var uid uint64
	var err error
	uid, err = repo.Cache.UidGetByAddr(ctx, addr)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return 0, hcode.ErrServer
		} else {
			var address *entity.Address
			address, err = repo.Db.AddressGetByAddr(ctx, addr)
			if err != nil || address.Uid <= 0 {
				return 0, hcode.ErrUserNotFound
			}
			uid = address.Uid
			_ = repo.Cache.AddressSave(ctx, address)
		}
	}
	return uid, nil
}

func (repo RepositoryStruct) AccountRegister(ctx context.Context, uid uint64, currency, address string) error {
	var err error
	var AddressCreateReq = &entity.Address{
		Id:          primitive.NewObjectID(),
		Uid:         uid,
		Address:     address,
		Currency:    currency,
		AccountType: entity.AccountType_Ordinary,
	}
	if !AddressCreateReq.CheckAddressCreate() {
		log.GetLogger().Error("[AccountRegister] AddressCreateReq.CheckAddressCreate failed",
			zap.Any("req", AddressCreateReq))
		return hcode.ErrInternalParameter
	}
	err = repo.mongo.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			log.GetLogger().Error("[AccountRegister] collection.InsertOne failed",
				zap.Any("uid", uid),
				zap.Any("currency", currency),
				zap.Any("address", address),
				zap.Error(err))
			return hcode.ErrInternalDb
		}
		// 创建地址
		err = repo.Db.AddressCreate(ctx, AddressCreateReq)
		if err != nil {
			err = sessionContext.AbortTransaction(sessionContext)
			log.GetLogger().Error("[BalanceTransfer] Db.AddressCreate failed",
				zap.Any("AddressCreateReq", AddressCreateReq),
				zap.Error(err))
			return hcode.ErrServer
		}
		// 创建余额
		err = repo.Db.BalanceCreate(ctx, uid, config.CurrencyList[0].Name)
		if err != nil {
			err = sessionContext.AbortTransaction(sessionContext)
			log.GetLogger().Error("[BalanceTransfer] repo.Db.BalanceCreate failed",
				zap.Any("uid", uid),
				zap.Any("currency", config.CurrencyList[0].Name),
				zap.Error(err))
			return err
		}
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.GetLogger().Error("[BalanceTransfer] sessionContext.CommitTransaction failed",
				zap.Error(err))
			return err
		}
		return nil
	})
	return err
}
