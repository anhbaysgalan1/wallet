package db

import (
	"context"
	"errors"
	"github.com/leaf-rain/wallet/internal/block_chain/model"
	"github.com/leaf-rain/wallet/pkg/context_db"
	"github.com/leaf-rain/wallet/pkg/hcode"
	"github.com/leaf-rain/wallet/pkg/log"
	"github.com/leaf-rain/wallet/pkg/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (r *repository) AddressCreate(ctx context.Context, addr *model.AddressPrivate) (err error) {
	mg := context_db.GetMongoTx(ctx)
	if mg == nil {
		mg = r.mongo
	}
	if addr == nil || len(addr.Address) == 0 || len(addr.Private) == 0 {
		return hcode.ErrInternalParameter
	}
	if addr.Id.IsZero() {
		addr.Id = primitive.NewObjectID()
	}
	addr.Status = model.AddressStatus_UnUsed
	collection := addressPrivateGetCollection(mg)
	var now = tool.GetTimeUnixMilli()
	addr.CreateTime, addr.UpdateTime = now, now
	if _, err := collection.InsertOne(ctx, addr); err != nil {
		log.GetLogger().Error("[AddressCreate] InsertOne failed", zap.Any("req", addr), zap.Error(err))
		return err
	}
	return err
}

func (r *repository) AddressSet(ctx context.Context, addr *model.AddressPrivate) (err error) {
	mg := context_db.GetMongoTx(ctx)
	if mg == nil {
		mg = r.mongo
	}
	collection := addressPrivateGetCollection(mg)
	var update = bson.M{}
	if addr.Status != 0 {
		update["status"] = addr.Status
	}
	if len(addr.Remarks) != 0 {
		update["remarks"] = addr.Remarks
	}
	update["update_time"] = tool.GetTimeUnixMilli()
	result, err := collection.UpdateOne(ctx, bson.M{"_id": addr.Id}, update)
	if err != nil {
		log.GetLogger().Error("[AddressSet] failed", zap.Any("req", addr), zap.Error(err))
		return hcode.ErrServer
	}
	if result.ModifiedCount != 1 {
		if err != nil {
			log.GetLogger().Error("[DbSetAddressStatus] failed", zap.Any("req", addr), zap.Error(err))
			return hcode.ErrDbExec
		}
	}
	return nil
}

func (r *repository) AddressSCreate(ctx context.Context, addrS []model.AddressPrivate) (err error) {
	mg := context_db.GetMongoTx(ctx)
	if mg == nil {
		mg = r.mongo
	}
	collection := addressPrivateGetCollection(mg)
	var value = make([]interface{}, len(addrS))
	for i, item := range addrS {
		value[i] = item
	}
	result, err := collection.InsertMany(ctx, value)
	if err != nil {
		log.GetLogger().Error("[GetAddressByPage] failed", zap.Any("addrS", addrS), zap.Error(err))
		return err
	}
	log.GetLogger().Info("[AddressSCreate] success", zap.Any("success total:", len(result.InsertedIDs)))
	return nil
}

func (r *repository) AddressGetByFilter(ctx context.Context, filter bson.D, opt *options.FindOptions) (addrS []model.AddressPrivate, err error) {
	mg := context_db.GetMongoTx(ctx)
	if mg == nil {
		mg = r.mongo
	}
	if filter == nil {
		filter = bson.D{}
	}
	collection := addressPrivateGetCollection(mg)
	cursor, err := collection.Find(ctx, filter, opt)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.GetLogger().Error("[GetAddressByFilter] failed", zap.Any("filter", filter), zap.Any("opt", opt), zap.Error(err))
		}
		return addrS, err
	}
	//延迟关闭游标
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			log.GetLogger().Error("[GetAddressByFilter] failed", zap.Any("filter", filter), zap.Any("opt", opt), zap.Error(err))
		}
	}()
	err = cursor.All(ctx, &addrS)
	if err != nil {
		log.GetLogger().Error("[GetAddressByFilter] failed", zap.Any("filter", filter), zap.Any("opt", opt), zap.Error(err))
		return addrS, err
	}
	return addrS, nil
}
