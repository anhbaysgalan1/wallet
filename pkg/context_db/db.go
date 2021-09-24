package context_db

import (
	"context"
	"github.com/leaf-rain/wallet/common"
	"github.com/leaf-rain/wallet/pkg/database/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetMongoTx(ctx context.Context, client *mongo.Client) context.Context {
	return context.WithValue(ctx, common.MongoTx, client)
}

func GetMongoTx(ctx context.Context) (mc *mongo.Client) {
	if txDB, ok := ctx.Value(common.MongoTx).(*mongo.Client); ok {
		return txDB
	}
	return nil
}

func SetRedisTx(ctx context.Context, rd *redis.Client) context.Context {
	return context.WithValue(ctx, common.MongoTx, rd)
}

func GetRedisTx(ctx context.Context) (rd *redis.Client) {
	if txDB, ok := ctx.Value(common.MongoTx).(*redis.Client); ok {
		return txDB
	}
	return nil
}
