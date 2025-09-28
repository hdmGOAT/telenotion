package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"telenotion/internal/telegram"
)

func TelegramHandler(cmd *telegram.Commands) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var u telegram.Update
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		go func() {
			defer cancel()
			cmd.HandleCommands(ctx, u)
		}()

		w.WriteHeader(http.StatusOK)
	})
}
