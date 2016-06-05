package main

import (
	"flag"
	"fmt"
	"log"
	"os"
    "time"
    "github.com/artsafin/ratebot/botservice"
    "github.com/artsafin/ratebot/commands"
    "github.com/artsafin/ratebot/event"
    "github.com/schachmat/ingo"
)

var url = flag.String("url", "http://localhost", "URL that will respond with JSON rates. Format: [{\"symbol\": string, \"ask\": string, \"bid\": string}, ...]")
var tickPeriod = flag.Duration("tick", time.Second * 2, "How often to check --url for updates")
var botToken = flag.String("token", "", "Bot token provided by Telegram")
var debug = flag.Bool("debug", false, "Debug HTTP requests made to Telegram API")

func main() {
	if err := ingo.Parse("ratebot"); err != nil {
        log.Fatal(err)
    }

	src, err := event.NewSource(*url)
	if err != nil {
		log.Fatal(err)
	}

	if *botToken == "" {
		fmt.Println("error: --token must not be empty")
		os.Exit(1)
	}

	fmt.Printf("Ratebot v.1\nListening %s every %v\n", *url, *tickPeriod)

	bot, err := botservice.NewBotService(*botToken, *debug)
	if err != nil {
		log.Fatal(err)
	}

	botChan, err := bot.NewUpdateChannel(2)
	if err != nil {
		log.Fatal(err)
	}

	handlerCb := commands.NewCommandsHandler(src)
	go bot.ConsumeBotCommands(handlerCb, botChan)

	srcChan := event.NewListenChannel()
	go event.ListenSource(*tickPeriod, src, srcChan)

	
	firedChan := NewFiredChannel()
	go SendFiredAlerts(bot, firedChan)

	ProcessEvents(srcChan, firedChan)
}
