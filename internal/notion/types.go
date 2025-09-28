package notion

import "time"

type Status string

const (
	StatusNotStarted Status = "not started"
	StatusInProgress Status = "in progress"
	StatusDone       Status = "done"
	StatusNA         Status = "N/A"
)

type Progress struct {
	Hans string `json:"hans_progress"`
	Ira  string `json:"ira_progress"`
}

type Todo struct {
	Deadline time.Time `json:"deadline"`  
	Course   string   `json:"course"`    
	Name     string   `json:"name"`      
	Type     string   `json:"type"`     
	Progress Progress `json:"progress"`  
	DaysLeft string   `json:"days_left"`
	Notes    string   `json:"notes"`
}
