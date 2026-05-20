package profile

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
type Education struct{}

type Experience struct{}

type Preferences struct{}
