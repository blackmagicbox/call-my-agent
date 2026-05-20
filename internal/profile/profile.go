package profile

import "time"

type CandidateProfile struct {
	Name        string       `json:"name"`
	Title       string       `json:"title"`
	Location    string       `json:"location"`
	Summary     string       `json:"summary"`
	Skills      []string     `json:"skills"`
	Experience  []Experience `json:"experience"`
	Education   []Education  `json:"education"`
	Preferences Preferences  `json:"preferences"`
}

type Experience struct {
	Role     string   `json:"role"`
	Company  string   `json:"company"`
	Location string   `json:"location"`
	Start    string   `json:"start"`
	End      *string  `json:"end"`
	Bullets  []string `json:"bullets"`
}

type Preferences struct {
	TargetRoles  []string `json:"target_roles"`
	TargetLevels []string `json:"target_levels"`
	Locations    []string `json:"locations"`
	RedFlags     []string `json:"red_flags"`
	MinSalaryEUR int      `json:"min_salary_eur"`
}

type Education struct {
	Degree       string     `json:"degree"`      // "Bachelor's", "Master's", "PhD", "Associate's", "High School", etc.
	Field        string     `json:"field"`       // Major/field of study
	Institution  string     `json:"institution"` // School/university name
	Location     string     `json:"location"`    // City, Country
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"` // nil if currently enrolled
	GPA          *float64   `json:"gpa,omitempty"`
	Achievements []string   `json:"achievements,omitempty"` // Honors, awards, relevant coursework
}
