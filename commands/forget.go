package commands

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/artsafin/ratebot/botservice"
    "github.com/artsafin/ratebot/alerts"
)

func alertsKeyboard(alerts []*alerts.Alert) interface{} {
    if len(alerts) > 0 {
        btns := make([][]tgbotapi.KeyboardButton, 1)

        for k, al := range alerts {
            if k != 0 && k % 3 == 0 {
                btns = append(btns, make([]tgbotapi.KeyboardButton, 1))
            }
            btns[len(btns)-1] = append(btns[len(btns)-1], tgbotapi.NewKeyboardButton(al.StringShort()))
        }

        return tgbotapi.ReplyKeyboardMarkup{
            Keyboard: btns,
        }
    } else {
        return hideKeyboard()
    }
}

func hideKeyboard() tgbotapi.ReplyKeyboardHide {
    return tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
}

func forget(flow *forgetChatFlow, userId int, text string) (chatFlow, *botservice.HandlerReply, error) {
    
    if flow != nil {
        foundAlert := alerts.FindUserAlert(userId, func(a *alerts.Alert) bool { return a.StringShort() == text })
        if foundAlert == nil {
            return nil, nil, botservice.BotPublicErr("Alert not found, try repeating operation again")
        }

        alerts.RemoveAlert(foundAlert)

        userAlerts := alerts.GetAlertsByUserId(userId)

        var retFlow chatFlow
        if len(userAlerts) == 0 {
            retFlow = nil
        } else {
            retFlow = flow
        }

        return retFlow, &botservice.HandlerReply{
                Text: "Removed " + foundAlert.StringShort(),
                Configurator: func(cfg *tgbotapi.MessageConfig) {
                    cfg.ReplyMarkup = alertsKeyboard(userAlerts)
                },
            }, nil
    } else {
        userAlerts := alerts.GetAlertsByUserId(userId)
        if len(userAlerts) == 0 {
            return nil, nil, botservice.BotPublicErr("You don't have registered alerts for removal")
        }

        retFlow := &forgetChatFlow{}

        return retFlow, &botservice.HandlerReply{
                Text: "What alert would you like to unsubscribe from? Type /cancel to quit",
                Configurator: func(cfg *tgbotapi.MessageConfig) {
                    cfg.ReplyMarkup = alertsKeyboard(userAlerts)
                },
            }, nil
    }
}
