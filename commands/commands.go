package commands

import (
    "strings"
    "fmt"
    "github.com/artsafin/ratebot/botservice"
    "github.com/artsafin/ratebot/event"
)

func isFirstWord(text string, word string) bool {
    fields := strings.Fields(text)
    if len(fields) > 0 && (fields[0] == word || fields[0] == "/" + word) {
        return true
    }
    return false
}

func NewCommandsHandler(src *event.Source) botservice.HandlerCb {
    var flows = make(map[int64]chatFlow, 5)

    return func (in *botservice.BotInput) (*botservice.HandlerReply, error) {
        if isFirstWord(in.Text, "cancel") {
            flow, reply := cancel(flows[in.Chat], in.Chat)
            flows[in.Chat] = flow
            return reply, nil
        }

        if isFirstWord(in.Text, "help") {
            return &botservice.HandlerReply{Text: help(0)}, nil
        }

        if isFirstWord(in.Text, "instruments") || isFirstWord(in.Text, "symbols") {
            return &botservice.HandlerReply{Text: instruments(src.KnownSymbols)}, nil
        }

        if isFirstWord(in.Text, "alerts") {
            return alertlist(in.User)
        }

        if flow, isForgetFlow := flows[in.Chat].(*forgetChatFlow); isFirstWord(in.Text, "forget") || isForgetFlow {
            flow, reply, err := forget(flow, in.User, in.Text)
            flows[in.Chat] = flow
            return reply, err
        }

        if alertParams, err := parseAlertInput(in.Text, src.KnownSymbols); err == nil {
            return alert(in.Chat, in.User, alertParams)
        } else {
            fmt.Println("parseAlertInput:", err)
            return nil, err
        }

        return nil, botservice.BotErr("command not found")
    }
}
