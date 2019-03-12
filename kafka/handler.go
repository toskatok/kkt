package kafka

import (
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type kktConsumerGroupHandler struct {
	lck             *sync.RWMutex
	repeatedOffsets map[string]bool
}

func (h kktConsumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	logrus.Infof("Hello world of kafka in generation (%d)", sess.GenerationID())
	return nil
}
func (h kktConsumerGroupHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	logrus.Infof("Bye bye kafka in generation (%d)", sess.GenerationID())
	return nil
}
func (h kktConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		logrus.Infof("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)

		h.lck.RLock()
		ok := h.repeatedOffsets[fmt.Sprintf("%d:%d", msg.Partition, msg.Offset)]
		h.lck.RUnlock()

		if ok {
			logrus.Errorf("Duplicated message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		} else {
			h.lck.Lock()
			h.repeatedOffsets[fmt.Sprintf("%d:%d", msg.Partition, msg.Offset)] = true
			h.lck.Unlock()
		}

		logrus.Debugln(string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}
