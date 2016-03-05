package twitch

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type apiSegment struct {
	parent *apiSegment
	path   string
}

func (s *apiSegment) root() *api {
	if s.parent != nil {
		return s.parent.root()
	}
	return s.(*api)
}

func (s *apiSegment) fetch(v interface{}, method string, path string, params url.Values) error {
	url := s.path + path
	if s.parent != nil {
		return s.parent.fetch(v, method, url, params)
	}

	data := bytes.NewBufferString(params.Encode())
	req := http.NewRequest(method, url, data)
	resp, err := s.(*api).Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
