package props

import (
	"context"
	"encoding/json"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"time"
	"tp_wallet/pkg/log"
)

// PushNftCreateSuccess 推送nft创建成功
func (ps PropsStruct) PushNftCreateSuccess(ctx context.Context, msg walletPb.Metadata) error {
	msg.TransferType = walletPb.TransferType_NftCreate
	var js, _ = json.Marshal(msg)
	var err error
	var partition int32
	var offset int64
	partition, offset, err = ps.Producer.SendMessage(&sarama.ProducerMessage{
		Topic:     ps.KafkaConfig.TopicProps,
		Value:     sarama.ByteEncoder(js),
		Timestamp: time.Now(),
	})
	if err != nil {
		log.GetLogger().Error("[PushNftCreateSuccess] failed", zap.Any("msg", msg), zap.Error(err))
		return err
	} else {
		log.GetLogger().Info("[PushNftCreateSuccess] success", zap.Int32("partition", partition), zap.Int64("offset", offset), zap.Any("msg", msg))
		return nil
	}
}

// PushNftCashSuccess 推送nft提现成功
func (ps PropsStruct) PushNftCashSuccess(ctx context.Context, msg walletPb.Metadata) error {
	msg.TransferType = walletPb.TransferType_NftCASH
	var js, _ = json.Marshal(msg)
	var err error
	var partition int32
	var offset int64
	partition, offset, err = ps.Producer.SendMessage(&sarama.ProducerMessage{
		Topic:     ps.KafkaConfig.TopicProps,
		Value:     sarama.ByteEncoder(js),
		Timestamp: time.Now(),
	})
	if err != nil {
		log.GetLogger().Error("[PushNftCreateSuccess] failed", zap.Any("msg", msg), zap.Error(err))
		return err
	} else {
		log.GetLogger().Info("[PushNftCreateSuccess] success", zap.Int32("partition", partition), zap.Int64("offset", offset), zap.Any("msg", msg))
		return nil
	}
}

// PushNftCharge 推送nft充值
func (ps PropsStruct) PushNftCharge(ctx context.Context, msg walletPb.Metadata) error {
	msg.TransferType = walletPb.TransferType_NftCHARGE
	var js, _ = json.Marshal(msg)
	var err error
	var partition int32
	var offset int64
	partition, offset, err = ps.Producer.SendMessage(&sarama.ProducerMessage{
		Topic:     ps.KafkaConfig.TopicProps,
		Value:     sarama.ByteEncoder(js),
		Timestamp: time.Now(),
	})
	if err != nil {
		log.GetLogger().Error("[PushNftCreateSuccess] failed", zap.Any("msg", msg), zap.Error(err))
		return err
	} else {
		log.GetLogger().Info("[PushNftCreateSuccess] success", zap.Int32("partition", partition), zap.Int64("offset", offset), zap.Any("msg", msg))
		return nil
	}
}

func (ps PropsStruct) Close() {
	ps.ConsumerCancelFunc()
}
