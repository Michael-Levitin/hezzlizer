package main

import (
	"fmt"
	"runtime"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

func main() {
	//// Подключение к ClickHouse
	//connect, err := clickhouse.Open("tcp://clickhouse-server:9000?debug=false")
	//if err != nil {
	//	log.Fatal().Err(err).Msg("error connecting to ClickHouse")
	//}
	//defer connect.Close()

	// Подключение к NATS
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting client to Nats")
	}
	log.Info().Msg("client connected to Nats")

	//Подписка на NATS тему для получения логов
	sub, err := nc.Subscribe("goods", func(msg *nats.Msg) {
		// Обработка полученного сообщения
		log.Info().Msg(fmt.Sprintf("Получено сообщение: \n%s", msg.Data))
		//stmt, err := connect.Prepare(string(msg.Data))

	})

	// Run forever
	runtime.Goexit()
	defer nc.Close()
	defer sub.Unsubscribe()
}
