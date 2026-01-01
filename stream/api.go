package stream

import (
	"bytes"
	"io"
	"net/http"
)

type API struct {
	url, token string
	client     *http.Client
}

func NewAPI(url, token string) *API {
	return &API{
		url:    url,
		token:  token,
		client: &http.Client{},
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
