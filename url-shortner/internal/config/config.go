package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yaml:"env" env-default:"local"`
	StoragePatch string `yaml:"storage_patch" env-required:"true"`
	HTTPServer   `yaml:"http_server"`
}

type HTTPServer struct {
	Adress      string        `yaml:"address" env-default:"localhost:8081"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad() *Config {
	configPatch := os.Getenv("CONFIG_PATH")

	if configPatch == "" {
		//log.Fatal("NOT CONFIGPATH")
		configPatch = "./url-shortner/config/local.yaml"
	}

	if _, err := os.Stat(configPatch); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exost", configPatch)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPatch, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
