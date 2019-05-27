package main

import (
	"log"
	"os"
	"strings"
	"time"

	nsqio "github.com/nsqio/go-nsq"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func PublishMessage(command, topicName string) {
	config := nsqio.NewConfig()
	w, _ := nsqio.NewProducer("flickpi.home:4150", config)
	//w.SetLogger(nsqio.nullLogger, nsqio.LogLevelInfo)
	defer w.Stop()

	log.Printf("Publishing topic '%s' command '%s'", topicName, command)

	err := w.Publish(topicName, []byte("play"))
	if err != nil {
		log.Fatalf("error %s", err)
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("FOO"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if strings.Contains(update.Message.Text, "osalie") {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			PublishMessage("play", "rosalie")
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "use rosalie")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

	}
}
