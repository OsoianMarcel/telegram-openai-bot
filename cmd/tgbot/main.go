package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/OsoianMarcel/tg-bot/internal/gptclient"
	"github.com/OsoianMarcel/tg-bot/internal/stats"
	"github.com/OsoianMarcel/tg-bot/internal/tgclient"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx := context.Background()

	st := stats.New(os.Getenv("STATS_FILE"))

	if st.IsFileSet() {
		log.Printf("Load the statistics from the file (%s)...\n", st.GetFilePath())
		if err := st.LoadFromFile(); err != nil {
			log.Panicln(err)
		}
	}

	tc := tgclient.New(os.Getenv("TG_APITOKEN"))
	gptc := gptclient.New(os.Getenv("GPT_AUTH_TOKEN"))

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
		stats, err := st.Stats.GenString()

		if err != nil {
			log.Println(err)
			req.Reply("Something went wrong. Can not generate the stats.")
			return
		}

		req.Reply(stats)
	})

	tc.SetCommandNotFoundHandler(func(req *tgclient.Request) {
		req.Reply("Command not found.")
	})

	tc.SetTextHandler(func(req *tgclient.Request) {
		text := req.Update.Message.Text
		from := req.Update.Message.From
		userId := from.ID

		// Validate the message.
		if len(text) < 2 {
			st.Stats.IncrAiInvalidErrors(userId)
			req.Reply("The message is too short.")
			return
		}

		if len(text) > 1024 {
			st.Stats.IncrAiInvalidErrors(userId)
			req.Reply("The message is too long (max: 1024 characters).")
			return
		}

		// Send typing action.
		req.SendAction("typing")

		// Ask the OpenAI.
		aiUserId := strconv.FormatInt(userId, 10)
		res, err := gptc.AskAI(ctx, text, aiUserId)
		if err != nil || len(res) == 0 {
			if err != nil {
				log.Println(err)
			}

			if errors.Is(err, context.DeadlineExceeded) {
				req.Reply("Error: AI request timeout occurred...\nPlease, try again.")
				st.Stats.IncrAiTimeoutErrors(userId)
				return
			}

			req.Reply("Error: The AI is unavailable or has no response...\nPlease, try again.")
			st.Stats.IncrAiErrors(userId)
			return
		}

		req.Reply(res)
		st.Stats.IncrAiResponses(userId)
	})

	var wg sync.WaitGroup

	// Listen for stop signal.
	wg.Add(1)
	go func() {
		defer wg.Done()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit
		log.Println("Graceful shutdown in progress...")

		log.Println("Shutdown the telegram client...")
		tc.Shutdown()

		log.Println("Exit the graceful shutdown goroutine.")
	}()

	// The main thread is blocked by the method Listen until the client is not closed.
	tc.Listen()
	log.Println("The telegram client has been stopped.")

	if st.IsFileSet() {
		log.Printf("Write the statistics to the file (%s)...\n", st.GetFilePath())
		if err := st.WriteToFile(); err != nil {
			log.Println(err)
		}
	}

	log.Println("Exit.")
}
