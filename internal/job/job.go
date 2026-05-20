package job

import "time"

type Job struct {
	ID             string        `json:"id"`
	Title          string        `json:"title"`
	Company        string        `json:"company"`
	Location       string        `json:"location"`
	Remote         string        `json:"remote"` // "remote", "hybrid", "onsite"
	Seniority      string        `json:"seniority"`
	Description    string        `json:"description"`
	RequiredSkills []string      `json:"required_skills"`
	NiceToHave     []string      `json:"nice_to_have"`
	Compensation   *Compensation `json:"compensation"`
	URL            string        `json:"url"`
	SeenAt         time.Time     `json:"seen_at"`
	Status         string        `json:"status"` // "to_evaluate", "to_apply", "applied", "rejected", "archived"
}

type Compensation struct {
	Min      int    `json:"min"`
	Max      int    `json:"max"`
	Currency string `json:"currency"`
}
