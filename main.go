package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/binsabit/tgSaveLinkBot/consumer/event_consumer"
	"github.com/binsabit/tgSaveLinkBot/events/telegram"
	"github.com/binsabit/tgSaveLinkBot/storage/files"
	tgClient "github.com/binsabit/tgSaveLinkBot/telegram"
)

const (
	host = "api.telegram.org"
)

func main() {

	savepath, token := mustPathAndToken()
	tgClient := tgClient.New(host, token)

	s := files.New(savepath)
	every := time.Duration(24 * 60 * time.Minute)
	ticker := time.NewTicker(every)
	go randomSender(ticker, tgClient, savepath, s)
	eventProcessor := telegram.New(tgClient, s)

	consumer := event_consumer.New(eventProcessor, eventProcessor)
	consumer.Start()

}

func mustPathAndToken() (string, string) {
	token := flag.String("tg-bot-token", "", "token to access telegram bot")
	savepath := flag.String("save-path", "", "file path to save links")

	flag.Parse()

	if *token == "" {
		log.Fatal("bot token must be provided")
	}
	if *savepath == "" {
		log.Fatal("file path must be provided")
	}
	_ = os.Mkdir(*savepath, os.ModePerm)

	return *savepath, *token
}

func randomSender(ticker *time.Ticker, tgClient *tgClient.Client, savepath string, s *files.Storage) {
	for range ticker.C {
		fmt.Println("int the ticker")
		chats, err := ioutil.ReadDir(savepath)
		if err != nil {
			continue
		}
		if len(chats) == 0 {
			continue
		}
		for _, each := range chats {
			random := rand.Intn(len(chats))
			if random%3 == 0 {
				continue
			}
			go func(fs.FileInfo) {

				chatId, _ := strconv.Atoi(each.Name())
				link, err := s.GetRandom(chatId)
				if err != nil {
					return
				}
				err = tgClient.SendMessage(link.ChatId, link.Content)
				if err != nil {
					return
				}
			}(each)
		}
	}
}
