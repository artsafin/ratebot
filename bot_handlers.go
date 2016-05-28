package main

import (
	"strings"
    "strconv"
	"github.com/go-telegram-bot-api/telegram-bot-api"
    "fmt"
    "github.com/artsafin/ratebot/alerts"
)

func isFirstWord(text string, word string) bool {
    fields := strings.Fields(text)
    if len(fields) > 0 && (fields[0] == word || fields[0] == "/" + word) {
        return true
    }
    return false
}

func BotHandler(text string, uch *UserChat, update *tgbotapi.Update) (*HandlerReply, error) {
    if isFirstWord(text, "cancel") {
        return cancel(uch.chat), nil
    }

    if isFirstWord(text, "help") {
        return &HandlerReply{text: help(0)}, nil
    }

    if isFirstWord(text, "instr") {
        return &HandlerReply{text: instr()}, nil
    }

    if isFirstWord(text, "alerts") {
        return alertList(uch.chat)
    }

    if isFirstWord(text, "forget") || IsFlow(uch.chat, "forget") {
        return forget(uch.chat, text)
    }

    if alertParams, err := parseAlertInput(text); err == nil {
        return alert(uch.chat, alertParams)
    } else {
        fmt.Println("parseAlertInput:", err)
    }

	return nil, BotErr("command not found")
}

func help(helpSection byte) string {
    head := `
Bot is able to notify you on currency rate changes.
All commands are available both with or without leading slash.

/help (or help) - Show this message
/instr (or instr) - Show supported instruments
/alerts (or alerts) - Show currently set alerts
/forget (or forget) - Unsubscribe from alert

`
    newa := `How to set new alerts:
    <code>/alert INSTRUMENT OPERATION VALUE</code>
    <code>alert INSTRUMENT OPERATION VALUE</code>
or just:
    <code>INSTRUMENT OPERATION VALUE</code>

where:
<code>INSTRUMENT</code> is one of the supported instruments (see /instr), case insensitive
<code>OPERATION</code> is one of:
    = (or eq, equals)
    &lt; (or lt, less than)
    &lt;= (or lte, less than or equals)
    &gt; (or gt, greater than)
    &gt;= (or gte, greater than or equals)
<code>VALUE</code> is a decimal number with dot as a decimal part separator

`
    if helpSection == 0 {
        return head + newa
    }

    var res string
    if helpSection & 1 == 1 {
        res += head
    }
    if helpSection & 2 == 2 {
        res += newa
    }
    return res
}

func instr() string {
    return "Supported instruments:\n\n" + strings.Join(alerts.GetInstr(), "\n")
}

func cancel(chatId int64) *HandlerReply {
    flows[chatId] = ChatFlow{""}

    return &HandlerReply{
                    text: "Cancelled current operation",
                    configurator: func(cfg *tgbotapi.MessageConfig) {
                        cfg.ReplyMarkup = tgbotapi.ReplyKeyboardHide {
                            HideKeyboard: true,
                        }
                    },
                }
}

func forget(chatId int64, text string) (*HandlerReply, error) {
    userAlerts, err := alerts.GetAlertsByChatId(chatId)
    if err != nil {
        return nil, BotPublicErr("No registered alerts")
    }

    kbConfig := func(userAlerts []*alerts.Alert, err error) func(*tgbotapi.MessageConfig) {
        return func(cfg *tgbotapi.MessageConfig) {
            if err != nil {
                cfg.ReplyMarkup = tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
            } else {
                btnsL := make([]tgbotapi.KeyboardButton, len(userAlerts)/2)
                btnsR := make([]tgbotapi.KeyboardButton, len(userAlerts)/2+1)
                for k, s := range userAlerts {
                    if k <= len(userAlerts)/2 {
                        btnsL[k] = tgbotapi.KeyboardButton{Text:s.String()}
                    } else {
                        btnsR[k] = tgbotapi.KeyboardButton{Text:s.String()}
                    }
                }

                cfg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
                    Keyboard: [][]tgbotapi.KeyboardButton{
                        btnsL,
                        btnsR,
                    },
                }
            }
        }
    }

    if IsNoFlow(chatId) {
        SetFlow(chatId, "forget")

        return &HandlerReply{
                text: "What alert would you like to unsubscribe from? Type /cancel to quit",
                configurator: kbConfig(userAlerts, err),
            }, nil
    } else if IsFlow(chatId, "forget") {
        for _, v := range userAlerts {
            if v.String() == text {
                alerts.RemoveAlert(*v)
                userAlerts, err = alerts.GetAlertsByChatId(chatId)
                if err != nil {
                    ResetFlow(chatId)
                }

                return &HandlerReply{
                    text: "Removed " + v.String(),
                    configurator: kbConfig(userAlerts, err),
                }, nil
            }
        }
    }

    return nil, nil
}

func alertList(chatId int64) (*HandlerReply, error) {
    userAlerts, err := alerts.GetAlertsByChatId(chatId)
    if err != nil {
        return nil, BotPublicErr("No registered alerts")
    }

    msg := strings.Join(alerts.AlertsArrayToStrings(userAlerts), "\n")
    msg = strings.Replace(msg, "<", "&lt;", -1)
    msg = strings.Replace(msg, ">", "&gt;", -1)

    return &HandlerReply{text: msg}, nil
}

func parseAlertInput(text string) ([]string, error) {
    text = strings.TrimSpace(strings.Replace(strings.Replace(text, "alert ", "", 1), "/alert ", "", 1))

    result := make([]string, 3)

    for _, v := range alerts.GetInstr() {
        if len(text) >= len(v) && strings.ToLower(text[:len(v)]) == strings.ToLower(v) {
            result[0] = v
            text = strings.TrimSpace(text[len(v):])
            break
        }
    }

    if result[0] == "" {
        return result, fmt.Errorf("Instrument not found")
    }

    for _, v := range alerts.GetOps() {
        if len(text) >= len(v) && strings.ToLower(text[:len(v)]) == strings.ToLower(v) {
            result[1] = v
            text = strings.TrimSpace(text[len(v):])
            break
        }
    }

    if result[1] == "" {
        return result, fmt.Errorf("Operation not found")
    }

    if text == "" {
        return result, fmt.Errorf("Value not found")
    }

    result[2] = text

    return result, nil
}

func alert(chatId int64, fields []string) (*HandlerReply, error) {

    fmt.Println("alert:", fields)

    if len(fields) != 3 {
        return nil, BotPublicErr("<em>Invalid format for alert command</em>\n\n" + help(2))
    }

    opValue, err := strconv.ParseFloat(fields[2], 32)

    if err != nil {
        return nil, BotPublicErr("<em>Unable to parse VALUE</em>\n\n" + help(2))
    }

    alert, err := alerts.NewAlert(chatId, fields[0], fields[1], float32(opValue))
    if err != nil {
        return nil, BotPublicErr("<em>" + err.Error() + "</em>")
    }

    alerts.AddAlert(alert)

    return &HandlerReply{text: "Registered"}, nil
}