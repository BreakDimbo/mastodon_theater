package main

import (
	"sync"
	"theater/bot"
	"theater/config"
	cons "theater/const"
)

var wg sync.WaitGroup

func main() {
	actors := make(map[string]*bot.Actor)
	actorsName := config.ActorsOnPlay()

	// you can DIY your actors here by adding different handlers
	for _, name := range actorsName {
		var actor *bot.Actor
		switch name {
		case cons.Kurisu:
			actor = bot.New(name, bot.LoveHandler)
		case cons.Okabe:
			actor = bot.New(name, bot.BlockHandler, bot.UnblockHandler)
		case cons.Itaru:
			actor = bot.New(name, bot.FoodHandler)
		default:
			actor = bot.New(name)
		}

		actors[name] = actor
		wg.Add(1)
		go actor.Act(&wg)
		if name == cons.Okabe || name == cons.Kurisu || name == cons.Itaru {
			go actor.ListenAudiences(actors)
		}
	}

	wg.Add(1)
	go sendLine(actors)

	wg.Wait()
}
