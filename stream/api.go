package stream

import (
	"bytes"
	"io"
	"net/http"
)

type APIState int

const (
	StateIdle APIState = iota
	StateFetching
	StateProcessing
	StateUpdating
)

type API struct {
	url, token string
	user       int64
	client     *http.Client
	status     APIState
	jobs       []Job
	jobItems   map[int64][]JobItem
}

func NewAPI(url, token string, user int64) *API {
	return &API{
		url:      url,
		token:    token,
		user:     user,
		client:   &http.Client{},
		status:   StateIdle,
		jobs:     make([]Job, 0),
		jobItems: make(map[int64][]JobItem),
	}
}

func (a *API) get(path string) ([]byte, error) {
	return a.query(http.MethodGet, path, nil)
}

func (a *API) post(path string, content []byte) ([]byte, error) {
	return a.query(http.MethodPost, path, content)
}

// Is it that easy...?
func (a *API) put(path string, content []byte) ([]byte, error) {
	return a.query(http.MethodPut, path, content)
}

func (a *API) query(requestType, path string, content []byte) ([]byte, error) {

	var (
		req *http.Request
		err error
	)

	reqUrl := a.url + path

	if len(content) > 0 {
		reqBody := bytes.NewReader(content)
		req, err = http.NewRequest(requestType, reqUrl, reqBody)
	} else {
		req, err = http.NewRequest(requestType, reqUrl, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+a.token)
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}
