package main

import (
	"strings"

	"github.com/blewb/bubblebeam/stream"
)

func (m *model) SearchJobs() {

	found := make([]stream.Job, 0)
	term := strings.ToLower(strings.TrimSpace(m.searchInput.Value()))

	if len(term) < 2 {
		m.searchJobs = found
		return
	}

	for _, job := range m.api.GetJobs() {

		if strings.Contains(job.Search, term) {
			found = append(found, job)
		}

	}

	m.searchJobs = found

}
