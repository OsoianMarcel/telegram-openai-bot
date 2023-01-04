package tgclient

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Request struct {
	Update   *tgbotapi.Update
	WorkerID int
	BotAPI   *tgbotapi.BotAPI
}

func (req *Request) Reply(message string) {
	msg := tgbotapi.NewMessage(req.Update.Message.Chat.ID, message)
	msg.ReplyToMessageID = req.Update.Message.MessageID

	if _, err := req.BotAPI.Send(msg); err != nil {
		log.Println(err)
	}
}

func (req *Request) Send(message string) {
	msg := tgbotapi.NewMessage(req.Update.Message.Chat.ID, message)

	if _, err := req.BotAPI.Send(msg); err != nil {
		log.Println(err)
	}
}

func (req *Request) SendAction(name string) {
	chatAction := tgbotapi.NewChatAction(req.Update.Message.Chat.ID, name)

	if _, err := req.BotAPI.Request(chatAction); err != nil {
		log.Println(err)
	}
}
