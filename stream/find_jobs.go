package stream

import (
	_ "embed"
	"fmt"
)

//go:embed find_jobs.json
var findjobsJSON []byte

func (a *API) FindJobs() {

	found, err := a.post("/search?search_view=7&include_statistics=false", findjobsJSON)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(found))

}
