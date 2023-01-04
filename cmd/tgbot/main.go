package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/OsoianMarcel/tg-bot/internal/gptclient"
	"github.com/OsoianMarcel/tg-bot/internal/stats"
	"github.com/OsoianMarcel/tg-bot/internal/tgclient"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx := context.Background()

	st := stats.NewStats()

	tc := tgclient.NewClient(os.Getenv("TG_APITOKEN"))
	gptc := gptclient.NewClient(os.Getenv("GPT_AUTH_TOKEN"))

	tc.AddCommandHandler("start", func(req *tgclient.Request) {
		req.Send(
			"Welcome! | Bine ai venit! | ¡Bienvenidos! | Добро пожаловать!\n\n" +
				"- Ask me a question!\n" +
				"- Întreabă-mă ceva!\n" +
				"- Hazme una pregunta!\n" +
				"- Задайте мне вопрос!\n",
		)
	})

	tc.AddCommandHandler("help", func(req *tgclient.Request) {
		req.Send(
			"This is a bot mode for communication with artificial intelligence (OpenAI).\n\n" +

				"How does it work?\n" +
				"Just send a message in chat, and the message will be forwarded to the AI, " +
				"and soon you'll get the reply.\n" +
				"You can ask him questions, or even ask him to help you with something.\n\n" +

				"What languages does the AI understand?\n" +
				"The primary language is English, but the AI understand other languages too" +
				", give it a try.\n\n" +

				"How good is this AI?\n" +
				"The AI is still in its early stages, however, and it has a long way to go " +
				"before it can match the intelligence of a human.\n\n" +

				"How can I contact you?\n" +
				"Type the command /feedback fallowed by your message, and the bot will forward " +
				"your message to the owner of this bot.\n\n" +

				"Do you collect the messages?\n" +
				"No, I do not collect your messages, the bot just forward your message to the AI.\n" +
				"The bot may collect some statistic information about usage (ex.: number of messages, " +
				"number of words etc.), but the message itself is not saved anywhere.",
		)
	})

	tc.AddCommandHandler("feedback", func(req *tgclient.Request) {
		adminChatIDStr, exists := os.LookupEnv("TG_ADMIN_CHATID")
		if !exists {
			req.Reply("This feature is currently unavailable.")
			return
		}

		adminChatID, err := strconv.ParseInt(adminChatIDStr, 10, 0)
		if err != nil {
			log.Println(err)
			return
		}

		adminMessage := req.Update.Message.CommandArguments()
		if len(adminMessage) < 16 {
			req.Reply("Your message is too short. Write minimum 16 characters.")
			return
		}

		msg := tgbotapi.NewForward(
			adminChatID,
			req.Update.Message.From.ID,
			req.Update.Message.MessageID,
		)
		if _, err := req.BotAPI.Send(msg); err != nil {
			log.Println(err)
		} else {
			req.Reply(
				"Your message has been delivered.\n" +
					"Thank you for using our service.")
		}
	})

	tc.AddCommandHandler("me", func(req *tgclient.Request) {
		reply := fmt.Sprintf(
			"Username: %s\n"+
				"User ID: %d",
			req.Update.Message.From.UserName,
			req.Update.Message.From.ID,
		)
		req.Reply(reply)
	})

	tc.AddCommandHandler("who", func(req *tgclient.Request) {
		reply := fmt.Sprintf("I am worker #%d", req.WorkerID)
		req.Reply(reply)
	})

	tc.AddCommandHandler("stats", func(req *tgclient.Request) {
		reply := fmt.Sprintf(
			"All messages: %d\n"+
				"Invalid messages: %d\n"+
				"AI errors: %d\n"+
				"Request chars: %d\n"+
				"Response chars: %d\n",
			st.GetAiAllMessages(),
			st.GetAiInvalidMessages(),
			st.GetAiErrors(),
			st.GetAiRequestChars(),
			st.GetAiResponseChars(),
		)
		req.Reply(reply)
	})

	tc.AddCommandNotFoundHandler(func(req *tgclient.Request) {
		req.Reply("Command not found.")
	})

	tc.AddTextHandler(func(req *tgclient.Request) {
		st.IncAiAllMessages()

		text := req.Update.Message.Text
		username := req.Update.Message.From.UserName

		// validate the message
		if len(text) < 2 {
			st.IncAiInvalidMessages()
			req.Reply("The message is too short.")
			return
		}

		if len(text) > 1024 {
			st.IncAiInvalidMessages()
			req.Reply("The message is too long (max: 1024 characters).")
			return
		}

		// send typing action
		req.SendAction("typing")

		// ask the OpenAI
		res, err := gptc.AskAI(ctx, text, username)
		if err != nil || len(res) == 0 {
			st.IncAiErrors()

			req.Reply("Try again...")
			if err != nil {
				log.Println(err)
			}
			return
		}

		st.AddAiRequestChars(uint32(len(text)))
		st.AddAiResponseChars(uint32(len(res)))

		// reply
		req.Reply(res)
	})

	// listen for stop signal
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit
		log.Println("Graceful shutdown in progress...")
		tc.Shutdown()
	}()

	tc.Listen()

	log.Println("Exit.")
}
