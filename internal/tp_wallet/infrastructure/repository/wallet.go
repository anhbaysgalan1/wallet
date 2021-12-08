package repository

import (
	"context"
	"gluttonous/internal/wallet/domain/entity"
	"gluttonous/internal/wallet/domain/ob"
	"gluttonous/internal/wallet/domain/repo"
	walletTool "gluttonous/internal/wallet/infrastructure/tool"
	"gluttonous/pkg/hcode"
	"gluttonous/pkg/log"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"
)

var _ repo.WalletRepo = (*Wallet)(nil)

type Wallet struct {
	repository
}

func (w *Wallet) getWalletTable() *mongo.Collection {
	return w.mgo.Database("gs").Collection("Wallet")
}

func (w *Wallet) CreateWallet(ctx context.Context, wallet *entity.Wallet) error {
	if _, err := w.getWalletTable().InsertOne(ctx, wallet); err != nil {
		log.GetLogger().Error("[CreateWallet] InsertOne", zap.Any("wallet", wallet), zap.Error(err))
		return hcode.MgoExecErr
	}
	return nil
}

func (w *Wallet) UpdateWallet(ctx context.Context, wallet *entity.Wallet) error {
	filter, updates := wallet.GetUpdates()
	res, err := w.getWalletTable().UpdateOne(ctx, filter, updates)
	if err != nil {
		log.GetLogger().Error("[UpdateWallet] UpdateOne", zap.Any("wallet", wallet), zap.Error(err))
		return hcode.MgoExecErr
	}
	if res.ModifiedCount != 1 || res.MatchedCount != 1 {
		return hcode.MgoExecErr
	}
	return nil
}

func (w *Wallet) RemoveWallet(ctx context.Context, repo repo.WalletSpecificationRepo) error {
	if err := repo.Check(ctx); err != nil {
		return err
	}
	res, err := w.getWalletTable().DeleteOne(ctx, repo.ToSql(ctx))
	if err != nil {
		log.GetLogger().Error("[RemoveWallet] DeleteOne", zap.Any("repo", repo.ToSql(ctx)), zap.Error(err))
		return hcode.MgoExecErr
	}
	if res.DeletedCount != 1 {
		return hcode.ResourcesNotFindErr
	}
	return nil
}

func (w *Wallet) QueryWalletSingle(ctx context.Context, repo repo.WalletSpecificationRepo) (wallet *entity.Wallet, err error) {
	if err = repo.Check(ctx); err != nil {
		return
	}
	wallet = new(entity.Wallet)
	if err = w.getWalletTable().FindOne(ctx, repo.ToSql(ctx)).Decode(wallet); err != nil {
		log.GetLogger().Error("[QueryWalletSingle] FindOne", zap.Any("data", repo.ToSql(ctx)), zap.Error(err))
		err = hcode.MgoExecErr
		return
	}
	return
}

func (w *Wallet) QueryWallets(ctx context.Context, repo repo.WalletSpecificationRepo) ([]*entity.Wallet, error) {
	if err := repo.Check(ctx); err != nil {
		return nil, err
	}
	var (
		wallet []*entity.Wallet
	)
	cursor, err := w.getWalletTable().Find(ctx, repo.ToSql(ctx))
	if err != nil {
		log.GetLogger().Error("[QueryWallets] Find", zap.Any("data", repo.ToSql(ctx)), zap.Error(err))
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		log.GetLogger().Error("[QueryWallets] cursor.Err()", zap.Any("data", repo.ToSql(ctx)), zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		data := &entity.Wallet{}
		if err = cursor.Decode(data); err != nil {
			log.GetLogger().Error("[QueryWallets] cursor.Decode()", zap.Any("data", repo.ToSql(ctx)), zap.Error(err))
			continue
		}
		wallet = append(wallet, data)
	}
	return wallet, nil
}

func (w *Wallet) SetWalletSys(ctx context.Context, wallet *ob.SysWallet, amount string) error {
	key, field := wallet.GetKeyAndField()
	err := w.rds.HSet(key, field, amount).Err()
	if err != nil {
		return err
	}
	return nil
}

func (w *Wallet) UpdateWalletSys(ctx context.Context, wallet *ob.SysWallet, amount string, isAdd bool) error {
	balance, err := w.GetWalletSysSingle(ctx, wallet)
	if err != nil {
		return err
	}
	var lastBalance string
	if isAdd {
		lastBalance, err = walletTool.AddBalance(balance, amount)
		if err != nil {
			return err
		}
	} else {
		lastBalance, err = walletTool.EnoughBalance(balance, amount)
		if err != nil {
			return err
		}
	}
	return w.SetWalletSys(ctx, wallet, lastBalance)
}

func (w *Wallet) GetWalletSysSingle(ctx context.Context, wallet *ob.SysWallet) (string, error) {
	key, field := wallet.GetKeyAndField()
	amount, err := w.rds.HGet(key, field).Result()
	if err != nil {
		return "", err
	}
	return amount, nil
}
