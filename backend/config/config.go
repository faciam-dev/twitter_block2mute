package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const projectDirName = "backend"

type Config struct {
	DB struct {
		Type     string `env:"DB_TYPE"`
		Host     string `env:"DB_HOST"`
		Port     string `env:"DB_PORT"`
		Username string `env:"DB_USERNAME"`
		Password string `env:"DB_PASSWORD"`
		DBName   string `env:"DB_DBNAME"`
	}
	Routing struct {
		AllowOrigins []string `env:"ALLOW_ORIGINS" envSeparator:","`
		Port         string   `env:"ROUTING_PORT"`
	}
	Session struct {
		Name   string `env:"SESSION_NAME"`
		Secret string `env:"SESSION_SECRET"`
	}
	Twitter struct {
		ConsumerKey    string `env:"TWITTER_API_KEY"`
		ConsumerSecret string `env:"TWITTER_API_KEY_SECRET"`
		CallbackUrl    string `env:"TWITTER_CALLBACK_URL"`
	}
}

func NewConfig(envFilename string) *Config {
	loadEnv(envFilename)
	/*
		if err := godotenv.Load(".env"); err != nil {
			// TODO: エラーログ出力
			log.Print("No .env file found")
		}
	*/

	fmt.Println()
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0] + "=" + pair[1])
	}

	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		panic(err)
	}

	return &conf
}

func loadEnv(envFilename string) {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/` + envFilename)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
