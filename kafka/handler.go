package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type kktConsumerGroupHandler struct {
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
		logrus.Debugln(string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}
