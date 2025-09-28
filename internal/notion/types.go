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
	DaysLeft int `json:"days_left"`
	Notes    string   `json:"notes"`
}


type DatabaseQueryResponse struct {
    Results []PageItem `json:"results"`
}

type PageItem struct {
    ID         string `json:"id"`
    URL        string `json:"url"`
    Properties struct {
        Deadline struct {
            Date struct {
                Start string `json:"start"`
                End   string `json:"end"`
            } `json:"date"`
        } `json:"deadline"`

        Name struct {
            Title []struct {
                PlainText string `json:"plain_text"`
            } `json:"title"`
        } `json:"name"`

        Type struct {
            Select struct {
                Name string `json:"name"`
            } `json:"select"`
        } `json:"type"`

        DaysLeft struct {
            Formula struct {
                String string `json:"string"`
            } `json:"formula"`
        } `json:"days left"`

        HansProgress struct {
            Status struct {
                Name string `json:"name"`
            } `json:"status"`
        } `json:"Hans Progress"`

        IraProgress struct {
            Status struct {
                Name string `json:"name"`
            } `json:"status"`
        } `json:"Ira Progress"`

        Notes struct {
            RichText []struct {
                PlainText string `json:"plain_text"`
            } `json:"rich_text"`
        } `json:"notes"`

        Course struct {
            Relation []struct {
                ID string `json:"id"`
            } `json:"relation"`
        } `json:"course"`
    } `json:"properties"`
}

