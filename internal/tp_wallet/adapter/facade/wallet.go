package anticorrosive

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	blockDto "tp_wallet/internal/block_chain/dto"
	walletCommon "tp_wallet/internal/common"
	"tp_wallet/internal/wallet/dto"
	walletJob "tp_wallet/internal/wallet/job"
	"tp_wallet/internal/wallet/repository"
	"tp_wallet/internal/wallet/repository/cache"
	"tp_wallet/internal/wallet/repository/db"
	walletSrv "tp_wallet/internal/wallet/service"
	"tp_wallet/pkg/database/redis"
	redisCommon "tp_wallet/pkg/redisCache/common"
)

func NewWalletRepository(r *redis.Client, d *mongo.Client, bSrv blockDto.BlockChainSrv) repository.Repository {
	repo := &repository.RepositoryStruct{}
	repo.Db = db.NewWalletDb(d)
	repo.Cache = cache.NewCache(r)
	repo.Lock = redisCommon.NewRedisLock(r)
	repo.BlockChainSrv = bSrv
	return repo
}

func NewWalletSrv(r *redis.Client, repo repository.Repository, bSrv blockDto.BlockChainSrv) dto.WalletSrvServer {
	return &walletSrv.WalletSrv{Repo: repo, BlockChainSrv: bSrv, Lock: redisCommon.NewRedisLock(r)}
}

func NewWalletJobSrv(redis *redis.Client, d *mongo.Client, walletSrv dto.WalletSrvServer, confBlockKafka *walletCommon.ConfigTransferKafka) dto.WalletJob {
	var err error
	job := &walletJob.WalletJob{
		WalletSrv: walletSrv,
	}
	job.ConsumerCtx, job.ConsumerCancelFunc = context.WithCancel(context.Background())
	h2oConf := sarama.NewConfig()
	h2oConf.Version = sarama.V2_2_0_0
	h2oConf.Consumer.Return.Errors = true
	h2oConf.Consumer.Offsets.Initial = sarama.OffsetNewest
	job.H2OConsumer, err = sarama.NewConsumerGroup(confBlockKafka.KafkaAddr, confBlockKafka.GroupH2OTransaction, h2oConf)
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
