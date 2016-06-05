package commands

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/artsafin/ratebot/botservice"
)

func cancel(flow chatFlow, chatId int64) (chatFlow, *botservice.HandlerReply) {
    return nil, &botservice.HandlerReply{
                    Text: "Cancelled current operation",
                    Configurator: func(cfg *tgbotapi.MessageConfig) {
                        cfg.ReplyMarkup = tgbotapi.ReplyKeyboardHide {
                            HideKeyboard: true,
                        }
                    },
                }
}