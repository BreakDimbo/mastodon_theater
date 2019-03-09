package client

import (
	"context"
	"log"
	"strconv"
	"theater/config"
	gomastodon "theater/go-mastodon"
	zlog "theater/log"
	"time"

	"github.com/cenkalti/backoff"
)

type Bot struct {
	Normal *gomastodon.Client
	WS     *gomastodon.WSClient
}

// New new bot client which should be called only when init
func New(config *config.MastodonClientInfo) (*Bot, error) {
	c := gomastodon.NewClient(&gomastodon.Config{
		Server:       config.Sever,
		ClientID:     config.ID,
		ClientSecret: config.Secret,
	})

	clientAuth := func() error {
		return c.Authenticate(context.Background(), config.Email, config.Password)
	}

	bkStrategy := backoff.NewExponentialBackOff()
	bkStrategy.MaxInterval = 1 * time.Hour
	err := backoff.Retry(clientAuth, bkStrategy)
	if err != nil {
		log.Fatalf("[Fatal]: authenticate error of mastodon client: %s\n", err)
		return nil, err
	}

	bc := &Bot{Normal: c, WS: c.NewWSClient()}
	return bc, nil
}

func (bc *Bot) RawPost(toot *gomastodon.Toot) (*gomastodon.Status, error) {
	status, err := bc.Normal.PostStatus(context.Background(), toot)
	if err != nil {
		zlog.SLogger.Errorf("post toot: %s error: %s", toot, err)
		return nil, err
	}
	return status, nil
}

func (bc *Bot) BlockAccount(accountID string) (gomastodon.ID, error) {
	id := gomastodon.ID(accountID)
	status, err := bc.Normal.AccountBlock(context.Background(), id)
	if err != nil {
		zlog.SLogger.Errorf("block id: %s error: %s", id, err)
		return "", err
	}
	return status.ID, nil
}

func (bc *Bot) UnBlockAccount(accountID string) (gomastodon.ID, error) {
	id := gomastodon.ID(accountID)
	status, err := bc.Normal.AccountUnblock(context.Background(), id)
	if err != nil {
		zlog.SLogger.Errorf("unblock id: %s error: %s", id, err)
		return "", err
	}
	return status.ID, nil
}

func (bc *Bot) Post(toot string) (gomastodon.ID, error) {
	status, err := bc.Normal.PostStatus(context.Background(), &gomastodon.Toot{
		Status: toot,
	})
	if err != nil {
		zlog.SLogger.Errorf("post toot: %s error: %s", toot, err)
		return "", err
	}
	return status.ID, nil
}

func (bc *Bot) PostSpoiler(spolier string, toot string) (gomastodon.ID, error) {
	status, err := bc.Normal.PostStatus(context.Background(), &gomastodon.Toot{
		Status:      toot,
		SpoilerText: spolier,
	})
	if err != nil {
		zlog.SLogger.Errorf("post toot: %s error: %s", toot, err)
		return "", err
	}
	return status.ID, nil
}

func (bc *Bot) DeleteToot(id string) error {
	ctx := context.Background()
	fbotTootID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		zlog.SLogger.Errorf("parse id: %s error: %s", id, err)
		return err
	}
	return bc.Normal.DeleteStatus(ctx, int64(fbotTootID))
}
