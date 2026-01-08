package kafkaconsumer

import (
	"context"
	"errors"

	"github.com/segmentio/kafka-go"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func New(cfg config.Kafka) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.Brokers,
			Topic:   cfg.Topic,
			GroupID: cfg.ConsumerGroup,
		}),
	}
}

func (k *KafkaConsumer) Consume(
	ctx context.Context,
	processFn func(ctx context.Context, payload []byte) error,
) error {
	defer func() {
		if err := k.reader.Close(); err != nil {
			logger.FromContext(ctx).WithError(err).Error("close kafka reader")
		}
	}()

	for {
		msg, err := k.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			logger.FromContext(ctx).
				WithError(err).
				Error("error while read message")
		}

		if err := processFn(ctx, msg.Value); err != nil {
			logger.FromContext(ctx).
				WithFields(logger.Fields{
					"kafka_topic":     msg.Topic,
					"kafka_partition": msg.Partition,
					"kafka_offset":    msg.Offset,
					"kafka_key":       string(msg.Key),
					"kafka_msg":       string(msg.Value),
				}).
				WithError(err).
				Error("error while process message")
		}
	}
}
