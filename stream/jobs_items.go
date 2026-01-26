package stream

import (
	"encoding/json"
	"fmt"
)

func (a *API) GetJobItems(id int64) ([]JobItem, error) {

	found, err := a.get(fmt.Sprintf("/jobs/%d/job_items", id))
	if err != nil {
		return nil, err
	}

	rawJobItems := make([]ParsedJobItem, 0)

	err = json.Unmarshal(found, &rawJobItems)
	if err != nil {
		return nil, err
	}

	jobItems := make([]JobItem, 0, len(rawJobItems))
	var itemUser int64

	for _, rji := range rawJobItems {

		jius, err := a.GetJobItemUsers(rji.ID)
		if err != nil {
			// TODO: Handle the error?!
			continue
		}

		itemUser = 0
		for _, ju := range jius {
			if ju.UserID == a.user {
				itemUser = ju.ID
				break
			}
		}

		if itemUser == 0 {
			continue
		}

		jobItems = append(jobItems, JobItem{
			Name:           rji.Name,
			Description:    rji.Description,
			ID:             rji.ID,
			User:           itemUser,
			PlannedMinutes: rji.PlannedMinutes,
			LoggedMinutes:  rji.LoggedMinutes,
		})

	}

	return jobItems, nil

}

func (a *API) GetJobItemUsers(id int64) ([]ParsedJobItemUser, error) {

	found, err := a.get(fmt.Sprintf("/job_items/%d/job_item_users", id))
	if err != nil {
		return nil, err
	}

	jobItemUsers := make([]ParsedJobItemUser, 0)

	err = json.Unmarshal(found, &jobItemUsers)
	if err != nil {
		return nil, err
	}

	return jobItemUsers, nil

}
