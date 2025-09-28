package api

import (
	"encoding/json"
	"log"
	"net/http"

	"telenotion/internal/telegram"
)

func TelegramHandler(tg *telegram.TelegramClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var u telegram.Update
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if u.Message != nil {
			log.Printf("Received: %q", u.Message.Text)
			if err := tg.SendMessage(r.Context(), u.Message.Chat.ID,
				"You said: "+u.Message.Text); err != nil {
				log.Println("telegram send error:", err)
			}
		}
		w.WriteHeader(http.StatusOK)
	})
}
