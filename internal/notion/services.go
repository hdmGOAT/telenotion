package notion

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	Client *NotionClient
}

func NewService(c *NotionClient) *Service {
	return &Service{Client: c}
}

func (s *Service) UpcomingTodos(ctx context.Context, dbID string, days int) ([]Todo, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 0, days)

	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	filter := map[string]any{
		"and": []any{
			map[string]any{"property": "deadline", "date": map[string]any{"on_or_after": startStr}},
			map[string]any{"property": "deadline", "date": map[string]any{"on_or_before": endStr}},
		},
	}


	todos, err := s.Client.GetToDos(ctx, dbID, filter)
	if err == nil {
	    return todos, nil
	}
	fmt.Printf("[Warning] Filtered GetToDos failed: %v, falling back to unfiltered\n", err)

	all, err2 := s.Client.GetToDos(ctx, dbID, nil)
	if err2 != nil {
	    return nil, fmt.Errorf("filtered GetToDos failed: %w; unfiltered GetToDos also failed: %v", err, err2)
	}


	var out []Todo
	for _, t := range all {
		if t.Deadline.IsZero() {
			continue
		}
		if !t.Deadline.Before(startDate) && !t.Deadline.After(endDate) {
			out = append(out, t)
		}
	}
	return out, nil
}

func (s *Service) DueToday(ctx context.Context, dbID string) ([]Todo, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayStr := today.Format("2006-01-02")

	filter := map[string]any{
		"property": "deadline",
		"date":     map[string]any{"equals": todayStr},
	}

	todos, err := s.Client.GetToDos(ctx, dbID, filter)
	if err == nil {
		return todos, nil
	}

	all, err2 := s.Client.GetToDos(ctx, dbID, nil)
	if err2 != nil {
		return nil, err
	}

	var out []Todo
	for _, t := range all {
		if t.Deadline.IsZero() {
			continue
		}
		y1, m1, d1 := t.Deadline.Date()
		y2, m2, d2 := today.Date()
		if y1 == y2 && m1 == m2 && d1 == d2 {
			out = append(out, t)
		}
	}
	return out, nil
}

func (s *Service) Pending(ctx context.Context, dbID string) ([]Todo, error) {
	filter := map[string]any{
		"or": []any{
			map[string]any{
				"and": []any{
					map[string]any{"property": "Hans Progress", "status": map[string]any{"does_not_equal": "submitted"}},
					map[string]any{"property": "Hans Progress", "status": map[string]any{"does_not_equal": "N/A"}},
				},
			},
			map[string]any{
				"and": []any{
					map[string]any{"property": "Ira Progress", "status": map[string]any{"does_not_equal": "submitted"}},
					map[string]any{"property": "Ira Progress", "status": map[string]any{"does_not_equal": "N/A"}},
				},
			},
		},
	}

	todos, err := s.Client.GetToDos(ctx, dbID, filter)
	if err == nil {
		return todos, nil
	}

	all, err2 := s.Client.GetToDos(ctx, dbID, nil)
	if err2 != nil {
		return nil, err
	}

	var out []Todo
	for _, t := range all {
		hans := t.Progress.Hans
		ira := t.Progress.Ira

		hansPending := hans != "submitted" && hans != "N/A"
		iraPending  := ira  != "submitted" && ira  != "N/A"

		if hansPending || iraPending {
			out = append(out, t)
		}
	}
	return out, nil
}
