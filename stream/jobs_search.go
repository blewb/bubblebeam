package stream

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

//go:embed jobs_search.json
var jobsSearchJSON []byte

func (a *API) LoadJobs() error {

	a.status = StateFetching

	/*
		// Temporarily saved locally for speed
		found, err := a.post("/search?search_view=7&include_statistics=false", jobsSearchJSON)
		if err != nil {
			a.status = StateIdle
			return err
		}

		err = os.WriteFile("temp/jobs.json", found, 0644)
		return err
	*/

	found, err := os.ReadFile("temp/jobs.json")
	if err != nil {
		a.status = StateIdle
		return err
	}

	a.status = StateProcessing

	var search JobSearch
	err = json.Unmarshal(found, &search)
	if err != nil {
		a.status = StateIdle
		return err
	}

	a.jobs = make([]Job, 0, len(search.Results))

	for _, job := range search.Results {
		a.jobs = append(a.jobs, Job{
			ID:      job.ID,
			Name:    job.Name,
			Number:  job.Number,
			Company: job.Company.Name,
			Search:  strings.ToLower(fmt.Sprintf("%s|%s|%s", job.Number, job.Name, job.Company)),
		})
	}

	a.status = StateIdle
	return nil

}

/* Not needed?
func (a *API) GetJob(id string) {

	found, err := a.get("/jobs/" + id)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(found))

}
*/
