package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/responses"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func getBody(t *testing.T, resp *http.Response) string {
	if resp == nil || resp.Body == nil {
		return ""
	}

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "failed reading body")

	return string(b)
}

func TestMoveItem(t *testing.T) {
	type testCase struct {
		testName         string
		setMock          func(*items.MockPersister)
		expectCode       int
		expectedResponse responses.Success
		body             MoveBody
	}

	testCases := []testCase{
		{
			testName: "success",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().MoveItem("1234", "in", 123).Return(nil)
			},
			expectCode:       200,
			expectedResponse: responses.Success{Success: true},
			body: MoveBody{
				ID:        "1234",
				Direction: "in",
			},
		},
		{
			testName: "internal error",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().MoveItem("1234", "in", 123).Return(errors.New("sorry"))
			},
			expectCode: 500,
			body: MoveBody{
				ID:        "1234",
				Direction: "in",
			},
		},
		{
			testName: "item not found error",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().MoveItem("1234", "in", 123).Return(items.ItemNotFoundErr)
			},
			expectCode: 404,
			body: MoveBody{
				ID:        "1234",
				Direction: "in",
			},
		},
		{
			testName:   "missing id in body",
			setMock:    func(ip *items.MockPersister) {},
			expectCode: 400,
			body: MoveBody{
				ID:        "",
				Direction: "in",
			},
		},
		{
			testName:   "missing direction in body",
			setMock:    func(ip *items.MockPersister) {},
			expectCode: 400,
			body: MoveBody{
				ID:        "1234",
				Direction: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			ip := items.NewMockPersister(mc)
			tc.setMock(ip)

			server := setupServerAuthenticated(ip, t)
			defer server.Close()

			resp, err := sendPost(server.URL+"/api/item/move", tc.body)
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

func TestMoveItemBadBody(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	ip := items.NewMockPersister(mc)
	server := setupServerAuthenticated(ip, t)
	defer server.Close()

	resp, err := sendPost(server.URL+"/api/item/move", `{"direction": "in", "id": 123}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteItem(t *testing.T) {
	type testCase struct {
		testName         string
		setMock          func(*items.MockPersister)
		expectCode       int
		expectedResponse responses.Success
		id               string
	}

	testCases := []testCase{
		{
			testName: "success",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().DeleteItem("1", 123).Return(nil)
			},
			expectCode:       200,
			expectedResponse: responses.Success{Success: true},
			id:               "1",
		},
		{
			testName:   "missing id",
			setMock:    func(ip *items.MockPersister) {},
			expectCode: 400,
			id:         "",
		},
		{
			testName: "internal error",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().DeleteItem("1", 123).Return(errors.New("oops"))
			},
			expectCode: 500,
			id:         "1",
		},
		{
			testName: "item not found",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().DeleteItem("1", 123).Return(items.ItemNotFoundErr)
			},
			expectCode: 404,
			id:         "1",
		},
		{
			testName:   "too many param",
			setMock:    func(ip *items.MockPersister) {},
			expectCode: 400,
			id:         "1&id=12",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			ip := items.NewMockPersister(mc)
			tc.setMock(ip)

			server := setupServerAuthenticated(ip, t)
			defer server.Close()

			resp, err := sendDelete(server.URL + "/api/item?id=" + tc.id)
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

func TestSearchItems(t *testing.T) {
	type testCase struct {
		testName         string
		setMock          func(*items.MockPersister)
		expectCode       int
		expectedResponse items.ItemDetailList
	}

	testCases := []testCase{
		{
			testName: "success",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().SearchItems("foo").Return(items.ItemDetailList{
					{
						ID:              "1",
						Name:            "foo",
						Category:        "fi",
						PictureURL:      "bar",
						Details:         "fum",
						Location:        "bah",
						LastPerformedBy: "humbug",
						Quantity:        1,
						Status:          "checked in",
					},
					{
						ID:              "2",
						Name:            "foo",
						Category:        "fi",
						PictureURL:      "bar",
						Details:         "fum",
						Location:        "bah",
						LastPerformedBy: "humbug",
						Quantity:        1,
						Status:          "checked in",
					},
				}, nil)
			},
			expectCode: 200,
			expectedResponse: items.ItemDetailList{
				{
					ID:              "1",
					Name:            "foo",
					Category:        "fi",
					PictureURL:      "bar",
					Details:         "fum",
					Location:        "bah",
					LastPerformedBy: "humbug",
					Quantity:        1,
					Status:          "checked in",
				},
				{
					ID:              "2",
					Name:            "foo",
					Category:        "fi",
					PictureURL:      "bar",
					Details:         "fum",
					Location:        "bah",
					LastPerformedBy: "humbug",
					Quantity:        1,
					Status:          "checked in",
				},
			},
		},
		{
			testName: "error",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().SearchItems("foo").Return(nil, errors.New("some error"))
			},
			expectCode: 500,
		},
		{
			testName: "error",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().SearchItems("foo").Return(nil, errors.New("some error"))
			},
			expectCode: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			ip := items.NewMockPersister(mc)
			tc.setMock(ip)

			server := setupServerAuthenticated(ip, t)
			defer server.Close()

			resp, err := sendGet(server.URL + "/api/item/info?q=foo")
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

func TestSearchItemsMissingQuery(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	ip := items.NewMockPersister(mc)
	server := setupServerAuthenticated(ip, t)
	defer server.Close()

	resp, err := sendGet(server.URL + "/api/item/info")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAddItem(t *testing.T) {
	type testCase struct {
		testName   string
		setMock    func(*items.MockPersister)
		sendBody   AddBody
		expectCode int
	}

	testCases := []testCase{
		{
			testName: "success",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().AddItem(items.ItemDetail{
					ID:              "1",
					Name:            "foo",
					Category:        "fi",
					PictureURL:      "bar",
					Details:         "fum",
					Location:        "bah",
					Quantity:        1,
					LastPerformedBy: "123",
					Status:          "checked in",
				}, false).Return(nil)
			},
			sendBody: AddBody{
				ID:         "1",
				Name:       "foo",
				Category:   "fi",
				PictureURL: "bar",
				Details:    "fum",
				Location:   "bah",
				Quantity:   1,
			},
			expectCode: 200,
		},
		{
			testName: "Missing in body",
			setMock:  func(ip *items.MockPersister) {},
			sendBody: AddBody{
				ID:         "1",
				Name:       "",
				Category:   "",
				PictureURL: "bar",
				Details:    "fum",
				Location:   "bah",
				Quantity:   1,
			},
			expectCode: 400,
		},
		{
			testName: "item already exists",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().AddItem(items.ItemDetail{
					ID:              "1",
					Name:            "foo",
					Category:        "fi",
					PictureURL:      "bar",
					Details:         "fum",
					Location:        "bah",
					Quantity:        1,
					LastPerformedBy: "123",
					Status:          "checked in",
				}, false).Return(items.ItemAlreadyExistsErr)
			},
			sendBody: AddBody{
				ID:         "1",
				Name:       "foo",
				Category:   "fi",
				PictureURL: "bar",
				Details:    "fum",
				Location:   "bah",
				Quantity:   1,
			},
			expectCode: 400,
		},
		{
			testName: "internal error",
			setMock: func(ip *items.MockPersister) {
				ip.EXPECT().AddItem(items.ItemDetail{
					ID:              "1",
					Name:            "foo",
					Category:        "fi",
					PictureURL:      "bar",
					Details:         "fum",
					Location:        "bah",
					Quantity:        1,
					LastPerformedBy: "123",
					Status:          "checked in",
				}, false).Return(errors.New("sorry"))
			},
			sendBody: AddBody{
				ID:         "1",
				Name:       "foo",
				Category:   "fi",
				PictureURL: "bar",
				Details:    "fum",
				Location:   "bah",
				Quantity:   1,
			},
			expectCode: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			ip := items.NewMockPersister(mc)
			tc.setMock(ip)

			server := setupServerAuthenticated(ip, t)
			defer server.Close()

			resp, err := sendPost(server.URL+"/api/item", tc.sendBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectCode, resp.StatusCode)

			if tc.expectCode == 200 {
				//b, err := json.Marshal(tc.expectedResponse)
				assert.NoError(t, err)
				//assert.JSONEq(t, string(b), string(getBody(t, resp)))
			}
		})
	}
}

func TestAddItemBadBody(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	ip := items.NewMockPersister(mc)
	server := setupServerAuthenticated(ip, t)
	defer server.Close()

	resp, err := sendPost(server.URL+"/api/item", `{"id": 123}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
