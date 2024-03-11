package producer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

type Producer struct {
	Bootstrap string
	Topic     string
	In        chan kafka.Message
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewProducer(Bootstrap, Topic string) (*Producer, error) {
	if len(Bootstrap) <= 0 {
		return nil, fmt.Errorf("invalid Bootstrap (%s)", Bootstrap)
	}

	if len(Topic) <= 0 {
		return nil, fmt.Errorf("invalid Topic (%s)", Topic)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	return &Producer{
		Bootstrap: Bootstrap,
		Topic:     Topic,
		In:        make(chan kafka.Message),
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (p *Producer) Open() error {
	p.wg.Add(1)
	go p.run()

	return nil
}

func (p *Producer) Close() {
	p.cancel()
	p.wg.Wait()
}

func (p *Producer) newWriter() *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          strings.Split(p.Bootstrap, ","),
		Topic:            p.Topic,
		BatchSize:        100,
		BatchBytes:       1000 * 1000 * 100,
		BatchTimeout:     1 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
		ErrorLogger:      kafka.LoggerFunc(func(msg string, a ...any) { fmt.Printf("[KAFKA] "+msg, a...); fmt.Println() }),
		Async:            true,
		ReadTimeout:      1 * time.Second,
	})
	writer.AllowAutoTopicCreation = true

	return writer
}

func (p *Producer) run() {
	writer := p.newWriter()
	defer p.wg.Done()
	defer writer.Close()
	defer close(p.In)

	for {
		select {
		case <-p.ctx.Done():
			return
		case msg := <-p.In:
			err := writer.WriteMessages(p.ctx, msg)
			if err != nil {
				fmt.Println("WriteMessage fail : ", err)
				writer.Close()
				writer = p.newWriter()
				continue
			}
		}
	}
}
