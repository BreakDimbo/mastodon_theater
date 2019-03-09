package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	cons "theater/const"

	"github.com/BurntSushi/toml"
)

var config TomlConfig

type TomlConfig struct {
	Title string
	ENV   string

	ScriptFile   string   `toml:"script_file"`
	ActorsOnPlay []string `toml:"actors_on_play"`
	TheaterName  string   `toml:"theater_name"`

	ActorA MastodonClientInfo `toml:"actor_a"`
	ActorB MastodonClientInfo `toml:"actor_b"`
	ActorC MastodonClientInfo `toml:"actor_c"`
	ActorD MastodonClientInfo `toml:"actor_d"`
	ActorE MastodonClientInfo `toml:"actor_e"`
	ActorF MastodonClientInfo `toml:"actor_f"`
	ActorG MastodonClientInfo `toml:"actor_g"`
	ActorH MastodonClientInfo `toml:"actor_h"`
	ActorI MastodonClientInfo `toml:"actor_i"`
	ActorJ MastodonClientInfo `toml:"actor_j"`
	ActorK MastodonClientInfo `toml:"actor_k"`
	ActorL MastodonClientInfo `toml:"actor_l"`
	ActorM MastodonClientInfo `toml:"actor_m"`
	ActorN MastodonClientInfo `toml:"actor_n"`
}

type MastodonClientInfo struct {
	ID       string `toml:"client_id"`
	Secret   string `toml:"client_secret"`
	Sever    string `toml:"server"`
	Email    string `toml:"client_email"`
	Password string `toml:"client_password"`
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

	if err := toml.Unmarshal(dat, &config); err != nil {
		log.Fatal(err)
	}
	config.ENV = *runingEnv
}

func GetRuntimeEnv() string {
	return config.ENV
}

func ActorsOnPlay() []string {
	return config.ActorsOnPlay
}

func TheaterName() string {
	return config.TheaterName
}

func ActorBotClientInfo(name string) (MastodonClientInfo, error) {
	switch name {
	case cons.Okabe:
		return config.ActorA, nil
	case cons.Mayuri:
		return config.ActorB, nil
	case cons.Itaru:
		return config.ActorC, nil
	case cons.Kurisu:
		return config.ActorD, nil
	case cons.Moeka:
		return config.ActorE, nil
	case cons.Ruka:
		return config.ActorF, nil
	case cons.NyanNyan:
		return config.ActorG, nil
	case cons.Suzuha:
		return config.ActorH, nil
	case cons.Maho:
		return config.ActorI, nil
	case cons.Kagari:
		return config.ActorJ, nil
	case cons.Yuki:
		return config.ActorK, nil
	case cons.Tennouji:
		return config.ActorL, nil
	case cons.Nae:
		return config.ActorM, nil
	case cons.Nakabachi:
		return config.ActorN, nil
	default:
		return MastodonClientInfo{}, fmt.Errorf("no such actor %s", name)
	}
}

func ScriptFilePath() string {
	return config.ScriptFile
}
