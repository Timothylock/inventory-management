package service

import (
	"encoding/json"
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
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true}, nil).AnyTimes()
			},
			expectCode: 200,
		},
		{
			testName: "bad login",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: false}, nil).AnyTimes()
			},
			expectCode: 401,
		},
		{
			testName: "bad login",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true}, errors.New("error")).AnyTimes()
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

func TestGetUsers(t *testing.T) {
	type testCase struct {
		testName         string
		setMock          func(up *users.MockPersister)
		sendBody         LoginBody
		expectCode       int
		expectedResponse users.MultipleUsers
	}

	testCases := []testCase{
		{
			testName: "bad login",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: false}, nil).AnyTimes()
			},
			expectCode: 401,
		},
		{
			testName: "Not Sys Admin",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: false}, nil).AnyTimes()
			},
			expectCode: 401,
		},
		{
			testName: "Success",
			setMock: func(up *users.MockPersister) {
				u := users.MultipleUsers{
					{
						ID: 123,
					},
					{
						ID: 126,
					},
				}

				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: true}, nil).AnyTimes()
				up.EXPECT().GetUsers().Return(u, nil)
			},
			expectCode: 200,
			expectedResponse: users.MultipleUsers{
				{
					ID: 123,
				},
				{
					ID: 126,
				},
			},
		},
		{
			testName: "Success",
			setMock: func(up *users.MockPersister) {
				u := users.MultipleUsers{
					{
						ID: 123,
					},
					{
						ID: 126,
					},
				}

				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: true}, nil).AnyTimes()
				up.EXPECT().GetUsers().Return(u, errors.New("shoot"))
			},
			expectCode: 500,
			expectedResponse: users.MultipleUsers{
				{
					ID: 123,
				},
				{
					ID: 126,
				},
			},
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

			resp, err := sendGet(server.URL + "/api/users")
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

func TestAddUser(t *testing.T) {
	type testCase struct {
		testName   string
		setMock    func(up *users.MockPersister)
		sendBody   UserBody
		expectCode int
	}

	sb := UserBody{
		Username:   "someuser",
		Password:   "somepassword",
		Email:      "someEmail",
		IsSysAdmin: "true",
	}

	testCases := []testCase{
		{
			testName: "success",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: true}, nil).AnyTimes()
				up.EXPECT().AddUser("someuser", "someEmail", "somepassword", true, false).Return(nil)
			},
			sendBody:   sb,
			expectCode: 200,
		},
		{
			testName: "not sysadmin",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: false}, nil).AnyTimes()
			},
			sendBody:   sb,
			expectCode: 401,
		},
		{
			testName: "not logged in",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: false, IsSysAdmin: false}, nil).AnyTimes()
			},
			sendBody:   sb,
			expectCode: 401,
		},
		{
			testName: "internal error",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: true}, nil).AnyTimes()
				up.EXPECT().AddUser("someuser", "someEmail", "somepassword", true, false).Return(errors.New("sorry"))
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

			resp, err := sendPost(server.URL+"/api/user/add", tc.sendBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectCode, resp.StatusCode)
		})
	}
}

func TestAddUserBadBody(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	ip := items.NewMockPersister(mc)
	server := setupServerAuthenticated(ip, t)
	defer server.Close()

	resp, err := sendPost(server.URL+"/api/user/add", `{"username": 123}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteUser(t *testing.T) {
	type testCase struct {
		testName   string
		setMock    func(up *users.MockPersister)
		uHeader    string
		expectCode int
	}

	testCases := []testCase{
		{
			testName: "success",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: true, ID: 12345}, nil).AnyTimes()
				up.EXPECT().GetUserByUsername("someuser", 12345).Return(users.User{ID: 123, Valid: true}, nil)
				up.EXPECT().DeleteUser(123, 12345).Return(nil)
			},
			uHeader:    "u=someuser",
			expectCode: 200,
		},
		{
			testName: "not sys admin",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: false, ID: 12345}, nil).AnyTimes()
			},
			uHeader:    "u=someuser",
			expectCode: 401,
		},
		{
			testName: "missing user",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: true, ID: 12345}, nil).AnyTimes()
			},
			uHeader:    "",
			expectCode: 400,
		},
		{
			testName: "cannot get user from username",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: true, ID: 12345}, nil).AnyTimes()
				up.EXPECT().GetUserByUsername("someuser", 12345).Return(users.User{}, errors.New("some error"))
			},
			uHeader:    "u=someuser",
			expectCode: 500,
		},
		{
			testName: "user is system",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: true, ID: 12345}, nil).AnyTimes()
				up.EXPECT().GetUserByUsername("someuser", 12345).Return(users.User{ID: 0, Valid: true}, nil)
			},
			uHeader:    "u=someuser",
			expectCode: 500,
		},
		{
			testName: "error",
			setMock: func(up *users.MockPersister) {
				up.EXPECT().GetUserByToken(gomock.Any()).Return(users.User{Valid: true, IsSysAdmin: true, ID: 12345}, nil).AnyTimes()
				up.EXPECT().GetUserByUsername("someuser", 12345).Return(users.User{ID: 123, Valid: true}, nil)
				up.EXPECT().DeleteUser(123, 12345).Return(errors.New("some error"))
			},
			uHeader:    "u=someuser",
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

			resp, err := sendDelete(server.URL + "/api/user/delete?" + tc.uHeader)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectCode, resp.StatusCode)
		})
	}
}
