package main

import (
	"log"
	"os"
	"tbot/controller"
	"tbot/handler"

	tgram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botToken, exist := os.LookupEnv("BOT_TOKEN")
	if !exist {
		log.Fatal("$BOT_TOKEN is not set up")
	}
	if botToken == "" {
		log.Fatal("$BOT_TOKEN is empty")
	}

	bot, err := tgram.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	close := make(chan struct{})

	u := tgram.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	up := handler.NewUpdate(controller.NewResponder(bot), close)
	//TODO: fan in/fan out
	for update := range updates {
		//TODO: check
		go up.Handle(update)
	}

	close <- struct{}{}
}
