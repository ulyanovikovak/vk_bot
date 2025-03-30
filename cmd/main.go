package main

import (
	"log"
	"net/http"
	"time"
	"vk_bot/bot"
	"vk_bot/logger"
	"vk_bot/poll"
	"vk_bot/storage"

	tarantool "github.com/tarantool/go-tarantool"
)

func main() {
	logger.Init()
	logger.Log.Info("Start Vote Bot")

	var conn *tarantool.Connection
	var err error
	opts := tarantool.Opts{User: "guest"}

	for i := 0; i < 10; i++ {
		conn, err = tarantool.Connect("tarantool:3301", opts)
		if err == nil {
			break
		}
		logger.Log.Info("Ожидание Tarantool, попытка:", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Ошибка подключения к Tarantool: %s", err)
	}
	defer conn.Close()

	tStorage := storage.NewTarantoolStorage(conn)
	pService := poll.NewPollService(tStorage)

	http.HandleFunc("/vote", bot.Handler(pService))
	logger.Log.Info("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
