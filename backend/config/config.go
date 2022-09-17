package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
    DB struct {
        Type string      `env:"DB_TYPE"`
        Host string      `env:"DB_HOST"`
        Port string      `env:"DB_PORT"`
        Username string  `env:"DB_USERNAME"`
        Password string  `env:"DB_PASSWORD"`
        DBName string    `env:"DB_DBNAME"`
    }
    Routing struct {
        Port string `env:"ROUTING_PORT"`
    }
    Session struct {
        Name string    `env:"SESSION_NAME"`
        Secret string  `env:"SESSION_SECRET"`
    }
    Twitter struct {
        ConsumerKey string     `env:"TWITTER_API_KEY"`
        ConsumerSecret string  `env:"TWITTER_API_KEY_SECRET"`
        CallbackUrl string     `env:"TWITTER_CALLBACK"`
    }
}

func NewConfig() *Config {
    if err := godotenv.Load(); err != nil {
        // TODO: エラーログ出力
        log.Print("No .env file found")
    }

    var conf Config
    if err := env.Parse(&conf); err != nil {
        panic(err)
    }

    return &conf
}