package config

import (
	"log"
	"os"
	"regexp"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const projectDirName = "backend"

type Config struct {
	ReleaseMode string `env:"RELEASE_MODE"`
	Logger      struct {
		LoggerLevel            string   `env:"LOGGER_LEVEL"`
		LoggerOutputPaths      []string `env:"LOGGER_OUTPUT_PATHS" envSeparator:","`
		LoggerErrorOutputPaths []string `env:"LOGGER_ERROR_OUTPUT_PATHS" envSeparator:","`
	}
	DB struct {
		Type     string `env:"DB_TYPE"`
		Host     string `env:"DB_HOST"`
		Port     string `env:"DB_PORT"`
		Username string `env:"DB_USERNAME"`
		Password string `env:"DB_PASSWORD"`
		DBName   string `env:"DB_DBNAME"`
	}
	Routing struct {
		AllowOrigins   []string `env:"ALLOW_ORIGINS" envSeparator:","`
		AllowHeaders   []string `env:"ALLOW_HEADERS" envSeparator:","`
		AllowMethods   []string `env:"ALLOW_METHODS" envSeparator:","`
		ExposeHeaders  []string `env:"EXPOSE_HEADERS" envSeparator:","`
		TrustedProxies []string `env:"TRUSTED_PROXIES" envSeparator:","`
		MaxAge         int      `env:"CORS_MAX_AGE"`
		Port           string   `env:"ROUTING_PORT"`
		CsrfSecure     bool     `env:"CSRF_SECURE"`
	}
	Session struct {
		Name   string `env:"SESSION_NAME"`
		Secret string `env:"SESSION_SECRET"`
	}
	Twitter struct {
		ConsumerKey    string `env:"GOTWI_API_KEY"`
		ConsumerSecret string `env:"GOTWI_API_KEY_SECRET"`
		CallbackUrl    string `env:"TWITTER_CALLBACK_URL"`
	}
}

func NewConfig(envFilename string) *Config {
	loadEnv(envFilename)

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
