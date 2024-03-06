package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Michael-Levitin/hezzlizer/config"
	"github.com/Michael-Levitin/hezzlizer/internal/database"
	"github.com/Michael-Levitin/hezzlizer/internal/delivery"
	"github.com/Michael-Levitin/hezzlizer/internal/logic"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// загружаем конфиг
	config.Init()
	sc := config.New()
	//logger := zerolog.New(os.Stdout)
	zerolog.SetGlobalLevel(sc.LogLevel)

	// подключаемся к базе данных
	dbAdrr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.DbUsername, sc.DbPassword, sc.DbHost, sc.DbPort, sc.DbName)
	db, err := pgxpool.New(context.TODO(), dbAdrr)
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to database")
	}
	log.Info().Msg("connected to database")
	defer db.Close()

	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to Redis")
	}
	log.Info().Msg("connected to Redis")
	defer client.Close()

	// Подключение к NATS
	//nc, err := nats.Connect("nats://nats-server:4222")
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting server to Nats")
	}
	log.Info().Msg("server connected to Nats")
	defer nc.Close()

	conn := database.NewNatsSender(nc)                          // подключаем Nats
	go conn.Send()                                              // включаем Send в горутине
	hezzlDB := database.NewHezzlDB(db)                          // подключаем бд
	redisDB := database.NewRedisDB(client)                      // подключаем Redis
	hezzlLogic := logic.NewHezzlLogic(hezzlDB, redisDB)         // подключаем бд и Redis к логике...
	hezzlServer := delivery.NewHezzlServer(hezzlLogic, redisDB) // ... а логику в delivery

	http.HandleFunc("/good/create", hezzlServer.GoodCreate)
	http.HandleFunc("/good/update", hezzlServer.GoodUpdate)
	http.HandleFunc("/good/remove", hezzlServer.GoodRemove)
	http.HandleFunc("/goods/list", hezzlServer.GoodsList)
	http.HandleFunc("/good/reprioritize", hezzlServer.GoodReprioritize)
	log.Info().Msg("server is running...")
	err = http.ListenAndServe(":8080", nil)
	log.Fatal().Err(err).Msg("http server crashed")
}
