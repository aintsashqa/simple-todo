package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	conf Config
	once sync.Once
)

type Config struct {
	Server   Server
	Database Database
}

type Server struct {
	Host           string        `env:"APP_HOST" env-default:"127.0.0.1"`
	Port           int           `env:"APP_PORT" env-default:"8080"`
	ReadTimeout    time.Duration `env:"APP_READ_TIMEOUT" env-default:"1s"`
	WriteTimeout   time.Duration `env:"APP_WRITE_TIMEOUT" env-default:"1s"`
	MaxHeaderBytes uint8         `env:"APP_MAX_HEADER_BYTES" env-default:"1"`
}

type Database struct {
	Host     string `env:"DB_HOST" env-default:"127.0.0.1"`
	Port     int    `env:"DB_PORT" env-default:"5432"`
	Username string `env:"DB_USERNAME" env-default:"root"`
	Password string `env:"DB_PASSWORD" env-default:"secret"`
	Name     string `env:"DB_NAME" env-required:"true"`
}

func GetConfig() Config {
	once.Do(func() {
		if err := cleanenv.ReadConfig(".env", &conf); err != nil {
			panic(fmt.Sprintf("unable read environment variables to config: %s", err))
		}
	})

	return conf
}
