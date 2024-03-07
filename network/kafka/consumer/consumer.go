package reader

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Handler func(kafka.Message) error

type Consumer struct {
	Bootstrap string
	GroupID   string
	Topic     string
	handler   Handler
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewConsumer(Bootstrap, GroupID, Topic string, handler Handler) (*Consumer, error) {
	if len(Bootstrap) <= 0 {
		return nil, fmt.Errorf("invalid Bootstrap (%s)", Bootstrap)
	}

	if len(GroupID) <= 0 {
		return nil, fmt.Errorf("invalid GroupID (%s)", GroupID)
	}

	if len(Topic) <= 0 {
		return nil, fmt.Errorf("invalid Topic (%s)", Topic)
	}

	if handler == nil {
		return nil, fmt.Errorf("invalid handler")
	}

	ctx, cancel := context.WithCancel(context.TODO())
	return &Consumer{
		Bootstrap: Bootstrap,
		GroupID:   GroupID,
		Topic:     Topic,
		handler:   handler,
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (c *Consumer) Open() error {
	c.wg.Add(1)
	go c.run()

	return nil
}

func (c *Consumer) Close() {
	c.cancel()
	c.wg.Wait()
}

func (c *Consumer) newReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:               strings.Split(c.Bootstrap, ","),
		GroupID:               c.GroupID,
		Topic:                 c.Topic,
		MinBytes:              10e3,
		MaxBytes:              10e6,
		MaxWait:               1 * time.Second,
		ReadLagInterval:       -1,
		RebalanceTimeout:      1 * time.Second,
		ErrorLogger:           kafka.LoggerFunc(func(msg string, a ...any) { fmt.Printf("[KAFKA ERROR] "+msg, a...); fmt.Println() }),
		OffsetOutOfRangeError: true,
		// Logger:                kafka.LoggerFunc(func(msg string, a ...any) { fmt.Printf("[KAFKA LOG] "+msg, a...); fmt.Println() }),
	})
}

func (c *Consumer) run() {
	reader := c.newReader()
	defer c.wg.Done()
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(c.ctx)
		if errors.Is(err, context.Canceled) {
			fmt.Println("ReadMessage context Canceled.")
			return
		}

		if err != nil {
			fmt.Println("ReadMessage fail : ", err)
			reader.Close()
			reader = c.newReader()
			continue
		}

		if err := c.handler(msg); err != nil {
			fmt.Println("handler fail : ", err)
		}
	}
}
