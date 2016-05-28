package main

import (
	"flag"
	"fmt"
	"os"
    // "github.com/artsafin/ratebot/alerts"
)

var url = flag.String("url", "localhost:80", "URL that will respond with JSON rates")
var botToken = flag.String("token", "", "bot token")
var debug = flag.Bool("debug", false, "run bot in debug mode")

func main() {
	flag.Parse()

	if *botToken == "" {
		fmt.Println("error: --token must not be empty")
		os.Exit(1)
	}

	bot := NewBotService(*botToken, *debug)

	bot.ConsumeBotCommands(BotHandler)
}
