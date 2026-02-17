package main

import (
	"strings"

	"github.com/blewb/bubblebeam/stream"
)

func searchJobs(api *stream.API, term string) []stream.Job {

	term = strings.ToLower(strings.TrimSpace(term))

	if len(term) < 2 {
		return nil
	}

	var found []stream.Job

	for _, job := range api.GetJobs() {
		if strings.Contains(job.Search, term) {
			found = append(found, job)
		}
	}

	return found

}
