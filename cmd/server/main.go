package main

import (
	"context"
	"fmt"
	"github.com/Michael-Levitin/hezzlizer/config"
	"github.com/Michael-Levitin/hezzlizer/internal/database"
	"github.com/Michael-Levitin/hezzlizer/internal/delivery"
	"github.com/Michael-Levitin/hezzlizer/internal/logic"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
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

	hezzlDB := database.NewHezzlDB(db)                 // подключаем бд
	hezzlLogic := logic.NewHezzlLogic(hezzlDB)         // подключаем бд к логике...
	hezzlServer := delivery.NewHezzlServer(hezzlLogic) // ... а логику в delivery

	http.HandleFunc("/good/create", hezzlServer.GoodCreate)
	http.HandleFunc("/good/update", hezzlServer.GoodUpdate)
	http.HandleFunc("/good/remove", hezzlServer.GoodRemove)
	http.HandleFunc("/goods/list", hezzlServer.GoodsList)
	http.HandleFunc("/good/reprioritize", hezzlServer.GoodReprioritize)
	log.Info().Msg("server is running...")
	err = http.ListenAndServe(":8080", nil)
	log.Fatal().Err(err).Msg("http server crashed")
}
