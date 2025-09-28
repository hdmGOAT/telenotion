package notion

import (
	"time"
)

func MapNotionToTodos(resp DatabaseQueryResponse) ([]Todo, error) {
	now := time.Now()
	var out []Todo

	for _, r := range resp.Results {
		var deadline time.Time
		if s := r.Properties.Deadline.Date.Start; s != "" {
			if d, err := time.Parse("2006-01-02", s); err == nil {
				deadline = d
			}
		}

		daysLeft := 0
		if !deadline.IsZero() {
			d := int(deadline.Sub(now).Hours() / 24)
			if d > 0 {
				daysLeft = d
			}
		}

		name := ""
		if len(r.Properties.Name.Title) > 0 {
			name = r.Properties.Name.Title[0].PlainText
		}

		course := ""
		if len(r.Properties.Course.Relation) > 0 {
			course = r.Properties.Course.Relation[0].ID
		}

		notes := ""
		if len(r.Properties.Notes.RichText) > 0 {
			for _, n := range r.Properties.Notes.RichText {
				if notes != "" {
					notes += "\n"
				}
				notes += n.PlainText
			}
		}
		t := Todo{
			Deadline: deadline,
			Course:   course,
			Name:     name,
			Type:     r.Properties.Type.Select.Name,
			Progress: Progress{
				Hans: r.Properties.HansProgress.Status.Name,
				Ira:  r.Properties.IraProgress.Status.Name,
			},
			DaysLeft: daysLeft,
			Notes:    notes,
		}

		out = append(out, t)
	}

	return out, nil
}
