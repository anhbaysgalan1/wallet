package job

import (
	"context"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"github.com/Shopify/sarama"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"tp_wallet/internal/common"
	"tp_wallet/internal/wallet/adapter/props"
	walletSrv "tp_wallet/internal/wallet/service"
	"tp_wallet/pkg/database/redis"
	"tp_wallet/pkg/log"
	redisCommon "tp_wallet/pkg/redisCache/common"
)

type WalletJob struct {
	Cron                    *cron.Cron
	WalletSrv               walletPb.WalletSrvServer
	WalletTransferToBlockId cron.EntryID
	BlockTransferToWallet   cron.EntryID
	H2OConsumer             sarama.ConsumerGroup
	NftTransferConsumer     sarama.ConsumerGroup
	NftCreateConsumer       sarama.ConsumerGroup
	ConsumerCtx             context.Context
	ConsumerCancelFunc      context.CancelFunc
	KafkaConfig             *common.ConfigTransferKafka
	Lock                    redisCommon.RedisLock
	PropsSrv                props.PropsSrv
}

// 这里只提供给测试方法使用
func newWalletJobSrv(redis *redis.Client, walletSrv walletSrv.WalletSrv, confBlockKafka *common.ConfigTransferKafka) WalletJob {
	var err error
	job := WalletJob{
		WalletSrv: &walletSrv,
	}
	job.ConsumerCtx, job.ConsumerCancelFunc = context.WithCancel(context.Background())
	h2oConf := sarama.NewConfig()
	h2oConf.Version = sarama.V2_2_0_0
	h2oConf.Consumer.Return.Errors = true
	h2oConf.Consumer.Offsets.Initial = sarama.OffsetNewest
	job.H2OConsumer, err = sarama.NewConsumerGroup(confBlockKafka.KafkaAddr, confBlockKafka.GroupCurrencyTransaction, h2oConf)
	nftTransferConf := sarama.NewConfig()
	nftTransferConf.Version = sarama.V2_2_0_0
	nftTransferConf.Consumer.Return.Errors = true
	nftTransferConf.Consumer.Offsets.Initial = sarama.OffsetNewest
	job.NftTransferConsumer, err = sarama.NewConsumerGroup(confBlockKafka.KafkaAddr, confBlockKafka.GroupNftTransaction, nftTransferConf)
	nftCreateConf := sarama.NewConfig()
	nftCreateConf.Version = sarama.V2_2_0_0
	nftCreateConf.Consumer.Return.Errors = true
	nftCreateConf.Consumer.Offsets.Initial = sarama.OffsetNewest
	job.NftCreateConsumer, err = sarama.NewConsumerGroup(confBlockKafka.KafkaAddr, confBlockKafka.GroupNftCreate, nftCreateConf)
	if err != nil {
		panic(err)
	}
	job.Cron = cron.New(cron.WithSeconds())
	job.KafkaConfig = confBlockKafka
	job.Lock = redisCommon.NewRedisLock(redis)
	return job
}

func (wj WalletJob) Run() {
	var err error
	var ctx = context.Background()
	// 异步上传区块链
	_, err = wj.Cron.AddFunc("*/5 * * * * ?", func() {
		wj.JobCurrencyWalletTransferToBlock(ctx)
	})
	if err != nil {
		log.GetLogger().Error("[JobStart] JobWalletTransferToBlock job start failed", zap.Error(err))
	} else {
		log.GetLogger().Info("[JobStart] JobWalletTransferToBlock job start success")
	}
	_, err = wj.Cron.AddFunc("*/5 * * * * ?", func() {
		wj.JobNftWalletTransferToBlock(ctx)
	})
	if err != nil {
		log.GetLogger().Error("[JobStart] JobNftWalletTransferToBlock job start failed", zap.Error(err))
	} else {
		log.GetLogger().Info("[JobStart] JobNftWalletTransferToBlock job start success")
	}
	wj.JobBlockTransferToWallet()
	wj.JobNftCharge()
	wj.JobNftCreate()
	go wj.Cron.Start()
}

func (wj WalletJob) Close() {
	wj.Cron.Stop()
	wj.ConsumerCancelFunc()
	wj.H2OConsumer.Close()
	wj.NftTransferConsumer.Close()
	wj.NftCreateConsumer.Close()
}
