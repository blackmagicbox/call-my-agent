package profile

type CandidateProfile struct {
	Name        string       `json:"name"`
	Title       string       `json:"title"`
	Location    string       `json:"location"`
	Summary     string       `json:"summary"`
	Skills      []string     `json:"skills"`
	Experience  []Experience `json:"experience"`
	Education   []Education  `json:"education"`
	Languages   []Language   `json:"languages"`
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
	TargetRoles      []string `json:"target_roles"`
	TargetLevels     []string `json:"target_levels"`
	Locations        []string `json:"locations"`
	RemotePreference string   `json:"remote_preference"` // "required", "preferred", "not_interested"
	RedFlags         []string `json:"red_flags"`
	MinSalary        int      `json:"min_salary"`
	SalaryCurrency   string   `json:"salary_currency"` // "EUR", "USD", "GBP", etc.
}

// Proficiency uses CEFR levels (A1–C2) or "native"/"fluent" for natural language.
type Language struct {
	Language    string `json:"language"`
	Proficiency string `json:"proficiency"` // "native", "fluent", "C2", "C1", "B2", "B1", "A2", "A1"
}

type Education struct {
	Degree       string   `json:"degree"`      // "Bachelor's", "Master's", "PhD", "Associate's", "High School", etc.
	Field        string   `json:"field"`       // Major/field of study
	Institution  string   `json:"institution"` // School/university name
	Location     string   `json:"location"`    // City, Country
	StartDate    string   `json:"start_date"`
	EndDate      *string  `json:"end_date"`
	Achievements []string `json:"achievements,omitempty"` // Honors, awards, relevant coursework
}
