package main

import (
	"flag"
	"fmt"
	"os"
    // "github.com/artsafin/ratebot/alerts"
)

var srcAddr = flag.String("addr", "localhost:8080", "source service address")
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
