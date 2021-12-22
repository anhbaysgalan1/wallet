package operation

import (
	"testing"
)

func TestPushMongoTransaction(t *testing.T) {
	pd, err := NewProducer([]string{"172.12.12.165:9092"})
	if err != nil {
		panic(err)
	}

	PushMongoTransaction(pd, TopicH20ForScan, []byte("hello"))
	PushMongoTransaction(pd, TopicNftForScan, []byte("hello"))
}
