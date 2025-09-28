package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"telenotion/internal/notion"
)

type Commands struct {
	Bot     *TelegramClient
	Service *notion.Service
	DBID    string
}

func NewCommands(bot *TelegramClient, svc *notion.Service, dbID string) *Commands {
	return &Commands{Bot: bot, Service: svc, DBID: dbID}
}

func (c *Commands) HandleCommands(ctx context.Context, u Update) {
	if u.Message == nil || !strings.HasPrefix(u.Message.Text, "/") {
		return
	}

	parts := strings.Fields(strings.TrimPrefix(u.Message.Text, "/"))
	if len(parts) == 0 {
		return
	}
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "upcoming":
		c.handleUpcoming(ctx, u, parts[1:])
	case "today":
		c.handleToday(ctx, u)
	case "pending":
		c.handlePending(ctx, u)
	default:
		c.Bot.SendMessage(ctx, u.Message.Chat.ID, "Unknown command")
	}
}

func (c *Commands) handleUpcoming(ctx context.Context, u Update, args []string) {
	days := 7
	if len(args) > 0 {
		if n, err := fmt.Sscanf(args[0], "%d", &days); n != 1 || err != nil {
			days = 7
		}
	}
	todos, err := c.Service.UpcomingTodos(ctx, c.DBID, days)
	if err != nil {
		log.Println("upcoming:", err)
		c.Bot.SendMessage(ctx, u.Message.Chat.ID, "error  fetfhing pending todos")
		return
	}
	if len(todos) == 0 {
		c.Bot.SendMessage(ctx, u.Message.Chat.ID, "no upcoming tasks")
		return
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Tasks in next %d days:\n", days))
	for _, t := range todos {
		b.WriteString(fmt.Sprintf("• %s (due %s)\n",
			t.Name, t.Deadline.Format("Jan 2")))
	}
	c.Bot.SendMessage(ctx, u.Message.Chat.ID, b.String())
}

func (c *Commands) handleToday(ctx context.Context, u Update) {
	todos, err := c.Service.DueToday(ctx, c.DBID)
	if err != nil {
		log.Println("today:", err)
		c.Bot.SendMessage(ctx, u.Message.Chat.ID, "ewrror fetching todays tasks")
		return
	}
	if len(todos) == 0 {
		c.Bot.SendMessage(ctx, u.Message.Chat.ID, "nmothing due today")
		return
	}
	var b strings.Builder
	b.WriteString("Due today:\n")
	for _, t := range todos {
		b.WriteString(fmt.Sprintf("• %s (%s)\n", t.Name,
			t.Deadline.Format(time.Kitchen)))
	}
	c.Bot.SendMessage(ctx, u.Message.Chat.ID, b.String())
}

func (c *Commands) handlePending(ctx context.Context, u Update) {
	todos, err := c.Service.Pending(ctx, c.DBID)
	if err != nil {
		log.Println("pending:", err)
		c.Bot.SendMessage(ctx, u.Message.Chat.ID, "eerror fetching pending tasks")
		return
	}
	if len(todos) == 0 {
		c.Bot.SendMessage(ctx, u.Message.Chat.ID, "no pending tasks")
		return
	}
	var b strings.Builder
	b.WriteString("Pending tasks:\n")
	for _, t := range todos {
		b.WriteString(fmt.Sprintf("• %s (Hans:%s Ira:%s)\n",
			t.Name, t.Progress.Hans, t.Progress.Ira))
	}
	c.Bot.SendMessage(ctx, u.Message.Chat.ID, b.String())
}
