package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"theater/bot"
	"theater/bredis"
	"theater/config"
	cons "theater/const"
	"theater/log"
	"time"

	"github.com/go-redis/redis"
)

const (
	NightStart = 12
	NightEnd   = 9
	Timeout    = 365 * 24 * 2 * time.Hour
)

func sendLine(actors map[string]*bot.Actor) {
	defer func() {
		for _, actor := range actors {
			close(actor.LineCh)
		}
	}()
	defer wg.Done()

	var id int
	var prevEp string

	filename := config.ScriptFilePath()
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	input := bufio.NewScanner(f)
	for input.Scan() {
		id++
		content := input.Text()
		ep, name, line, err := parseText(content)
		if err != nil {
			log.SLogger.Errorf("parse text:[%s] error: %v", content, err)
			continue
		}

		if ep != prevEp {
			id = 1
			prevEp = ep
		}

		acted, err := checkActed(ep, id)
		if acted || err != nil {
			continue
		}

		for checkNight() {
			time.Sleep(5 * time.Minute)
		}

		actor, ok := actors[name]
		if !ok {
			log.SLogger.Errorf("not find actor by name: %s on line id: %d", name, id)
			continue
		}

		select {
		case actor.LineCh <- line:
			log.SLogger.Infof("acts ep %s id %d", ep, id)
		default:
			log.SLogger.Errorf("actor %s LineCh blocked with line id: %d", actor.Name, id)
		}
	}
}

/*
line example:

ep/name/line
*/
func parseText(content string) (string, string, string, error) {
	s := strings.Split(content, "/")
	if len(s) < 3 {
		return "", "", "", fmt.Errorf("split content [%s] len less 3 error", content)
	}

	ep, name, line := s[0], s[1], s[2]
	return ep, name, line, nil
}

// it is based on the line number(id)
func checkActed(ep string, id int) (bool, error) {
	key := fmt.Sprintf("%s:%s", cons.Stein, ep)
	value, err := bredis.Client.Get(key).Result()
	if err == nil {
		valueInt, _ := strconv.Atoi(value)

		if id <= valueInt {
			return true, nil
		}

		// TODO: improve this place
		time.Sleep(config.ActInterVal())

		err := bredis.Client.Set(key, id, Timeout).Err()
		if err != nil {
			log.SLogger.Errorf("set ep %s with id %d from redis error: %v", ep, id, err)
			return false, err
		}
		return false, nil

	} else if err == redis.Nil {
		err := bredis.Client.Set(key, id, Timeout).Err()
		if err != nil {
			log.SLogger.Errorf("set ep %s with id %d from redis error: %v", ep, id, err)
			return false, err
		}
		return false, nil
	}

	log.SLogger.Errorf("get ep %s with id %d from redis error: %v", ep, id, err)
	return false, err
}

func checkNight() bool {
	now := time.Now()

	if now.Hour() >= NightStart && now.Hour() < NightEnd {
		return true
	}
	return false
}
