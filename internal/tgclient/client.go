package tgclient

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(request *Request)

// Client struct instance.
type Client struct {
	tgAPIToken string
	TgClient   *tgbotapi.BotAPI

	commandHandlers        map[string]HandlerFunc
	commandNotFoundHandler HandlerFunc
	textHandler            HandlerFunc

	updateChan chan tgbotapi.Update
	wg         sync.WaitGroup
}

// Creates a new Telegram client.
func New(tgAPIToken string) *Client {
	tgClient, err := tgbotapi.NewBotAPI(tgAPIToken)
	if err != nil {
		log.Panic(err)
	}

	webhook, err := tgClient.GetWebhookInfo()
	if err != nil {
		log.Panic(err)
	}

	if webhook.IsSet() {
		log.Panicf(
			"Conflict: can't use getUpdates method while Webhook (%s) is active; "+
				"use deleteWebhook to delete the webhook and try again.",
			webhook.URL,
		)
	}

	return &Client{
		tgAPIToken:      tgAPIToken,
		TgClient:        tgClient,
		commandHandlers: make(map[string]HandlerFunc),
	}
}

// Add a new command handler.
func (client *Client) AddCommandHandler(cmdName string, handlerFunc HandlerFunc) {
	client.commandHandlers[cmdName] = handlerFunc
}

// Set command not found handler.
func (client *Client) SetCommandNotFoundHandler(handlerFunc HandlerFunc) {
	client.commandNotFoundHandler = handlerFunc
}

// Set text handler.
func (client *Client) SetTextHandler(handlerFunc HandlerFunc) {
	client.textHandler = handlerFunc
}

// Message worker.
func (client *Client) messageWorker(wId int) {
	defer client.wg.Done()
	log.Printf("Worker #%d started.", wId)
	for update := range client.updateChan {
		req := Request{
			WorkerID: wId,
			BotAPI:   client.TgClient,
			Update:   &update,
		}

		if update.Message.IsCommand() {
			cmdHandler, exists := client.commandHandlers[update.Message.Command()]
			if !exists {
				if client.commandNotFoundHandler != nil {
					client.commandNotFoundHandler(&req)
				}
				continue
			}

			cmdHandler(&req)
			continue
		}

		if client.textHandler == nil {
			continue
		}

		client.textHandler(&req)
	}
	log.Printf("Worker #%d ended.", wId)
}

// Start listening for Telegram updates.
func (client *Client) Listen() {
	log.Println("Start listener.")
	client.updateChan = make(chan tgbotapi.Update, UPDATE_CHANNEL_BUFFER)

	// Start the workers.
	for i := 0; i < MAX_MESSAGE_WORKERS; i++ {
		client.wg.Add(1)
		go client.messageWorker(i)
	}

	// Get update channel.
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := client.TgClient.GetUpdatesChan(u)

	for update := range updates {
		// Ignore any non-Message updates.
		if update.Message == nil {
			continue
		}

		// Ingore bot messages.
		if update.Message.From.IsBot {
			continue
		}

		// Push the update into the update channel.
		client.updateChan <- update
	}

	// Close the update channel which will stop the workers.
	close(client.updateChan)

	// Wait for the workers.
	client.wg.Wait()
	log.Println("Listener has been stopped.")
}

// Shutdown the telegram client.
func (client *Client) Shutdown() {
	client.TgClient.StopReceivingUpdates()
}
