package bot

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/alexkarpovich/go-chinese-dict-bot/scraper"
)

// Start : Entry point on Chinese Dict bot
func Start() {
	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		// универсальный ответ на любое сообщение
		reply := "Не знаю что сказать"
		if update.Message == nil {
			continue
		}

		// логируем от кого какое сообщение пришло
		log.Printf("[@%s] %s", update.Message.From.UserName, update.Message.Text)

		// свитч на обработку комманд
		// комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {
		case "start":
			reply = "Привет. Я телеграм-бот"
		default:
			info := scraper.Scrape(update.Message.Text)
			reply = fmt.Sprintf(`<b>%s</b>
			<i>курсив</i>
			<a href="http://www.example.com/">ссылка</a>
			<code>строчный моноширинный</code>
			<pre>блочный моноширинный (можно писать код)</pre>`, info)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "html"
		// отправляем
		bot.Send(msg)
	}
}
