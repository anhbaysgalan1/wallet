package job

import (
	"context"
	"encoding/json"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	scanCommon "tp_wallet/internal/block_chain/chain/scan/common"
	walletCommon "tp_wallet/internal/common"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/log"
)

// JobNftCharge 处理链上Nft到钱包
func (s *WalletJob) JobNftCharge() {
	log.GetLogger().Debug("job JobNftCharge Success ", zap.Any("Topic", s.KafkaConfig.TopicNftTransaction))
	go func() {
		for err := range s.NftTransferConsumer.Errors() {
			log.GetLogger().Error("[JobNftCharge] consumer Errors", zap.Error(err))
		}
	}()
	go func() {
		for {
			hand := KafkaNftJob{
				walletSrv: s.WalletSrv,
			}
			if err := s.NftTransferConsumer.Consume(s.ConsumerCtx, []string{s.KafkaConfig.TopicNftTransaction}, hand); err != nil {
				log.GetLogger().Error("[JobNftCharge] consumer Consume", zap.Error(err))
			}
			if s.ConsumerCtx.Err() != nil {
				log.GetLogger().Error("[JobNftCharge] consumer ctx", zap.Error(s.ConsumerCtx.Err()))
				return
			}
		}
	}()
}

type KafkaNftJob struct {
	walletSrv walletPb.WalletSrvServer
}

func (consumer KafkaNftJob) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer KafkaNftJob) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer KafkaNftJob) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var err error
	var ctx = context.TODO()
	for message := range claim.Messages() {
		var result scanCommon.Nft721TransactionForPush
		var receiveData scanCommon.PushInput
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
		if len(result.From) == 0 || len(result.To) == 0 || len(result.Hash) == 0 || len(result.Contract) == 0 || len(result.NftToken) == 0 {
			log.GetLogger().Error("[ConsumeClaim] Parameters failed", zap.Any("msg", result))
			sess.MarkMessage(message, "")
			continue
		}
		var newBill = &walletPb.BillInfo{
			BillType: int64(entity.BillType_Eip721),
			Hash:     result.Hash,
			FromAddr: result.From,
			ToAddr:   result.To,

			ContractRecord: &walletPb.ContractRecord{
				ContractAddr: result.Contract,
				NftToken:     result.NftToken,
			},
		}
		switch receiveData.Code {
		case scanCommon.H2OTransferCode:
			newBill.TransferType = walletPb.TransferType_CurrencyTransfer
		case scanCommon.H2OApprovalCode:
			newBill.TransferType = walletPb.TransferType_CurrencyApproval
		case scanCommon.H2OTransferFromCode:
			newBill.TransferType = walletPb.TransferType_CurrencyApprovalTransfer
		case scanCommon.RowingNftCreateCode:
			newBill.TransferType = walletPb.TransferType_NftCreate
		case scanCommon.RowingNftTransferCode:
			newBill.TransferType = walletPb.TransferType_NftTransfer
		case scanCommon.RowingNftApprovalCode:
			newBill.TransferType = walletPb.TransferType_NftApproval
		case scanCommon.RowingNftTransferFromCode:
			newBill.TransferType = walletPb.TransferType_NftApprovalTransfer
		case scanCommon.RacerNftCreateCode:
			newBill.TransferType = walletPb.TransferType_NftCreate
		case scanCommon.RacerNftTransferCode:
			newBill.TransferType = walletPb.TransferType_NftTransfer
		case scanCommon.RacerNftApprovalCode:
			newBill.TransferType = walletPb.TransferType_NftApproval
		case scanCommon.RacerNftTransferFromCode:
			newBill.TransferType = walletPb.TransferType_NftApprovalTransfer
		case scanCommon.UnknownCode:
			newBill.TransferType = walletPb.TransferType_TransferType_ABANDON
		}
		if result.Status == walletCommon.BlockSuccess {
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
