package telegram

import (
	"errors"

	"github.com/binsabit/tgSaveLinkBot/events"
	"github.com/binsabit/tgSaveLinkBot/storage"
	"github.com/binsabit/tgSaveLinkBot/telegram"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatId int
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch() ([]events.Event, error) {
	updates, err := p.tg.GetUpdates(p.offset)
	if err != nil {
		return nil, err
	}
	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, toEvent(u))
	}

	p.offset = updates[len(updates)-1].UpdateId + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errors.New("Unknown message type")
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := metaOf(event)
	if err != nil {
		return err
	}
	if err := p.doCmd(event.Text, meta.ChatId); err != nil {
		return errors.New("can't process message")
	}

	return nil
}

func metaOf(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errors.New("can't get Meta")
	}
	return res, nil
}
func toEvent(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatId: upd.Message.Chat.ChatId,
		}
	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == (telegram.Message{}) {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == (telegram.Message{}) {
		return events.Unknown
	}
	return events.Message
}
