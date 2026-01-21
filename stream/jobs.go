package stream

// Parsed items are pulled from JSON API
// and flattened into simple structs
type ParsedJob struct {
	Number  string `json:"number"`
	Name    string `json:"name"`
	Company struct {
		Name string `json:"name"`
	} `json:"company"`
	ID int64 `json:"id"`
}

type Job struct {
	Number  string
	Name    string
	Company string
	Search  string
	ID      int64
}

type JobSearch struct {
	Results []ParsedJob `json:"searchResults"`
}

type ParsedJobItem struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	ID             int64  `json:"id"`
	PlannedMinutes int    `json:"totalPlannedMinutes"`
	LoggedMinutes  int    `json:"totalLoggedMinutes"`
}

type JobItem struct {
	Name           string
	Description    string
	ID             int64
	User           int64 // This is a unique ID representing the bridge between an actual user and a job item
	PlannedMinutes int
	LoggedMinutes  int
}

func (a API) GetJobs() []Job {
	return a.jobs
}
