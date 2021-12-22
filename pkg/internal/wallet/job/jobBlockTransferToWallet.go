package job

import (
	"context"
	"encoding/json"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"tp_wallet/internal/block_chain/chain/scan/common"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/log"
)

// JobBlockTransferToWallet 处理链上数据到钱包
func (s *WalletJob) JobBlockTransferToWallet() {
	log.GetLogger().Debug("job JobBlockTransferToWallet Success ", zap.Any("Topic", s.KafkaConfig.TopicCurrencyTransaction))
	go func() {
		for err := range s.H2OConsumer.Errors() {
			log.GetLogger().Error("[JobBlockTransferToWallet] consumer Errors", zap.Error(err))
		}
	}()
	go func() {
		for {
			hand := KafkaH2OJob{
				walletSrv: s.WalletSrv,
			}
			if err := s.H2OConsumer.Consume(s.ConsumerCtx, []string{s.KafkaConfig.TopicCurrencyTransaction}, hand); err != nil {
				log.GetLogger().Error("[TransactionJob] consumer Consume", zap.Error(err))
			}
			if s.ConsumerCtx.Err() != nil {
				log.GetLogger().Error("[TransactionJob] consumer ctx", zap.Error(s.ConsumerCtx.Err()))
				return
			}
		}
	}()
}

type KafkaH2OJob struct {
	walletSrv walletPb.WalletSrvServer
}

func (consumer KafkaH2OJob) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer KafkaH2OJob) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer KafkaH2OJob) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var err error
	var ctx = context.TODO()
	for message := range claim.Messages() {
		var result common.H20TransactionForPush
		var receiveData common.PushInput
		err = json.Unmarshal(message.Value, &receiveData)
		if err != nil {
			log.GetLogger().Error("[ConsumeClaim] Unmarshal", zap.Error(err))
			sess.MarkMessage(message, "")
			continue
		}
		err = json.Unmarshal(receiveData.Data, &result)
		if err != nil {
			log.GetLogger().Error("[ConsumeClaim] Unmarshal", zap.Error(err))
			sess.MarkMessage(message, "")
			continue
		}
		log.GetLogger().Debug("Messages",
			zap.Any("msg", result),
			zap.Any("code", receiveData.Code),
			zap.Any("topic", message.Topic),
			zap.Any("offset", message.Offset),
			zap.Any("partition", message.Partition),
		)
		if len(result.From) == 0 || len(result.To) == 0 || len(result.Hash) == 0 || len(result.Amount) <= 0 {
			log.GetLogger().Error("[ConsumeClaim] Parameters failed", zap.Any("msg", result))
			sess.MarkMessage(message, "")
			continue
		}
		var newBill = &walletPb.BillInfo{
			BillType: int64(entity.BillType_Eip20),
			Hash:     result.Hash,
			FromAddr: result.From,
			ToAddr:   result.To,
			BalanceRecord: &walletPb.BalanceRecord{
				Currency:      result.Currency,
				ReceiveAmount: result.Amount,
			},
		}
		if result.Status == "success" {
			newBill.BillStatus = walletPb.BillStatus_Success
		} else {
			newBill.BillStatus = walletPb.BillStatus_Failed
		}
		if _, err = consumer.walletSrv.DealWithBill(ctx, newBill); err != nil {
			log.GetLogger().Error("[ConsumeClaim] failed", zap.Any("req", result), zap.Error(err))
		}
		sess.MarkMessage(message, "")
	}
	return nil
}
