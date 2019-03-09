package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

var cfg TomlConfig
var actorInfoMap map[string]*ActorInfo

type TomlConfig struct {
	Title string
	ENV   string

	LogPath      string        `toml:"log_path"`
	ScriptFile   string        `toml:"script_file"`
	ActorsOnPlay []string      `toml:"actors_on_play"`
	TheaterName  string        `toml:"theater_name"`
	ActInterVal  time.Duration `toml:"act_interval"`

	ActorA ActorInfo `toml:"actor_a"`
	ActorB ActorInfo `toml:"actor_b"`
	ActorC ActorInfo `toml:"actor_c"`
	ActorD ActorInfo `toml:"actor_d"`
	// add more actors if you like
}

type ActorInfo struct {
	ID       string `toml:"client_id"`
	Secret   string `toml:"client_secret"`
	Sever    string `toml:"server"`
	Email    string `toml:"client_email"`
	Password string `toml:"client_password"`
	Name     string `toml:"name"`
}

func init() {
	runingEnv := flag.String("env", "development", "running env")
	rootedPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	fmt.Printf("Running in %s env\n", *runingEnv)
	configPath := fmt.Sprintf("%s/config/%s.toml", rootedPath, *runingEnv)

	dat, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("read config file error: %s", err)
		return
	}

	if err := toml.Unmarshal(dat, &cfg); err != nil {
		log.Fatal(err)
	}
	cfg.ENV = *runingEnv

	actorInfoMap = map[string]*ActorInfo{
		cfg.ActorA.Name: &cfg.ActorA,
		cfg.ActorB.Name: &cfg.ActorB,
		cfg.ActorC.Name: &cfg.ActorC,
		cfg.ActorD.Name: &cfg.ActorD,
	}
}

func GetRuntimeEnv() string {
	return cfg.ENV
}

func LogPath() string {
	return cfg.LogPath
}

func ActorsOnPlay() []string {
	return cfg.ActorsOnPlay
}

func TheaterName() string {
	return cfg.TheaterName
}

func ActInterVal() time.Duration {
	return cfg.ActInterVal * time.Minute
}

func ActorBotClientInfo(name string) (ActorInfo, error) {
	if acInfo, ok := actorInfoMap[name]; ok {
		return *acInfo, nil
	}
	return ActorInfo{}, fmt.Errorf("no such actor %s", name)
}

func ScriptFilePath() string {
	return cfg.ScriptFile
}
