package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"telenotion/api"
	"telenotion/internal/notion"
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
	hp := &http.Client{
		Timeout: 30 * time.Second,
	}

	tg := telegram.NewTelegramClient(tgToken, hp)
	nt := notion.NewNotionClient(ntToken, hp)

        nC := notion.NewService(nt)	
	
	tC := telegram.NewCommands(tg, nC, dbID)
	http.Handle("/telegram", api.TelegramHandler(tC))

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
