package tgclient

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(request *Request)

type Client struct {
	tgAPIToken string
	tgClient   *tgbotapi.BotAPI

	commandHandlers        map[string]HandlerFunc
	commandNotFoundHandler HandlerFunc
	textHandler            HandlerFunc

	updateChan chan tgbotapi.Update
	wg         sync.WaitGroup
}

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
		tgClient:        tgClient,
		commandHandlers: make(map[string]HandlerFunc),
	}
}

func (client *Client) AddCommandHandler(cmdName string, handlerFunc HandlerFunc) {
	client.commandHandlers[cmdName] = handlerFunc
}

func (client *Client) AddCommandNotFoundHandler(handlerFunc HandlerFunc) {
	client.commandNotFoundHandler = handlerFunc
}

func (client *Client) AddTextHandler(handlerFunc HandlerFunc) {
	client.textHandler = handlerFunc
}

func (client *Client) handlerWorker(wId int) {
	defer client.wg.Done()
	log.Printf("Worker #%d started", wId)
	for update := range client.updateChan {
		req := Request{
			WorkerID: wId,
			BotAPI:   client.tgClient,
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
	log.Printf("Worker #%d ended", wId)
}

func (client *Client) Listen() {
	log.Println("Start listener")
	client.updateChan = make(chan tgbotapi.Update, UPDATE_CHANNEL_BUFFER)

	// start the workers
	for i := 0; i < MAX_HANDLER_WORKERS; i++ {
		client.wg.Add(1)
		go client.handlerWorker(i)
	}

	// get update channel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := client.tgClient.GetUpdatesChan(u)

	for update := range updates {
		// ignore any non-Message updates
		if update.Message == nil {
			continue
		}

		// ingore bot messages
		if update.Message.From.IsBot {
			continue
		}

		// push the update into the update channel
		client.updateChan <- update
	}

	// close the update channel which will stop the workers
	close(client.updateChan)

	// wait for the workers
	client.wg.Wait()
	log.Println("Listener has been stopped")
}

func (client *Client) Shutdown() {
	client.tgClient.StopReceivingUpdates()
}
