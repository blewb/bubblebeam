package stream

import (
	"encoding/json"
	"fmt"
)

func (a *API) GetJobItems(id string) {

	found, err := a.get(fmt.Sprintf("/jobs/%s/job_items", id))
	if err != nil {
		fmt.Println(err)
	}

	jobItems := make([]ParsedJobItem, 0)

	err = json.Unmarshal(found, &jobItems)
	if err != nil {
		fmt.Println(err)
	}

	for _, itm := range jobItems {
		fmt.Println(itm)
	}

}

func (a *API) GetJobItemUsers(id string) {

	found, err := a.get(fmt.Sprintf("/job_items/%s/job_item_users", id))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(found))

}
