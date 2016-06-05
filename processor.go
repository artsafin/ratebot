package main

import (
    "log"
    "fmt"
    "github.com/artsafin/ratebot/botservice"
    "github.com/artsafin/ratebot/alerts"
    "github.com/artsafin/ratebot/event"
)

type FiredAlert struct {
    alert alerts.Alert
    quote event.Quote
}

func NewFiredChannel() chan FiredAlert {
    return make(chan FiredAlert, 20)
}

func ProcessEvents(srcChan <-chan event.ChangedEvent, firedChan chan<- FiredAlert) {
    for ev := range srcChan {
        symAlerts := alerts.GetAlertsBySymbol(ev.Symbol)

        for _, alert := range symAlerts {
            if alert.Test(ev.New.Bid) {
                fmt.Println("fired:", alert, ev.New.Bid)
                firedChan <- FiredAlert{*alert, ev.New}

                fmt.Println("alert1:", alert)
                alert.TryExpire()
                fmt.Println("alert2:", alert)
            }
        }
    }
}

func createAlertMessage(fired FiredAlert) string {
    return fmt.Sprintf("Alert: %s\nBid: <b>%v</b>\nAsk: <b>%v</b>", fired.alert.String(), fired.quote.Bid, fired.quote.Ask)
}

func SendFiredAlerts(bot *botservice.BotSvc, firedChan <-chan FiredAlert) {
    for fired := range firedChan {
        message := bot.CreateMessage(fired.alert.ChatId, createAlertMessage(fired), nil)
        _, err := bot.Send(message)
        if err != nil {
            log.Println(err)
        }
    }
}