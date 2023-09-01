package main

import (
	"github.com/rlapenok/test/internal/bot"
	"github.com/rlapenok/test/internal/config"
)

func main() {
	//Загрузка конфига
	cfg := config.New()
	//Создание бота
	bot := bot.New(cfg.Token, cfg.Port)
	//Запуск бота
	bot.Run()
}
