package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"tp_wallet/internal/tp_wallet/domain/entity"
	"tp_wallet/pkg/hcode"
	"tp_wallet/pkg/log"
)

var _ repo.BillRepo = (*Bill)(nil)

type Bill struct {
	repository
}

func (b *Bill) getWalletTable() *mongo.Collection {
	return b.mgo.Database("gs").Collection("bill")
}

func (b *Bill) CreateBill(ctx context.Context, bill *entity.Bill) error {
	if _, err := b.getWalletTable().InsertOne(ctx, bill); err != nil {
		log.GetLogger().Error("[CreateBill] InsertOne", zap.Any("bill", bill), zap.Error(err))
		return hcode.ErrInternalDb
	}
	return nil
}

func (b *Bill) UpdateBill(ctx context.Context, bill *entity.Bill) error {
	filter, updates := bill.GetUpdates()
	res, err := b.getWalletTable().UpdateOne(ctx, filter, updates)
	if err != nil {
		log.GetLogger().Error("[UpdateBill] UpdateOne", zap.Any("bill", bill), zap.Error(err))
		return hcode.ErrInternalDb
	}
	if res.ModifiedCount != 1 || res.MatchedCount != 1 {
		return hcode.ErrInternalDb
	}
	return nil
}

func (b *Bill) RemoveBill(ctx context.Context, repo repo.BillSpecificationRepo) error {
	if err := repo.Check(ctx); err != nil {
		return err
	}
	res, err := b.getWalletTable().DeleteOne(ctx, repo.ToSql(ctx))
	if err != nil {
		log.GetLogger().Error("[RemoveBill] DeleteOne", zap.Any("repo", repo.ToSql(ctx)), zap.Error(err))
		return hcode.ErrInternalDb
	}
	if res.DeletedCount != 1 {
		return hcode.ErrInternalDb
	}
	return nil
}

func (b *Bill) QueryBillSingle(ctx context.Context, repo repo.BillSpecificationRepo) (bill *entity.Bill, err error) {
	if err = repo.Check(ctx); err != nil {
		return
	}
	bill = new(entity.Bill)
	if err = b.getWalletTable().FindOne(ctx, repo.ToSql(ctx)).Decode(bill); err != nil {
		log.GetLogger().Error("[QueryBillSingle] FindOne", zap.Any("data", repo.ToSql(ctx)), zap.Error(err))
		err = hcode.ErrInternalDb
		return
	}
	return
}

func (b *Bill) QueryBills(ctx context.Context, repo repo.BillSpecificationRepo) ([]*entity.Bill, error) {
	if err := repo.Check(ctx); err != nil {
		return nil, err
	}
	var (
		bills []*entity.Bill
	)
	cursor, err := b.getWalletTable().Find(ctx, repo.ToSql(ctx))
	if err != nil {
		log.GetLogger().Error("[QueryBills] Find", zap.Any("data", repo.ToSql(ctx)), zap.Error(err))
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		log.GetLogger().Error("[QueryBills] cursor.Err()", zap.Any("data", repo.ToSql(ctx)), zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		data := &entity.Bill{}
		if err = cursor.Decode(data); err != nil {
			log.GetLogger().Error("[QueryBills] cursor.Decode()", zap.Any("data", repo.ToSql(ctx)), zap.Error(err))
			continue
		}
		bills = append(bills, data)
	}
	return bills, nil
}
