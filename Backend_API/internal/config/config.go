package config

import (
	"flag"

	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type HttpServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct {
	Env             string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath     string `yaml:"storage_path" env-required:"true"`
	UserStoragePath string `yaml:"user_storage_path" env-required:"true"`
	JwtSecrateKey   string `yaml:"jwt_secrate_key" env-required:"true"`
	RedisPassword   string `yaml:"redisPassword" env-required:"true"`

	HttpServer HttpServer `yaml:"http_server"`
}

func MustDone() *Config {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flgpath := flag.String("config_path", "", "Path to config file")
		flag.Parse()
		configPath = *flgpath
		if configPath == "" {
			log.Fatal("Invalid config path from flag")
		}
	}
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		log.Fatal("Unable to find the config file")
	}

	var cnf Config

	err = cleanenv.ReadConfig(configPath, &cnf)
	if err != nil {
		log.Fatal("unable to parse the config file ")
	}

	return &cnf
}
