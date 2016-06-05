package commands

import (
    "strings"
    "github.com/artsafin/ratebot/botservice"
    "github.com/artsafin/ratebot/alerts"
)

func alertlist(userId int) (*botservice.HandlerReply, error) {
    userAlerts := alerts.GetAlertsByUserId(userId)
    if len(userAlerts) == 0 {
        return nil, botservice.BotPublicErr("You don't have registered alerts")
    }

    msg := strings.Join(alerts.AlertsArrayToStrings(userAlerts), "\n")

    return &botservice.HandlerReply{Text: msg}, nil
}
