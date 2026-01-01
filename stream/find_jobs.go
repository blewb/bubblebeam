package stream

import (
	"fmt"
	"net/http"
)

type SearchParams struct {
	WildcardSearch        string                `json:"wildcardSearch"`
	Offset                int                   `json:"offset"`
	MaxResults            int                   `json:"maxResults"`
	SortField             int                   `json:"sortField"`
	SortAscending         bool                  `json:"sortAscending"`
	FilterGroupCollection FilterGroupCollection `json:"filterGroupCollection"`
}
type Filters struct {
	ValueMatchTypeID int    `json:"valueMatchTypeId"`
	Value            string `json:"value"`
}
type FilterGroups struct {
	FilterGroupTypeID    int       `json:"filterGroupTypeId"`
	ConditionMatchTypeID int       `json:"conditionMatchTypeId"`
	Filters              []Filters `json:"filters"`
}
type FilterGroupCollection struct {
	ConditionMatchTypeID   int            `json:"conditionMatchTypeId"`
	FilterGroupCollections []any          `json:"filterGroupCollections"`
	FilterGroups           []FilterGroups `json:"filterGroups"`
}

func (a *API) FindJobs() {

	found, err := a.query(http.MethodPost, "/search?search_view=7")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(found))

}

func (a *API) GetBranches() {

	found, err := a.query(http.MethodGet, "/users")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(found))

}
