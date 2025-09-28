package main

import (
	"log"
	"net/http"
	"os"

	"telenotion/api"
	"telenotion/internal/telegram"
	"telenotion/internal/utils"
)

func main() {
	if err := utils.LoadEnv(".env"); err != nil {
		log.Println("Could not load .env:", err)
	}

	tgToken := os.Getenv("TELEGRAM_TOKEN")
	if tgToken == "" {
		log.Fatal("TELEGRAM_TOKEN not set")
	}

	ntToken := os.Getenv("NOTION_TOKEN")
	if ntToken == "" {
		log.Fatal("NOTION_TOKEN not set")
	}
	dbID := os.Getenv("NOTION_DB_ID")
	if dbID == "" {
		log.Fatal("NOTION_DB_ID not set")
	}

	tg := telegram.NewTelegramClient(tgToken, nil)

	http.Handle("/telegram", api.TelegramHandler(tg))

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
