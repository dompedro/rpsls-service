package rpslsapi

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Config config

type config struct {
	Server             ServerConfig
	DB                 DatabaseConfig
	Redis              RedisConfig
	RandomNumberServer string
	ScoreboardSize     int
	Environment        string
	DummyUserID        string // dummyUserID is a fixed userId to be used in single-player mode
}

type ServerConfig struct {
	Addr string
}

type DatabaseConfig struct {
	Uri      string
	Username string
	Password string
	Database string
	Realm    string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func LoadConfig() {
	env := os.Getenv("RPSLS_ENV")
	if "" == env {
		env = "development"
	}

	loadFiles(env)
	Config = config{
		Server: ServerConfig{
			Addr: os.Getenv("SERVER_ADDR"),
		},
		DB: DatabaseConfig{
			Uri:      os.Getenv("DB_URI"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_DATABASE"),
			Realm:    os.Getenv("DB_REALM"),
		},
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       intConfig("REDIS_DB"),
		},
		RandomNumberServer: os.Getenv("RANDOM_NUMBER_SERVER"),
		ScoreboardSize:     intConfig("RPSLS_SCOREBOARD_SIZE"),
		Environment:        env,
		DummyUserID:        "a4868d93-2d71-4ce4-b48c-c70e6a043851",
	}
}

func intConfig(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		panic(fmt.Errorf("env var %s must be an int", key))
	}
	return value
}

func loadFiles(env string) {
	_ = godotenv.Load(".env." + env + ".local")
	if "test" != env {
		_ = godotenv.Load(".env.local")
	}

	if err := godotenv.Load(".env." + env); err != nil {
		panic("failed to load env configuration file")
	}
	if err := godotenv.Load(); err != nil {
		panic("failed to load base configuration file")
	}
}
