package main

import (
	"log"
	"os"

	"github.com/alexkarpovich/go-chinese-dict-bot/bot"
	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	if os.Getenv("TELEGRAM_BOT_TOKEN") == "" {
		log.Print("TELEGRAM_BOT_TOKEN is required")
		os.Exit(1)
	}
}

func main() {
	bot.Start()
}
