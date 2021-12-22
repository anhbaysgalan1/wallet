package props

import (
	"context"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"github.com/Shopify/sarama"
	"tp_wallet/internal/common"
)

type PropsSrv interface {
	// PushNftCreateSuccess 推送nft创建成功
	PushNftCreateSuccess(ctx context.Context, msg walletPb.Metadata) error
	// PushNftCashSuccess 推送nft提现成功
	PushNftCashSuccess(ctx context.Context, msg walletPb.Metadata) error
	// PushNftCharge 推送nft充值
	PushNftCharge(ctx context.Context, msg walletPb.Metadata) error

	Close()
}

// PropsStruct 此路由为kafka通知道具服修改道具状态
type PropsStruct struct {
	ConsumerCtx        context.Context
	ConsumerCancelFunc context.CancelFunc
	KafkaConfig        *common.ConfigTransferKafka
	Producer           sarama.SyncProducer
}

func NewPropsSrv(kafkaConfig *common.ConfigTransferKafka) (PropsSrv, error) {
	var ctx, cancel = context.WithCancel(context.Background())
	var srv = PropsStruct{
		ConsumerCtx:        ctx,
		ConsumerCancelFunc: cancel,
		KafkaConfig:        kafkaConfig,
	}
	var err error
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true
	conf.Producer.Retry.Max = 1000
	conf.Version = sarama.V2_2_0_0
	srv.Producer, err = sarama.NewSyncProducer(kafkaConfig.KafkaAddr, conf)
	if err != nil {
		return nil, err
	}
	return &srv, nil
}
