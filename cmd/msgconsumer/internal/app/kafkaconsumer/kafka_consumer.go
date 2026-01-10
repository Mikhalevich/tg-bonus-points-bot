package kafkaconsumer

import (
	"context"
	"errors"
	"sync"

	"github.com/segmentio/kafka-go"

	"github.com/Mikhalevich/tg-coffee-shop-bot/cmd/msgconsumer/internal/config"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

type KafkaConsumer struct {
	reader      *kafka.Reader
	workerCount int
}

func New(cfg config.Kafka) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.Brokers,
			Topic:   cfg.Topic,
			GroupID: cfg.ConsumerGroup,
		}),
		workerCount: cfg.WorkersCount,
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

	var (
		dataChan = make(chan kafka.Message)
		wgr      = k.runWorkers(ctx, dataChan, processFn)
	)

	for {
		msg, err := k.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}

			logger.FromContext(ctx).
				WithError(err).
				Error("error while read message")
		}

		dataChan <- msg
	}

	close(dataChan)
	wgr.Wait()

	return nil
}

func (k *KafkaConsumer) runWorkers(
	ctx context.Context,
	dataChan <-chan kafka.Message,
	processFn func(ctx context.Context, payload []byte) error,
) *sync.WaitGroup {
	var wgr sync.WaitGroup

	for range k.workerCount {
		wgr.Add(1)

		go func() {
			defer wgr.Done()

			for msg := range dataChan {
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

				if err := k.reader.CommitMessages(ctx, msg); err != nil {
					logger.FromContext(ctx).
						WithFields(logger.Fields{
							"kafka_topic":     msg.Topic,
							"kafka_partition": msg.Partition,
							"kafka_offset":    msg.Offset,
							"kafka_key":       string(msg.Key),
							"kafka_msg":       string(msg.Value),
						}).
						WithError(err).
						Error("commit message")
				}
			}
		}()
	}

	return &wgr
}
