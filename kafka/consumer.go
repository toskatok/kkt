package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

// ConsumerGroup contains sarama group consumer plus the required information for creating it.
type ConsumerGroup struct {
	Topic         string
	ConsumerGroup sarama.ConsumerGroup

	done chan struct{}
}

// NewConsumerGroup creates new consumer group based on given configuration
func NewConsumerGroup(topic string, brokers []string, clientID string) (*ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.ClientID = clientID
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Return.Errors = true

	gc, err := sarama.NewConsumerGroup(brokers, fmt.Sprintf("%s-%s", "kkt", topic), config)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		Topic:         topic,
		ConsumerGroup: gc,
		done:          make(chan struct{}),
	}, nil
}

// Run runs the consumer group to consume from kafka
func (gc *ConsumerGroup) Run() {
	topics := []string{gc.Topic}
	handler := kktConsumerGroupHandler{
		lck:             &sync.RWMutex{},
		repeatedOffsets: make(map[string]bool),
	}

	go func() {
		for err := range gc.ConsumerGroup.Errors() {
			logrus.Errorf("sarama: %s", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := gc.ConsumerGroup.Consume(ctx, topics, handler)
		if err != nil {
			logrus.Errorf("samara consumer group: %s", err)
		}
		if err := gc.ConsumerGroup.Close(); err != nil {
			logrus.Errorf("samara consumer group close error: %s", err)
		}
		close(gc.done)
	}()

	<-gc.done
	cancel()
}

// Exit closes the consumer group and exits from its loop
func (gc *ConsumerGroup) Exit() {
	gc.done <- struct{}{}
	<-gc.done
}
