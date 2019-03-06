package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type kktConsumerGroupHandler struct {
}

func (h kktConsumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	logrus.Infof("Hello world of kafka (%s) in generation (%d)\n", "ride", sess.GenerationID())
	return nil
}
func (h kktConsumerGroupHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	logrus.Infof("Bye bye kafka (%s) in generation (%d)\n", "ride", sess.GenerationID())
	return nil
}
func (h kktConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		logrus.Infof("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		logrus.Debugln(msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}
