package bot

import (
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gookit/slog"
	"github.com/rlapenok/test/internal/grpc_server"
	auto_generate "github.com/rlapenok/test/internal/library/auto_generate/tg/proto"
)

type MyInterface interface {
	Run()
	PrintLog(chat_id int64) string
	ReadUpdate()
}

type TgBot struct {
	grpc_server *grpc_server.GrpcServer
	bot         *tgbotapi.BotAPI
	updates     tgbotapi.UpdatesChannel
}

func New(token string, port string) *TgBot {
	//Создание тг бота из пакета
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Fatal(err.Error())
		os.Exit(1)
	}
	//Сделать описание
	update_config := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(update_config)
	if err != nil {
		slog.Fatal(err.Error())
	}
	//Создание gprcServer
	grpc_server := grpc_server.NewServer(port)
	tg_bot := TgBot{bot: bot, updates: updates, grpc_server: grpc_server}
	return &tg_bot
}
func (bot *TgBot) Run() {
	slog.Info("Starting tg_bot...")
	//Создание wg
	wg := sync.WaitGroup{}
	wg.Add(1)
	//Запуск grpcServer
	go bot.grpc_server.Run()
	//Чтение входящих updates на бот
	bot.ReadUpdate()
	wg.Wait()

}
func (bot *TgBot) ReadUpdate() {
	//Создание WaitGroup для слежения за горутинами
	wg := sync.WaitGroup{}
	for update := range bot.updates {
		wg.Add(1)
		//Обработка каждого update в отдельной горутине
		go bot.HandleUpdates(update, &wg)
	}
	wg.Wait()

}
func (bot *TgBot) HandleUpdates(update tgbotapi.Update, wg *sync.WaitGroup) {
	//Разбираем тип update
	//Если это команда
	if update.Message != nil && update.Message.IsCommand() {
		chat_id := update.Message.Chat.ID
		{
			switch update.Message.Command() {
			case "start":
				{
					//Проверяем, имеет ли пользователь свой канал для общения
					exist := bot.grpc_server.ChechChatId(chat_id)
					//Если уже chat_id существует
					if exist {
						msg := tgbotapi.NewMessage(chat_id, "Вы уже запустили меня")
						bot.bot.Send(msg)
						defer wg.Done()
					} else {
						resp := "Начинаю печать логов"
						msg := tgbotapi.NewMessage(chat_id, resp)
						bot.bot.Send(msg)
						//Если не существует, то создаем для него канал для общения
						channel := make(chan *auto_generate.PrintLogRequest)
						bot.grpc_server.AddChatId(chat_id, channel)
						//Слушаем канал и отдаем сообщения пользователю
						for message := range channel {
							//Преобразуем в читаемый вид
							resp := "service_name:" + message.ServiceName + "\n" + "message:" + message.Message
							msg := tgbotapi.NewMessage(chat_id, resp)
							//Создаем для сообщения клавиатуру
							msg.ReplyMarkup = create_keyboard()
							//Отправляем сообщение
							bot.bot.Send(msg)

						}
						defer wg.Done()

					}
				}
			case "stop":
				{
					//Проверяем, имеет ли пользователь свой канал для общения
					exist := bot.grpc_server.ChechChatId(chat_id)
					if !exist {
						msg := tgbotapi.NewMessage(chat_id, "Вы уже остановили меня")
						bot.bot.Send(msg)
						defer wg.Done()
					} else {
						//Закрываем и удаляем канал
						bot.grpc_server.DeleteChatId(chat_id)
						resp := "Бот остановлен"
						msg := tgbotapi.NewMessage(chat_id, resp)
						bot.bot.Send(msg)
						defer wg.Done()
					}
				}
			case "help":
				{
					resp := "Список доступных команд:\n\n" + "/start - Старт бота\n" + "/help - Показать список доступных команд\n" + "/stop - Остановить бота \n"
					msg := tgbotapi.NewMessage(chat_id, resp)
					bot.bot.Send(msg)
				}
			default:
				{
					wg.Done()
				}
			}

		}
	} else if update.CallbackQuery != nil {
		call_back := update.CallbackQuery.Data
		switch call_back {
		case "0":
			slog.Info("CallBack:Approve")
		case "1":
			slog.Info("CallBack:Delete")
		}
		defer wg.Done()
	} else {
		wg.Done()
	}

}

func create_keyboard() tgbotapi.InlineKeyboardMarkup {
	callback_btn_approve := "0"
	callback_btn1_delete := "1"

	btn_1 := tgbotapi.InlineKeyboardButton{Text: "Approve", CallbackData: &callback_btn_approve}
	btn_2 := tgbotapi.InlineKeyboardButton{Text: "Delete", CallbackData: &callback_btn1_delete}
	arr := [2]tgbotapi.InlineKeyboardButton{btn_1, btn_2}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(arr[:])
	return keyboard
}
