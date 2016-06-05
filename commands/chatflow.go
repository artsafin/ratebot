package commands

type chatFlow interface {}

type forgetChatFlow struct {
    chatFlow
}

type nilFlow struct {
    chatFlow
}
