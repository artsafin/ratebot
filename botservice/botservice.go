package botservice

import (
	"fmt"
    "html"
    "strings"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotInput struct {
    Text string
    Chat int64
    User int
    Update *tgbotapi.Update
}

type BotSvc struct {
	bot *tgbotapi.BotAPI
}

type HandlerReply struct {
	Text         string
	Configurator func(*tgbotapi.MessageConfig)
}

type HandlerCb func(*BotInput) (*HandlerReply, error)

func NewBotService(token string, isDebug bool) (*BotSvc, error) {
	baseBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	baseBot.Debug = isDebug

	return &BotSvc{baseBot}, nil
}

type UserChat struct {
    user int
    chat int64
	// user *tgbotapi.User
	// chat *tgbotapi.Chat
}

func (svc *BotSvc) NewUpdateChannel(timeout int) (<-chan tgbotapi.Update, error) {
    u := tgbotapi.NewUpdate(0)
    u.Timeout = timeout

    return svc.bot.GetUpdatesChan(u)
}

func getMessage(update *tgbotapi.Update) *tgbotapi.Message {
    if update.Message != nil {
        return update.Message
    }
    if update.EditedMessage != nil {
        return update.EditedMessage
    }

    return nil
}

func (svc *BotSvc) ConsumeBotCommands(handler HandlerCb, updatesChan <-chan tgbotapi.Update) {
	for update := range updatesChan {
		svc.processCommand(&update, handler)
	}
}

func (svc *BotSvc) processCommand(update *tgbotapi.Update, handler HandlerCb) {
    message := getMessage(update)
    if message == nil {
        return
    }

	fmt.Printf("[%s] %s\n", message.From.UserName, message.Text)

    input := BotInput{message.Text, message.Chat.ID, message.From.ID, update}

	reply, err := handler(&input)

	if err == nil {
        if reply != nil {
		  svc.bot.Send(svc.CreateMessage(message.Chat.ID, reply.Text, reply.Configurator))
        }
		return
	}
    if _, ok := err.(BotPublicErr); ok {
        svc.bot.Send(svc.CreateMessage(message.Chat.ID, err.Error(), nil))
        return
    }

    if len(message.Text) >= 5 && message.Text[:5] == "/ping" {
        svc.pingReply(message)
    }

	svc.defaultReply(message)
}

func (svc *BotSvc) pingReply(inputMsg *tgbotapi.Message) {
    msg := fmt.Sprintf("I am @%v, %v #%v", svc.bot.Self.UserName, svc.bot.Self.FirstName, svc.bot.Self.ID)
    svc.bot.Send(svc.CreateMessage(inputMsg.Chat.ID, msg, nil))
}

func (svc *BotSvc) defaultReply(inputMsg *tgbotapi.Message) {
    svc.bot.Send(svc.CreateMessage(inputMsg.Chat.ID, "Sorry?", nil))
}

func (svc *BotSvc) CreateMessage(chatId int64, text string, configurator func(*tgbotapi.MessageConfig)) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "HTML"
	// msg.ReplyToMessageID = update.Message.MessageID

    if msg.ParseMode == "HTML" {
        msg.Text = escapeWhitelist(msg.Text)
        // msg.Text = strings.Replace(msg.Text, "&", "&amp;", -1)
     }

	if configurator != nil {
		configurator(&msg)
	}

	return msg
}

func escapeWhitelist(s string) string {
    s = html.EscapeString(s)

    whiteList := []string{"b", "strong", "em", "i", "a", "code", "pre"}

    for _, tag := range whiteList {
        s = strings.Replace(s, "&lt;" + tag + "&gt;", "<" + tag + ">", -1)
        s = strings.Replace(s, "&lt;" + tag + " ", "<" + tag + " ", -1)
        s = strings.Replace(s, "&lt;/" + tag + "&gt;", "</" + tag + ">", -1)
    }

    return s
}

func (svc *BotSvc) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
    return svc.bot.Send(c)
}