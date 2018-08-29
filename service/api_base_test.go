package service

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"bytes"

	"github.com/Timothylock/inventory-management/items"
)

func setupServer(ip items.Persister, t *testing.T) (*httptest.Server) {
	is := items.NewService(ip)

	serv := NewAPI(is)

	return httptest.NewServer(NewRouter(&serv))
}

func sendPost(url string, body interface{}) (*http.Response, error) {
	bs, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}
