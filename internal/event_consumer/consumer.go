package event_consumer

import (
	"layouts-orders-bot/internal/entity"
	"log"
	"time"
)

type Processor interface {
	Fetch(limit int) ([]entity.Event, error)
	Process(entity.Event) error
}

type Consumer struct {
	processor Processor
	batchSize int
}

func New(processor Processor, batchSize int) Consumer {
	return Consumer{
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.processor.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		c.handleEvents(gotEvents)
	}
}

func (c *Consumer) handleEvents(events []entity.Event) {
	for _, event := range events {
		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())
			continue
		}
	}
}
