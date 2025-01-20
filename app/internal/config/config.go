package config

import (
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

const (
	cfgPath = "./config/config.yaml"
)

type Config struct {
	Env       string `yaml:"env"`
	LogFolder string `yaml:"log_folder"`
	BX24      BX24   `yaml:"bx24"`
}

type BX24 struct {
	Domain         string         `yaml:"bx24_domain"`
	ServerHost     string         `yaml:"server_host"`
	ServerPort     string         `yaml:"server_port"`
	RequestCounter RequestCounter `yaml:"leacky_bucket"`
	Timeout        time.Duration  `yaml:"timeout"`
}

type RequestCounter struct {
	Decrement int `yaml:"decrement"`
	Max       int `yaml:"max"`
}

var (
	instance Config
	once     sync.Once
)

func GetConfig() Config {
	once.Do(func() {
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			panic(err)
		}

		if err := cleanenv.ReadConfig(cfgPath, &instance); err != nil {
			panic(err)
		}
	})

	return instance
}
