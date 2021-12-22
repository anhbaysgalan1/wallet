package repository

import (
	"context"
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"strings"
	"tp_wallet/config"
	"tp_wallet/internal/common"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/log"
	"tp_wallet/pkg/tool"
)

// BillCurrencyGetForAsync 获取链上异步未结算账单
func (repo RepositoryStruct) BillCurrencyGetForAsync(ctx context.Context) ([]*entity.Bill, error) {
	return repo.Db.BillH2OGetForQueue(ctx)
}

// BillNftGetForAsync 获取链上异步未结算账单
func (repo RepositoryStruct) BillNftGetForAsync(ctx context.Context) ([]*entity.Bill, error) {
	return repo.Db.BillNftGetForQueue(ctx)
}

// BillDealWith 处理一笔订单
func (repo RepositoryStruct) BillDealWith(ctx context.Context, bill *entity.Bill) (err error) {
	if bill == nil {
		log.GetLogger().Error("[BillDealWith] parameter failed", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	// 在这里组装公共数据数据
	if len(bill.FromAddr) > 0 {
		bill.FromUid, _ = repo.UidGetByAddress(ctx, bill.FromAddr)
	}
	if len(bill.ToAddr) > 0 {
		bill.ToUid, _ = repo.UidGetByAddress(ctx, bill.ToAddr)
	}
	switch bill.BillStatus {
	case walletPb.BillStatus_Pending:
		return repo.BillDealWithPending(ctx, bill)
	case walletPb.BillStatus_Success:
		return repo.BillDealWithSuccess(ctx, bill)
	case walletPb.BillStatus_Failed:
		return repo.BillDealWithFailed(ctx, bill)
	default:
		return hcode.ErrInternalParameter
	}
}

// BillDealWithPending 处理上链中订单
func (repo RepositoryStruct) BillDealWithPending(ctx context.Context, bill *entity.Bill) error {
	switch bill.BillType {
	case entity.BillType_Eip20:
		return repo.BillCurrencyDealWithPending(ctx, bill)
	case entity.BillType_Eip721:
		return repo.BillNftDealWithPending(ctx, bill)
	default:
		return hcode.ErrBillType
	}
}

func (repo RepositoryStruct) BillCurrencyDealWithPending(ctx context.Context, bill *entity.Bill) error {
	var err error
	if len(bill.Hash) <= 0 || bill.NumericalOrder == 0 || bill.Id.IsZero() {
		log.GetLogger().Error("[BillH2ODealWithPending] hash length or id", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	// 修改订单hash以及订单状态
	bill.DealWithBillPending()
	var filter = map[string]interface{}{"numerical_order": bill.NumericalOrder}
	err = repo.Db.BillSet(ctx, bill, filter)
	return err
}

func (repo RepositoryStruct) BillNftDealWithPending(ctx context.Context, bill *entity.Bill) error {
	var err error
	if len(bill.Hash) <= 0 || len(bill.FromAddr) == 0 || len(bill.ToAddr) == 0 {
		log.GetLogger().Error("[BillNftDealWithPending] hash length or id", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	bill.DealWithBillPending()
	var filter = map[string]interface{}{"numerical_order": bill.NumericalOrder}
	err = repo.Db.BillSet(ctx, bill, filter)
	return err
}

func (repo RepositoryStruct) BillDealWithSuccess(ctx context.Context, bill *entity.Bill) error {
	if len(bill.Hash) == 0 || len(bill.ToAddr) == 0 || len(bill.FromAddr) == 0 {
		return hcode.ErrInternalParameter
	}
	var billForStore *entity.Bill
	billForStore, _ = repo.Db.BillGetByHash(ctx, bill.Hash)
	if billForStore != nil { // 内部账单，直接赋值属性
		bill.NumericalOrder = billForStore.NumericalOrder
		bill.TransferType = billForStore.TransferType
		if bill.ContractRecord.NftToken != "" {
			billForStore.ContractRecord.NftToken = bill.ContractRecord.NftToken
		}
		bill.ContractRecord = billForStore.ContractRecord
		bill.BillType = billForStore.BillType
	} else { // 外部账单到系统内，只有充值，和授权两种可能，（目前授权暂时不考虑，待以后开发）
		switch bill.BillType {
		case entity.BillType_Eip20:
			if bill.ToUid <= common.UidStart {
				bill.TransferType = walletPb.TransferType_CurrencyCHARGE
			} else { // 可能是外部第三方转账的数据，无须在内部判断，
				log.GetLogger().Info("[BillDealWithSuccess] External third party transfer for currency", zap.Any("bill", bill))
			}
		case entity.BillType_Eip721, entity.BillType_Eip1155: // nft授权在这里添加判断
			if bill.ToUid <= common.UidStart {
				bill.TransferType = walletPb.TransferType_NftCHARGE
			}
		default:
			log.GetLogger().Error("[BillDealWithSuccess] bill type error", zap.Any("bill", bill))
			return hcode.ErrBillType
		}
		bill.BillType = 0 // 这里之后已经不需要了，减少更新数据时，修改这个字段，因为字段有索引
	}
	switch bill.TransferType { // 根据交易类型判断
	case walletPb.TransferType_CurrencyTransfer:
		log.GetLogger().Info("[BillDealWithSuccess] error for bill type", zap.Any("bill", bill))
		return hcode.ErrBillType
	case walletPb.TransferType_CurrencyCASH:
		bill.DealWithBillToCash(billForStore)
		return repo.BillDealWithSuccessForCurrencyCash(ctx, bill)
	case walletPb.TransferType_CurrencyCHARGE:
		return repo.BillDealWithSuccessForCurrencyCharge(ctx, bill)
	case walletPb.TransferType_NftTransfer:
		return repo.BillDealWithSuccessNftTransfer(ctx, bill)
	case walletPb.TransferType_NftApproval:
		// todo: 授权给我们，需要展示他授权给我们的内容
		log.GetLogger().Info("[BillDealWithSuccess] error for bill type", zap.Any("bill", bill))
		return hcode.ErrBillType
	case walletPb.TransferType_NftApprovalTransfer:
		// todo: 需要修改nft拥有者 & 删除授权给我们的nft道具信息
		log.GetLogger().Info("[BillDealWithSuccess] error for bill type", zap.Any("bill", bill))
		return hcode.ErrBillType
	case walletPb.TransferType_NftCreate:
		return repo.BillDealWithSuccessForNftCreate(ctx, bill)
	case walletPb.TransferType_NftCASH:
		return repo.BillDealWithSuccessNftCash(ctx, bill)
	case walletPb.TransferType_NftCHARGE:
		return repo.BillDealWithSuccessNftCharge(ctx, bill)
	default:
		log.GetLogger().Error("[BillDealWithSuccess] bill type check failed",
			zap.Any("inset bill", bill),
			zap.Any("store bill", billForStore))
		return hcode.ErrInternalParameter
	}
}

// BillDealWithSuccessForCurrencyCash 处理货币提现
func (repo RepositoryStruct) BillDealWithSuccessForCurrencyCash(ctx context.Context, bill *entity.Bill) error {
	var err error
	if bill.NumericalOrder == 0 {
		log.GetLogger().Error("[BillDealWithSuccessForCurrencyCash] internal parameter failed", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	// 修改订单状态
	err = repo.Db.BillSet(ctx, bill, map[string]interface{}{"numerical_order": bill.NumericalOrder})
	return err
}

// BillDealWithSuccessForCurrencyCharge 处理充值订单
func (repo RepositoryStruct) BillDealWithSuccessForCurrencyCharge(ctx context.Context, bill *entity.Bill) error {
	var err error
	if len(bill.Hash) == 0 || bill.FromUid == 0 || len(bill.FromAddr) == 0 || bill.ToUid == 0 || len(bill.ToAddr) == 0 || bill.BalanceRecord.IsEmpty() {
		log.GetLogger().Error("[BillDealWithSuccessForCurrencyCharge] hash, amount or from addr length", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	if _, ok := config.SysAccountMapByUid[bill.ToUid]; !ok {
		log.GetLogger().Error("[BillDealWithSuccessForCurrencyCharge] SysAccountMapByUid check failed", zap.Any("bill", bill))
		return hcode.ErrSysAccountNotFound
	} else {
		if bill.BalanceRecord.Currency == "" {
			bill.BalanceRecord.Currency = strings.Split(config.SysAccountMapByUid[bill.ToUid].SysToken, "_")[0] // 取下划线为分隔符第一个为币种
		} else {
			if bill.BalanceRecord.Currency != strings.Split(config.SysAccountMapByUid[bill.ToUid].SysToken, "_")[0] {
				log.GetLogger().Error("[BillDealWithSuccessForCurrencyCharge] currency check failed",
					zap.Any("bill", bill),
					zap.String("currency", strings.Split(config.SysAccountMapByUid[bill.ToUid].SysToken, "_")[0]))
			}
		}
	}
	bill.NumericalOrder = uint64(tool.GetSnowFlake().GetId())
	bill.BalanceRecord.Amount = bill.BalanceRecord.ReceiveAmount
	bill.Uid = bill.FromUid
	bill.BillType = entity.BillType_Eip20
	// 事务
	err = repo.Mongo.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			log.GetLogger().Error("[BillDealWithSuccessForCurrencyCharge] collection.InsertOne failed",
				zap.Error(err))
			return hcode.ErrInternalDb
		}
		// 修改用户金额
		bill.BalanceRecord.BeforeBalance, bill.BalanceRecord.AfterBalance, err = repo.BalanceSet(sessionContext, bill.FromUid, bill.BalanceRecord.ReceiveAmount, bill.BalanceRecord.Currency, true)
		if err != nil {
			return err
		}
		// 创建订单
		if err = repo.Db.BillCreate(sessionContext, bill); err != nil {
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				log.GetLogger().Error("[BillDealWithSuccessForCurrencyCharge] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return hcode.ErrInternalDb
			}
			log.GetLogger().Error("[BillDealWithSuccessForCurrencyCharge] Db.BillCreate failed",
				zap.Error(err))
			return err
		}
		// 给系统用户价钱
		_, _, err = repo.BalanceSet(sessionContext, bill.ToUid, bill.BalanceRecord.ReceiveAmount, bill.BalanceRecord.Currency, true)
		if err != nil {
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				log.GetLogger().Error("[BillDealWithSuccessForCurrencyCharge] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return hcode.ErrInternalDb
			}
			log.GetLogger().Error("[BillDealWithSuccessForCurrencyCharge] repo.BalanceSet failed",
				zap.Error(err))
			return err
		}
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.GetLogger().Error("[BillDealWithSuccessForCurrencyCharge] sessionContext.CommitTransaction failed",
				zap.Error(err))
			return err
		}
		return nil
	})
	return err
}

// BillDealWithSuccessForNftCreate 处理创建
func (repo RepositoryStruct) BillDealWithSuccessForNftCreate(ctx context.Context, bill *entity.Bill) error {
	var err error
	if len(bill.Hash) <= 0 || len(bill.ToAddr) == 0 || bill.ContractRecord.IsEmpty() || len(bill.ContractRecord.NftToken) == 0 || len(bill.ContractRecord.ContractType) == 0 {
		log.GetLogger().Error("[BillDealWithSuccessForNftCreate] hash length", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	// 获取nft拥有者
	var nftOwner *entity.NftOwner
	nftOwner, err = repo.Db.NftOwnerGetByCreateHash(ctx, bill.Hash)
	if err != nil {
		return err
	}
	nftOwner.NftToken = bill.ContractRecord.NftToken
	nftOwner.NftData.NftBlockToken = bill.ContractRecord.NftToken
	// 事务
	err = repo.Mongo.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			log.GetLogger().Error("[BillDealWithSuccessForNftCreate] collection.InsertOne failed",
				zap.Error(err))
			return hcode.ErrInternalDb
		}
		// 修改账单状态
		err = repo.Db.BillSet(sessionContext, &entity.Bill{
			BillStatus:     walletPb.BillStatus_Success,
			ContractRecord: bill.ContractRecord,
		}, map[string]interface{}{"hash": bill.Hash})
		if err != nil {
			return err
		}
		// 修改nft拥有者
		err = repo.Db.NftOwnerSetByGameToken(sessionContext, nftOwner)
		if err != nil {
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				log.GetLogger().Error("[BillDealWithSuccessForNftCreate] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return hcode.ErrInternalDb
			}
			log.GetLogger().Error("[BillDealWithSuccessForNftCreate] Db.NftOwnerSetByGameToken failed",
				zap.Error(err))
			return err
		}
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.GetLogger().Error("[BillDealWithSuccessForNftCreate] sessionContext.CommitTransaction failed",
				zap.Error(err))
			return err
		}
		return nil
	})
	// 推送kafka
	_ = repo.PropsSrv.PushNftCreateSuccess(ctx, walletPb.Metadata{
		NumericalOrder: bill.NumericalOrder,
		Uid:            bill.Uid,
		Cid:            bill.Cid,
		TransferType:   bill.TransferType,
		ContractRecord: &walletPb.ContractRecord{
			ContractType: bill.ContractRecord.ContractType,
			ContractAddr: bill.ContractRecord.ContractAddr,
			NftToken:     bill.ContractRecord.NftToken,
			GameToken:    bill.ContractRecord.GameToken,
			Num:          bill.ContractRecord.Num,
		},
	})
	return err
}

// BillDealWithSuccessNftCash 处理提现
func (repo RepositoryStruct) BillDealWithSuccessNftCash(ctx context.Context, bill *entity.Bill) error {
	var err error
	if len(bill.Hash) <= 0 || len(bill.ToAddr) == 0 {
		log.GetLogger().Error("[BillDealWithSuccessNftCash] hash length", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	// 事务
	err = repo.Mongo.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			log.GetLogger().Error("[BillDealWithSuccessNftCash] collection.InsertOne failed",
				zap.Error(err))
			return hcode.ErrInternalDb
		}
		// 修改账单状态
		err = repo.Db.BillSet(sessionContext, &entity.Bill{
			Hash:       bill.Hash,
			BillStatus: walletPb.BillStatus_Success,
		}, map[string]interface{}{"hash": bill.Hash})
		if err != nil {
			return err
		}
		// 修改nft拥有者
		switch bill.BillType {
		case entity.BillType_Eip721:
			if err = repo.Db.NftOwnerSetByToken(sessionContext, &entity.NftOwner{
				ContractToken: bill.ContractRecord.ContractType,
				NftToken:      bill.ContractRecord.NftToken,
				OwnerAddress:  bill.ToAddr,
				Status:        entity.NftOwnerStatus_Available,
			}); err != nil {
				if err := sessionContext.AbortTransaction(sessionContext); err != nil {
					log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.AbortTransaction failed",
						zap.Error(err))
					return hcode.ErrInternalDb
				}
				log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return err
			}
		case entity.BillType_Eip1155:
			var storeNftOwner *entity.NftOwner
			var storeUserNftOwner *entity.NftOwner
			storeNftOwner, err = repo.Db.NftOwnerGetByNftToken(ctx, bill.ContractRecord.ContractType, bill.ContractRecord.NftToken)
			if bill.BillType == entity.BillType_Eip721 {
				if storeNftOwner.OwnerAddress == bill.ToAddr {
					log.GetLogger().Warn("[BillDealWithSuccessNftCharge] bill already the owner", zap.Any("bill", bill))
					return nil
				}
			} else {
				storeUserNftOwner, _ = repo.Db.NftOwnerGetByNftTokenAndOwnerAddr(ctx, bill.ContractRecord.ContractType, bill.ContractRecord.NftToken, bill.ToAddr)
			}
			if err != nil {
				return err
			}
			// 扣除系统库存
			if err = repo.Db.NftOwnerInventory(sessionContext, storeNftOwner.ContractToken, storeNftOwner.NftToken, storeNftOwner.OwnerAddress, bill.ContractRecord.Num, false); err != nil {
				if err := sessionContext.AbortTransaction(sessionContext); err != nil {
					log.GetLogger().Error("[NftCreate] sessionContext.AbortTransaction failed",
						zap.Error(err))
					return hcode.ErrInternalDb
				}
				return err
			}
			if storeUserNftOwner != nil { // 已经存在，添加用户库存
				if err = repo.Db.NftOwnerInventory(sessionContext, storeUserNftOwner.ContractToken, storeUserNftOwner.NftToken, storeUserNftOwner.OwnerAddress, bill.ContractRecord.Num, false); err != nil {
					if err := sessionContext.AbortTransaction(sessionContext); err != nil {
						log.GetLogger().Error("[NftCreate] sessionContext.AbortTransaction failed",
							zap.Error(err))
						return hcode.ErrInternalDb
					}
					return err
				}
			} else { // 不存在，创建
				storeNftOwner.OwnerAddress = bill.ToAddr
				storeNftOwner.NftData.Num = bill.ContractRecord.Num
				storeNftOwner.Id = primitive.NewObjectID()
				if err = repo.Db.NftOwnerCreate(sessionContext, storeNftOwner); err != nil {
					log.GetLogger().Error("[NftCreate] sessionContext.AbortTransaction failed",
						zap.Error(err))
					return hcode.ErrInternalDb
				}
				return err
			}
		default:
			log.GetLogger().Error("[BillDealWithSuccessNftCash] bill type", zap.Any("bill", bill))
			return hcode.ErrBillType
		}
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.CommitTransaction failed",
				zap.Error(err))
			return err
		}
		return nil
	})
	if bill.BillType != entity.BillType_Eip20 {
		// 推送kafka
		_ = repo.PropsSrv.PushNftCashSuccess(ctx, walletPb.Metadata{
			NumericalOrder: bill.NumericalOrder,
			Uid:            bill.Uid,
			Cid:            bill.Cid,
			TransferType:   bill.TransferType,
			ContractRecord: &walletPb.ContractRecord{
				ContractType: bill.ContractRecord.ContractType,
				ContractAddr: bill.ContractRecord.ContractAddr,
				NftToken:     bill.ContractRecord.NftToken,
				GameToken:    bill.ContractRecord.GameToken,
				Num:          bill.ContractRecord.Num,
			},
			TpRecord: &walletPb.TpRecord{
				OrderId: bill.TpRecord.OrderId,
				Type:    bill.TpRecord.Type,
				Remarks: bill.TpRecord.Remarks,
				Data:    bill.TpRecord.Data,
			},
		})
	}
	return err
}

// BillDealWithSuccessNftCharge 处理充值订单
func (repo RepositoryStruct) BillDealWithSuccessNftCharge(ctx context.Context, bill *entity.Bill) error {
	if bill.ContractRecord.IsEmpty() || bill.FromUid == 0 || len(bill.FromAddr) == 0 || bill.ToUid == 0 || len(bill.ToAddr) == 0 || len(bill.Hash) == 0 {
		log.GetLogger().Error("[BillDealWithSuccessNftCharge] parameter check failed", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	var err error
	// 判断nft是否重复
	var sysUid uint64
	sysUid, err = repo.UidGetByAddress(ctx, bill.ContractRecord.ContractAddr)
	if err != nil {
		return err
	}
	var storeNftOwner *entity.NftOwner
	storeNftOwner, err = repo.Db.NftOwnerGetByNftToken(ctx, config.SysAccountMapByUid[sysUid].SysToken, bill.ContractRecord.NftToken)
	if bill.BillType == entity.BillType_Eip721 {
		if storeNftOwner.OwnerAddress == bill.ToAddr {
			log.GetLogger().Warn("[BillDealWithSuccessNftCharge] bill already the owner", zap.Any("bill", bill))
			return nil
		}
	}
	if err != nil {
		return err
	}
	bill.BillType = entity.BillType_Eip721
	bill.Id = primitive.NewObjectID()
	bill.NumericalOrder = uint64(tool.GetSnowFlake().GetId())
	bill.Uid = bill.FromUid
	bill.ContractRecord = entity.ContractRecord{
		ContractType: storeNftOwner.ContractToken,
		ContractAddr: config.SysAccountMapByUid[sysUid].ContractAddress,
		GameId:       storeNftOwner.GameId,
		NftToken:     storeNftOwner.NftToken,
		GameToken:    storeNftOwner.GameToken,
		Num:          storeNftOwner.NftData.Num,
	}
	// 事务
	err = repo.Mongo.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			log.GetLogger().Error("[BillDealWithSuccessNftCash] collection.InsertOne failed",
				zap.Error(err))
			return hcode.ErrInternalDb
		}
		// 创建订单
		if err = repo.Db.BillCreate(sessionContext, bill); err != nil {
			return err
		}
		// 修改nft拥有者, 这里有两个类型的nft
		switch bill.BillType {
		case entity.BillType_Eip721:
			if err = repo.Db.NftOwnerSetByToken(sessionContext, &entity.NftOwner{
				NftToken:      bill.ContractRecord.NftToken,
				ContractToken: config.SysAccountMapByUid[sysUid].SysToken,
				OwnerAddress:  bill.ToAddr,
				Status:        entity.NftOwnerStatus_Available,
			}); err != nil {
				if err := sessionContext.AbortTransaction(sessionContext); err != nil {
					log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.AbortTransaction failed",
						zap.Error(err))
					return hcode.ErrInternalDb
				}
				log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return err
			}
		case entity.BillType_Eip1155:
			if err = repo.Db.NftOwnerInventory(sessionContext, storeNftOwner.ContractToken, storeNftOwner.NftToken, config.SysAccountMap[bill.ContractRecord.ContractType].AddrIncome, bill.ContractRecord.Num, true); err != nil {
				if err := sessionContext.AbortTransaction(sessionContext); err != nil {
					log.GetLogger().Error("[NftCreate] sessionContext.AbortTransaction failed",
						zap.Error(err))
					return hcode.ErrInternalDb
				}
				return err
			}
		default:
			return hcode.ErrBillType
		}
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.CommitTransaction failed",
				zap.Error(err))
			return err
		}
		return nil
	})
	if bill.BillType != entity.BillType_Eip20 {
		// 推送kafka
		_ = repo.PropsSrv.PushNftCharge(ctx, walletPb.Metadata{
			NumericalOrder: bill.NumericalOrder,
			Uid:            bill.Uid,
			Cid:            bill.Cid,
			TransferType:   bill.TransferType,
			ContractRecord: &walletPb.ContractRecord{
				ContractType: bill.ContractRecord.ContractType,
				ContractAddr: bill.ContractRecord.ContractAddr,
				NftToken:     bill.ContractRecord.NftToken,
				GameToken:    bill.ContractRecord.GameToken,
				Num:          bill.ContractRecord.Num,
			},
		})
	}
	return err
}

// BillDealWithSuccessNftTransfer 处理nft外部交易
func (repo RepositoryStruct) BillDealWithSuccessNftTransfer(ctx context.Context, bill *entity.Bill) error {
	if bill.ContractRecord.IsEmpty() || bill.FromUid == 0 || len(bill.FromAddr) == 0 || len(bill.ToAddr) == 0 || len(bill.Hash) == 0 {
		log.GetLogger().Error("[BillDealWithSuccessNftCharge] parameter check failed", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	var err error
	// 判断nft是否重复
	var sysUid uint64
	sysUid, err = repo.UidGetByAddress(ctx, bill.ContractRecord.ContractAddr)
	if err != nil {
		return err
	}
	var storeNftOwner *entity.NftOwner
	var storeUserNftOwner *entity.NftOwner
	storeNftOwner, err = repo.Db.NftOwnerGetByNftToken(ctx, config.SysAccountMapByUid[sysUid].SysToken, bill.ContractRecord.NftToken)
	if bill.BillType == entity.BillType_Eip721 {
		if storeNftOwner.OwnerAddress == bill.ToAddr {
			log.GetLogger().Warn("[BillDealWithSuccessNftCharge] bill already the owner", zap.Any("bill", bill))
			return nil
		}
	} else {
		storeUserNftOwner, _ = repo.Db.NftOwnerGetByNftTokenAndOwnerAddr(ctx, config.SysAccountMapByUid[sysUid].SysToken, bill.ContractRecord.NftToken, bill.FromAddr)
	}
	if err != nil {
		return err
	}
	bill.BillType = entity.BillType_Eip721
	bill.Id = primitive.NewObjectID()
	bill.NumericalOrder = uint64(tool.GetSnowFlake().GetId())
	bill.Uid = bill.FromUid
	bill.ContractRecord = entity.ContractRecord{
		ContractType: storeNftOwner.ContractToken,
		ContractAddr: config.SysAccountMapByUid[sysUid].ContractAddress,
		GameId:       storeNftOwner.GameId,
		NftToken:     storeNftOwner.NftToken,
		GameToken:    storeNftOwner.GameToken,
		Num:          storeNftOwner.NftData.Num,
	}
	// 修改nft拥有者, 这里有两个类型的nft
	switch bill.BillType {
	case entity.BillType_Eip721:
		if err = repo.Db.NftOwnerSetByToken(ctx, &entity.NftOwner{
			NftToken:      bill.ContractRecord.NftToken,
			ContractToken: config.SysAccountMapByUid[sysUid].SysToken,
			OwnerAddress:  bill.ToAddr,
			Status:        entity.NftOwnerStatus_Available,
		}); err != nil {
			log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.AbortTransaction failed",
				zap.Error(err))
			return err
		}
	case entity.BillType_Eip1155:
		if err = repo.Db.NftOwnerInventory(ctx, storeUserNftOwner.ContractToken, storeUserNftOwner.NftToken, bill.FromAddr, bill.ContractRecord.Num, false); err != nil {
			return err
		}
	default:
		return hcode.ErrBillType
	}

	return err
}

// BillDealWithFailed 处理上链失败订单
func (repo RepositoryStruct) BillDealWithFailed(ctx context.Context, bill *entity.Bill) error {
	var err error
	// 获取账单
	bill, err = repo.Db.BillGetById(ctx, bill.Id)
	if err != nil {
		return err
	}
	// 事务
	err = repo.Mongo.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			log.GetLogger().Error("[BillDealWithSuccessNftCash] collection.InsertOne failed",
				zap.Error(err))
			return hcode.ErrInternalDb
		}
		var beforeBalance, afterBalance entity.Amount
		switch {
		case bill.BillType == entity.BillType_Eip20 && bill.TransferType == walletPb.TransferType_CurrencyCASH:
			// 提现，from是系统，to是用户，给to退钱
			beforeBalance, afterBalance, err = repo.BalanceSet(sessionContext, bill.ToUid, bill.BalanceRecord.Amount, bill.BalanceRecord.Currency, true)
			if err != nil {
				return err
			} else {
				log.GetLogger().Info("[BillDealWithFailed] success", zap.Uint64("uid", bill.ToUid),
					zap.String("before balance", string(beforeBalance)),
					zap.String("after balance", string(afterBalance)))
				return nil
			}
		case bill.BillType == entity.BillType_Eip721 && bill.TransferType == walletPb.TransferType_NftCASH:
			if err = repo.Db.NftOwnerSetByToken(sessionContext, &entity.NftOwner{
				NftToken:      bill.ContractRecord.NftToken,
				ContractToken: bill.ContractRecord.ContractType,
				OwnerAddress:  bill.ToAddr,
				Status:        entity.NftOwnerStatus_Available,
			}); err != nil {
				return err
			} else {
				log.GetLogger().Info("[BillDealWithFailed] success", zap.Uint64("uid", bill.ToUid),
					zap.String("before balance", string(beforeBalance)),
					zap.String("after balance", string(afterBalance)))
			}
		case bill.BillType == entity.BillType_Eip1155 && bill.TransferType == walletPb.TransferType_NftCASH:
			// 给用户返回nft
			if err = repo.Db.NftOwnerInventory(sessionContext, bill.ContractRecord.ContractType, bill.ContractRecord.NftToken, bill.FromAddr, bill.ContractRecord.Num, true); err != nil {
				if err := sessionContext.AbortTransaction(sessionContext); err != nil {
					log.GetLogger().Error("[NftCreate] sessionContext.AbortTransaction failed",
						zap.Error(err))
					return hcode.ErrInternalDb
				}
				return err
			}
		default:
			log.GetLogger().Error("[BillDealWithFailed] bill deal with failed for bill type", zap.Any("bill", bill))
			return hcode.ErrBillType
		}
		var remarks = entity.ErrorBillRemark{ // 存储错误信息
			To:              bill.ToUid,
			ToBeforeBalance: string(beforeBalance),
			ToAfterBalance:  string(afterBalance),
			Data:            bill.Remark,
		}.ToJson()
		// 修改订单状态
		err = repo.Db.BillSet(sessionContext, &entity.Bill{Id: bill.Id, BillStatus: walletPb.BillStatus_Failed, Remark: remarks}, map[string]interface{}{"_id": bill.Id})
		if err != nil {
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return hcode.ErrInternalDb
			}
			log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.AbortTransaction failed",
				zap.Error(err))
			return err
		}
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.GetLogger().Error("[BillDealWithSuccessNftCash] sessionContext.CommitTransaction failed",
				zap.Error(err))
			return err
		}
		return nil
	})
	return err
}
