package bot

import (
	"context"
	"sync"
	"theater/client"
	"theater/config"
	gomastodon "theater/go-mastodon"
	"theater/log"
)

// NotificationHandler define different behaviors for actors
type NotificationHandler func(*Actor, *gomastodon.Notification, interface{})

// Actor the bot for acting
type Actor struct {
	Name        string
	LineCh      chan string // receive the sentences from script
	BlockCh     chan string // receive the block user id from notification
	UnBlockCh   chan string
	NtfHandler  []NotificationHandler
	client      *client.Bot
	isListening bool
}

// New a actor bot
func New(name string, handlers ...NotificationHandler) *Actor {
	cfg, err := config.ActorBotClientInfo(name)
	if err != nil {
		panic(err)
	}
	c, err := client.New(&cfg)
	if err != nil {
		panic(err)
	}

	return &Actor{
		Name:        name,
		LineCh:      make(chan string),
		BlockCh:     make(chan string),
		UnBlockCh:   make(chan string),
		NtfHandler:  handlers,
		client:      c,
		isListening: false,
	}
}

// Act post the toot according to the script in turn
func (a *Actor) Act(wg *sync.WaitGroup) {
	defer wg.Done()
	isContine := true

	for isContine {
		select {
		case line, ok := <-a.LineCh:
			if !ok {
				isContine = false
				break
			}
			_, err := a.client.PostSpoiler(line, config.TheaterName())
			if err != nil {
				log.SLogger.Errorf("%s post line [%s] to mastodon error: %v", a.Name, line, err)
			}
		case accountID := <-a.BlockCh:
			a.client.BlockAccount(accountID)
			log.SLogger.Infof("block user %s ok", accountID)
		case accountID := <-a.UnBlockCh:
			a.client.UnBlockAccount(accountID)
		}
	}
}

// ListenAudiences some actors need to interact with audiences
// it will execute the handlers owned by actors for every mention notification
func (a *Actor) ListenAudiences(actors map[string]*Actor) {
	if a.isListening {
		return
	}

	a.isListening = true
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	userCh, err := a.client.WS.StreamingWSUser(ctx)
	if err != nil {
		log.SLogger.Errorf("new user ws connction error: %s", err)
		return
	}
	defer close(userCh)

	for ntf := range userCh {
		switch ntf.(type) {
		case *gomastodon.NotificationEvent:
			n := ntf.(*gomastodon.NotificationEvent)
			if n.Notification.Type != "mention" {
				return
			}
			for _, handler := range a.NtfHandler {
				ntf := n.Notification
				if ntf.Status == nil {
					return
				}
				log.SLogger.Debugf("start execute handler: %v", handler)
				handler(a, ntf, actors)
			}
		default:
			// log.SLogger.Infof("receive other event: %s", uq)
		}
	}
}
