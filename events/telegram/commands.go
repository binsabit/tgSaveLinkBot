package telegram

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/binsabit/tgSaveLinkBot/storage"
)

const (
	RandomCmd = "/random"
	HelpCmd   = "/help"
	StartCmd  = "/start"
)

func (p *Processor) doCmd(text string, chatID int) error {
	text = strings.TrimSpace(text)

	if isAddCmd(text) {
		return p.saveLink(chatID, text)
	}

	switch text {
	case RandomCmd:
		return p.sendRandom(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}
func (p *Processor) saveLink(ChatId int, text string) error {

	link := storage.Link{
		ChatId:  ChatId,
		Content: text,
	}
	isExists := p.storage.IsExists(link)
	if isExists {
		p.tg.SendMessage(ChatId, existsMsg)
	}
	if err := p.storage.Insert(link); err != nil {
		return err
	}

	p.tg.SendMessage(ChatId, savedMsg)

	return nil
}

func (p *Processor) sendRandom(chatId int) error {
	link, err := p.storage.GetRandom(chatId)
	if err != nil {
		fmt.Println("here")
		if errors.Is(err, storage.ErrNoLinkFound) {
			return p.tg.SendMessage(chatId, NoLinksMsg)

		}
		return err
	}

	if err := p.tg.SendMessage(chatId, link.Content); err != nil {
		return err
	}

	return nil
}
func (p *Processor) sendHelp(chatId int) error {
	return p.tg.SendMessage(chatId, helloMsg)
}

func (p *Processor) sendHello(chatId int) error {
	return p.tg.SendMessage(chatId, helloMsg)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
