package main

type ChatFlow struct {
    t string
}

var flows = make(map[int]ChatFlow)

func IsNoFlow(chatId int) bool {
    return IsFlow(chatId, "")
}

func IsFlow(chatId int, flow string) bool {
    return flows[chatId].t == flow
}

func SetFlow(chatId int, flow string) {
    flows[chatId] = ChatFlow{flow}
}

func ResetFlow(chatId int) {
    SetFlow(chatId, "")
}