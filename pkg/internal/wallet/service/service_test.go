package service

import (
	"context"
	"flag"
	"git.huoys.com/middle-end/kratos/pkg/conf/paladin"
	"go.uber.org/zap"
	"os"
	"testing"
	"tp_wallet/config"
	adapter2 "tp_wallet/internal/block_chain/adapter"
	"tp_wallet/internal/wallet/adapter/block_chain"
	"tp_wallet/internal/wallet/adapter/props"
	repository2 "tp_wallet/internal/wallet/repository"
	"tp_wallet/pkg/database/mongo"
	"tp_wallet/pkg/database/redis"
	"tp_wallet/pkg/log"
	"tp_wallet/pkg/redisCache/common"
	"tp_wallet/pkg/tool"
)

var walletSrv WalletSrv
var ctx = context.Background()

func TestMain(m *testing.M) {
	var err error
	flag.Set("conf", "../../../configs")
	flag.Parse()
	if err = paladin.Init(); err != nil {
		panic(err)
	}
	config.NewWalletBusiness()
	config.NewBlockBusiness()
	config.NewConfigFee()
	config.NewConfigCurrency()
	redisConfig, err := config.ConfRedis()
	if err != nil {
		panic(err)
	}
	client, err := redis.NewRedis(redisConfig)
	if err != nil {
		panic(err)
	}
	mongoConfig, err := config.ConfNewDB()
	if err != nil {
		panic(err)
	}
	configTransferKafka, err := config.NewConfigBlockKafka()
	if err != nil {
		panic(err)
	}
	mongoClient := mongo.NewMongo(mongoConfig)
	//blockChainSrv := adapter2.NewTestBlockChainSrv()
	blockChainSrv := adapter2.NewBlockChainSrv()
	block_chainBlockChainSrv := block_chain.NewBlockChainSrv(blockChainSrv)
	propsSrv, err := props.NewPropsSrv(configTransferKafka)
	if err != nil {
		panic(err)
	}
	repository := repository2.NewWalletRepository(client, mongoClient, block_chainBlockChainSrv, propsSrv)
	walletSrv = WalletSrv{
		Repo:          repository,
		BlockChainSrv: block_chainBlockChainSrv,
		Lock:          common.NewRedisLock(client),
	}
	// 初始化雪花算法
	if snowFlakeId, err := tool.NewWorkerIdTool(client); err != nil {
		log.GetLogger().Error("[Server Init] tool.NewWorkerIdTool failed", zap.Error(err))
		panic(err)
	} else {
		if err := tool.NewWorker(snowFlakeId); err != nil {
			log.GetLogger().Error("[Server Init] tool.NewWorker failed", zap.Int64("id", snowFlakeId),
				zap.Error(err))
			panic(err)
		}
		log.GetLogger().Info("tool.NewWorkerIdTool snow flake init success", zap.Int64("id", snowFlakeId))
	}
	ret := m.Run()
	os.Exit(ret)
}
