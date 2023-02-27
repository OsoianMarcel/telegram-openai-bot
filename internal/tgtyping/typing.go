package tgtyping

import (
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChatData struct {
	workers      uint16
	lastTypingOn time.Time
}

type TgTyping struct {
	botAPI        *tgbotapi.BotAPI
	chats         map[int64]ChatData
	mutex         sync.RWMutex
	ticker        *time.Ticker
	stopTicker    chan struct{}
	waitGroup     sync.WaitGroup
	tickerWorking bool
}

func New(botAPI *tgbotapi.BotAPI) *TgTyping {
	return &TgTyping{
		botAPI: botAPI,
		chats:  make(map[int64]ChatData),
	}
}

func (t *TgTyping) StarTicker() {
	if t.tickerWorking {
		return
	}
	t.tickerWorking = true
	t.ticker = time.NewTicker(time.Second)
	t.stopTicker = make(chan struct{})

	t.waitGroup.Add(1)
	go func() {
		defer t.waitGroup.Done()
		for {
			select {
			case <-t.stopTicker:
				return
			case <-t.ticker.C:
				t.tick()
			}
		}
	}()
}

func (t *TgTyping) Shutdown() {
	if !t.tickerWorking {
		return
	}

	t.ticker.Stop()
	t.stopTicker <- struct{}{}
	t.waitGroup.Wait()
}

func (t *TgTyping) tick() {
	t.mutex.RLock()
	chatIDs := make([]int64, 0, len(t.chats))
	for chatID := range t.chats {
		chatIDs = append(chatIDs, chatID)
	}
	t.mutex.RUnlock()

	for _, chatID := range chatIDs {
		t.sendTyping(chatID, false)
	}
}

func (t *TgTyping) sendTyping(chatID int64, forceSend bool) {
	t.mutex.Lock()
	chat, ok := t.chats[chatID]
	if !ok {
		chat = ChatData{}
	}
	if time.Since(chat.lastTypingOn).Seconds() <= 3 && !forceSend {
		t.mutex.Unlock()
		return
	}
	chat.lastTypingOn = time.Now()
	t.chats[chatID] = chat
	t.mutex.Unlock()

	chatAction := tgbotapi.NewChatAction(chatID, "typing")
	if _, err := t.botAPI.Request(chatAction); err != nil {
		log.Println(err)
	}
}

func (t *TgTyping) Typing(chatID int64) func() {
	// Add chat id to chats.
	t.mutex.Lock()
	chat, ok := t.chats[chatID]
	if !ok {
		chat = ChatData{}
	}
	chat.workers++
	t.chats[chatID] = chat
	t.mutex.Unlock()

	// Send typing immediately.
	t.sendTyping(chatID, false)

	var done bool

	return func() {
		// If done was already called then stop here.
		if done {
			return
		}
		done = true

		t.mutex.Lock()
		chat, ok := t.chats[chatID]
		// If no chats, stop here.
		if !ok {
			t.mutex.Unlock()
			return
		}
		// If number of chat workers is just one, then delete the chat.
		if chat.workers <= 1 {
			delete(t.chats, chatID)
			t.mutex.Unlock()
			return
		}
		// Otherwise decrease the number of active workers.
		chat.workers--
		t.chats[chatID] = chat
		t.mutex.Unlock()
	}
}
