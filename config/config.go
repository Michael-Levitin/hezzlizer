package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DbHost     string
	DbPort     string
	DbName     string
	DbUsername string
	DbPassword string
	LogLevel   zerolog.Level
}

func Init() {
	// загружаем данные из .env файла в систему
	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("no .env file found")
	} else {
		log.Info().Msg("loaded env values")
	}
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		DbHost:     getEnv("DB_HOST", "localhost"),
		DbPort:     getEnv("DB_PORT", "5432"),
		DbName:     getEnv("DB_NAME", "postgres"),
		DbUsername: getEnv("DB_USERNAME", "postgres"),
		DbPassword: getEnv("DB_PASSWORD", "postgres"),
		LogLevel:   getLevel("LOG_LEVEL", "trace"),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	var value string
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if value == "" && defaultVal == "" {
		log.Fatal().Msg(fmt.Sprint(key, " value not found"))
	}
	return defaultVal
}

func getLevel(key string, defaultVal string) zerolog.Level {
	var userLevel string
	if value, exists := os.LookupEnv(key); exists {
		userLevel = value
	}

	levelS := map[string]zerolog.Level{
		"trace":    zerolog.TraceLevel,
		"info":     zerolog.InfoLevel,
		"warn":     zerolog.WarnLevel,
		"error":    zerolog.ErrorLevel,
		"fatal":    zerolog.FatalLevel,
		"panic":    zerolog.PanicLevel,
		"nolevel":  zerolog.NoLevel,
		"disabled": zerolog.DebugLevel,
	}

	if level, exists := levelS[userLevel]; exists {
		log.Info().Msg(fmt.Sprint("setting log level to ", userLevel))
		return level
	}

	if level, exists := levelS[defaultVal]; exists {
		log.Warn().Msg(fmt.Sprint("user log level not found, setting default value - ", defaultVal))
		return level
	}
	log.Warn().Msg("log levels not found, setting log level to info")
	return levelS["info"]
}
