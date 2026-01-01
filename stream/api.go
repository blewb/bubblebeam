package stream

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

type API struct {
	url, token string
	client     *http.Client
}

func NewAPI(url, token string) *API {
	tr := &http.Transport{
		// This empty map trick prevents the package from
		// initializing the HTTP/2 protocol upgrade
		TLSNextProto:       make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	return &API{
		url:    url,
		token:  token,
		client: client,
	}
}

func (a *API) query(requestType, path string, value ...any) ([]byte, error) {

	var (
		req *http.Request
		err error
	)

	reqUrl := a.url + path

	if len(value) > 0 {
		js, jerr := json.Marshal(value[0])
		if jerr != nil {
			return nil, jerr
		}
		reqBody := bytes.NewReader(js)
		req, err = http.NewRequest(requestType, reqUrl, reqBody)
	} else {
		req, err = http.NewRequest(requestType, reqUrl, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+a.token)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "curl/7.79.1")

	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Println("")
	fmt.Printf("%s", dump)
	fmt.Println("")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dump, _ = httputil.DumpResponse(resp, true)
	fmt.Println("")
	fmt.Printf("%s", dump)
	fmt.Println("")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}
