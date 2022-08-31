package event_consumer

import (
	"log"

	"github.com/binsabit/tgSaveLinkBot/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
}

func New(fetcher events.Fetcher, processor events.Processor) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
	}
}

func (c Consumer) Start() {
	for {
		events, err := c.fetcher.Fetch()
		if err != nil {
			log.Println(err.Error())
		}

		if err := c.handleEvents(events); err != nil {
			log.Println(err)
			continue
		}

	}
}

func (c Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		if err := c.processor.Process(event); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}
