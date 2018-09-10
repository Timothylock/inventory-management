package service

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/users"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	type testCase struct {
		testName   string
		setMock    func(up *users.MockPersister)
		sendBody   LoginBody
		expectCode int
	}

	sb := LoginBody{
		Username: "someuser",
		Password: "somepassword",
	}

	testCases := []testCase{
		{
			testName: "success",
			setMock: func(up *users.MockPersister) {
				u := users.User{
					Valid: true,
				}
				up.EXPECT().GetUser("someuser", "somepassword").Return(u, nil)
			},
			sendBody:   sb,
			expectCode: 200,
		},
		{
			testName: "bad login",
			setMock: func(up *users.MockPersister) {
				u := users.User{
					Valid: false,
				}
				up.EXPECT().GetUser("someuser", "somepassword").Return(u, nil)
			},
			sendBody:   sb,
			expectCode: 401,
		},
		{
			testName: "error",
			setMock: func(up *users.MockPersister) {
				u := users.User{
					Valid: false,
				}
				up.EXPECT().GetUser("someuser", "somepassword").Return(u, errors.New("error"))
			},
			sendBody:   sb,
			expectCode: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			up := users.NewMockPersister(mc)
			tc.setMock(up)

			server := setupServer(nil, up, t)
			defer server.Close()

			resp, err := sendPost(server.URL+"/api/user/login", tc.sendBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectCode, resp.StatusCode)
		})
	}
}

func TestLoginBadBody(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	ip := items.NewMockPersister(mc)
	server := setupServerAuthenticated(ip, t)
	defer server.Close()

	resp, err := sendPost(server.URL+"/api/user/login", `{"username": 123}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestLoginCheck(t *testing.T) {
	type testCase struct {
		testName   string
		setMock    func(up *users.MockPersister)
		sendBody   LoginBody
		expectCode int
	}

	testCases := []testCase{
		{
			testName: "bad login",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().IsValidToken(gomock.Any()).Return(true, nil)
			},
			expectCode: 200,
		},
		{
			testName: "bad login",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().IsValidToken(gomock.Any()).Return(false, nil)
			},
			expectCode: 401,
		},
		{
			testName: "bad login",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().IsValidToken(gomock.Any()).Return(true, errors.New("error"))
			},
			expectCode: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			up := users.NewMockPersister(mc)
			tc.setMock(up)

			server := setupServer(nil, up, t)
			defer server.Close()

			resp, err := sendGet(server.URL + "/api/user/logincheck")
			assert.NoError(t, err)
			assert.Equal(t, tc.expectCode, resp.StatusCode)
		})
	}
}
