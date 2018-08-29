package main

import (
	"testing"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	"inventory-management/items"
	"inventory-management/service"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"bytes"
	"net/http"
	"io/ioutil"
	"inventory-management/responses"
	"errors"
)

func setupServer(ip items.Persister, t *testing.T) (*httptest.Server) {
	is := items.NewService(ip)

	serv := service.NewAPI(is)

	return httptest.NewServer(service.NewRouter(&serv))
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


func getBody(t *testing.T, resp *http.Response) string {
	if resp == nil || resp.Body == nil {
		return ""
	}

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "failed reading body")

	return string(b)
}

type MoveBody struct {
	Direction string `json:"direction"`
	ID        string `json:"id"`
}

func TestMoveItem(t *testing.T) {
	type testCase struct {
		testName     string
		setMock      func(*items.MockPersister)
		expectCode int
		expectedResponse responses.Success
	}

	testCases := []testCase{
		{
			testName: "success",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().MoveItem("1234", "in").Return(nil)
			},
			expectCode: 200,
			expectedResponse: responses.Success{Success:true},
		},
		{
			testName: "internal error",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().MoveItem("1234", "in").Return(errors.New("sorry"))
			},
			expectCode: 500,
		},
		{
			testName: "item not found error",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().MoveItem("1234", "in").Return(items.ItemNotFoundErr)
			},
			expectCode: 404,
		},
	}

	mb := MoveBody{
		ID: "1234",
		Direction: "in",
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			ip := items.NewMockPersister(mc)
			tc.setMock(ip)

			server := setupServer(ip, t)
			defer server.Close()

			resp, err := sendPost(server.URL+"/api/item/move", mb)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectCode, resp.StatusCode)

			if tc.expectCode == 200 {
				b, err := json.Marshal(tc.expectedResponse)
				assert.NoError(t, err)
				assert.JSONEq(t, string(b), string(getBody(t, resp)))
			}
		})
	}
}

