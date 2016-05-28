package main

type ChatFlow struct {
    t string
}

var flows = make(map[int64]ChatFlow)

func IsNoFlow(chatId int64) bool {
    return IsFlow(chatId, "")
}

func IsFlow(chatId int64, flow string) bool {
    return flows[chatId].t == flow
}

func SetFlow(chatId int64, flow string) {
    flows[chatId] = ChatFlow{flow}
}

func ResetFlow(chatId int64) {
    SetFlow(chatId, "")
}