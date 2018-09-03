package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Timothylock/inventory-management/config"
	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/upc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupServer(ip items.Persister, t *testing.T) *httptest.Server {
	cfg := config.Config{}

	is := items.NewService(ip)
	us := upc.NewService(cfg)

	serv := NewAPI(is, us)

	return httptest.NewServer(NewRouter(&serv, cfg))
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

func sendDelete(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func sendGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestParseBodyFail(t *testing.T) {
	testRequest := httptest.NewRequest(http.MethodPost, "/", errReader(0))
	assert.Error(t, parseBody(testRequest, nil))
}

func TestNotImplemented(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	ip := items.NewMockPersister(mc)
	server := setupServer(ip, t)
	defer server.Close()

	resp, err := sendPost(server.URL+"/api/user/add", "")
	assert.NoError(t, err)
	assert.Equal(t, "Not yet implemented\n", getBody(t, resp))
}
