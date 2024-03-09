package main

import (
	"database/sql"
	"fmt"
	"runtime"

	_ "github.com/mailru/go-clickhouse"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

func main() {
	// Подключение к ClickHouse
	connect, err := sql.Open("clickhouse", "http://127.0.0.1:8123/default")
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting client to clickhouse")
	}
	if err = connect.Ping(); err != nil {
		log.Fatal().Err(err).Msg("error sending ping to clickhouse")
	}
	log.Info().Msg("client connected to Clickhouse")

	// Подключение к NATS
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting client to Nats")
	}
	log.Info().Msg("client connected to Nats")

	//Подписка на NATS тему для получения логов
	sub, err := nc.Subscribe("goods", func(msg *nats.Msg) {
		// Обработка полученного сообщения
		log.Trace().Msg(fmt.Sprintf("Получено сообщение: %s", string(msg.Data)))
		_, err = connect.Exec(string(msg.Data))
		if err != nil {
			log.Warn().Err(err).Msg("failed executing query")
		}
	})

	// Run forever
	runtime.Goexit()
	defer nc.Close()
	defer sub.Unsubscribe()
	defer connect.Close()
}
