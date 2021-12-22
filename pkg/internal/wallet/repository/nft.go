package repository

import (
	"context"
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"tp_wallet/config"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/log"
)

func (repo RepositoryStruct) NftOwnerCreate(ctx context.Context, nft *entity.NftOwner) error {
	return repo.Db.NftOwnerCreate(ctx, nft)
}
func (repo RepositoryStruct) NftOwnerSetByToken(ctx context.Context, nft *entity.NftOwner) error {
	return repo.Db.NftOwnerSetByToken(ctx, nft)
}
func (repo RepositoryStruct) NftOwnerGetByAddr(ctx context.Context, address string, page *walletPb.Page) ([]*entity.NftOwner, error) {
	return repo.Db.NftOwnerGetByAddr(ctx, address, page)
}
func (repo RepositoryStruct) NftOwnerGetByNftToken(ctx context.Context, contractToken, token string) (*entity.NftOwner, error) {
	return repo.Db.NftOwnerGetByNftToken(ctx, contractToken, token)
}
func (repo RepositoryStruct) NftOwnerGetByGameToken(ctx context.Context, contractToken, token string) (*entity.NftOwner, error) {
	return repo.Db.NftOwnerGetByGameToken(ctx, contractToken, token)
}

func (repo RepositoryStruct) NftContractCreate(ctx context.Context, nft *entity.NftContract) error {
	return repo.Db.NftContractCreate(ctx, nft)
}
func (repo RepositoryStruct) NftContractSetByToken(ctx context.Context, nft *entity.NftContract) error {
	return repo.Db.NftContractSetByToken(ctx, nft)
}
func (repo RepositoryStruct) NftContractGetByAddr(ctx context.Context, address string, page *walletPb.Page) ([]*entity.NftContract, error) {
	return repo.Db.NftContractGetByAddr(ctx, address, page)
}
func (repo RepositoryStruct) NftContractGetByNftToken(ctx context.Context, contractToken, token string) (*entity.NftContract, error) {
	return repo.Db.NftContractGetByNftToken(ctx, contractToken, token)
}

func (repo RepositoryStruct) NftCreate(ctx context.Context, bill *entity.Bill, nftData entity.NftData) error {
	var err error
	var owner *entity.NftOwner
	if bill.BillType == entity.BillType_Eip1155 { // 如果是1155类型需要判断是否存在，如果存在，则只需要添加修改的数量，而不需要创建
		owner, _ = repo.Db.NftOwnerGetByNftTokenAndOwnerAddr(ctx, bill.ContractRecord.ContractType, bill.ContractRecord.GameToken, config.SysAccountMap[bill.ContractRecord.ContractType].AddrIncome)
	}
	// 事务
	err = repo.Mongo.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			log.GetLogger().Error("[NftCreate] sessionContext.StartTransaction failed",
				zap.Error(err))
			return hcode.ErrInternalDb
		}
		err = repo.Db.BillCreate(sessionContext, bill)
		if err != nil {
			return err
		}
		// 创建nft owner
		switch bill.BillType {
		case entity.BillType_Eip721:
			var newOwner = &entity.NftOwner{
				Id:            primitive.NewObjectID(),
				GameId:        bill.ContractRecord.GameId,
				NftToken:      "",
				GameToken:     bill.ContractRecord.GameToken,
				CreateHash:    bill.Hash,
				ContractToken: bill.ContractRecord.ContractType,
				OwnerAddress:  bill.ToAddr,
				NftData:       nftData,
				Uid:           bill.Uid,
				Status:        entity.NftOwnerStatus_Available,
			}
			if err = repo.NftOwnerCreate(sessionContext, newOwner); err != nil {
				if err := sessionContext.AbortTransaction(sessionContext); err != nil {
					log.GetLogger().Error("[NftCreate] sessionContext.AbortTransaction failed",
						zap.Error(err))
					return hcode.ErrInternalDb
				}
				return err
			}
		case entity.BillType_Eip1155:
			if owner != nil { // 存在，添加数量
				if err = repo.Db.NftOwnerInventory(ctx, owner.ContractToken, owner.NftToken, owner.OwnerAddress, bill.ContractRecord.Num, true); err != nil {
					if err := sessionContext.AbortTransaction(sessionContext); err != nil {
						log.GetLogger().Error("[NftCreate] sessionContext.AbortTransaction failed",
							zap.Error(err))
						return hcode.ErrInternalDb
					}
					return err
				}
			} else { // 不存在创建
				var newOwner = &entity.NftOwner{
					Id:            primitive.NewObjectID(),
					GameId:        bill.ContractRecord.GameId,
					NftToken:      "",
					GameToken:     bill.ContractRecord.GameToken,
					CreateHash:    bill.Hash,
					ContractToken: bill.ContractRecord.ContractType,
					OwnerAddress:  bill.ToAddr,
					NftData:       nftData,
					Uid:           bill.Uid,
					Status:        entity.NftOwnerStatus_Available,
				}
				if err = repo.NftOwnerCreate(sessionContext, newOwner); err != nil {
					if err := sessionContext.AbortTransaction(sessionContext); err != nil {
						log.GetLogger().Error("[NftCreate] sessionContext.AbortTransaction failed",
							zap.Error(err))
						return hcode.ErrInternalDb
					}
					return err
				}
			}
		default:
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				log.GetLogger().Error("[NftCreate] sessionContext.AbortTransaction failed",
					zap.Error(err))
				return hcode.ErrInternalDb
			}
			return hcode.ErrBillType
		}

		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.GetLogger().Error("[NftCreate] sessionContext.CommitTransaction failed",
				zap.Error(err))
			return err
		}
		return nil
	})
	return err
}
