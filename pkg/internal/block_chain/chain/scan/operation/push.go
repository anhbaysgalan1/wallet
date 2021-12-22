package operation

import (
	"github.com/Shopify/sarama"
)

var (
	Producer              sarama.SyncProducer
	TopicH20ForScan       string
	TopicNftForScan       string
	TopicCreateNftForScan string
)

func NewProducer(ip []string) (sarama.SyncProducer, error) {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true
	conf.Producer.Retry.Max = 1000
	conf.Version = sarama.V2_2_0_0
	producer, err := sarama.NewSyncProducer(ip, conf)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func PushMongoTransaction(producer sarama.SyncProducer, topic string, data []byte) error {
	_, _, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
