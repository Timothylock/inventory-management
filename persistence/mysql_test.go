package persistence

import (
	"encoding/json"
	"testing"

	"errors"

	"github.com/Timothylock/inventory-management/config"
	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/users"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	updateItem       = `UPDATE items.+`
	doesItemExist    = `SELECT count\(1\) FROM items.+`
	deleteItem       = `UPDATE items SET DELETED=1.+`
	GetUser          = `SELECT ID\, ISSYSADMIN\, EMAIL\, TOKEN\, USERNAME FROM users.+`
	addItem          = `INSERT INTO items`
	addUser          = `INSERT INTO users.+`
	addItemOverwrite = `UPDATE items.+`
	searchItems      = `SELECT search.ID AS ID, NAME, CATEGORY, PICTURE_URL, DETAILS, LOCATION, USERNAME, QUANTITY, STATUS FROM.+`
	deleteUser       = `UPDATE users SET ACTIVE=0.+`
)

func newTestDB(t *testing.T) (*MySQL, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	return &MySQL{
		conn: sqlx.NewDb(db, "sqlmock"),
	}, mock
}

func TestConnFailsNoneExistent(t *testing.T) {
	c := config.Config{
		DbUrl:  "localhost",
		DbUser: "nobody",
		DbPass: "foo",
		DbName: "blah",
	}

	_, err := NewMySQL(&c)
	assert.Error(t, err)
}

func TestMoveItem(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	type testCase struct {
		testName    string
		direction   string
		directionDB string
	}

	testCases := []testCase{
		{
			testName:    "check in item",
			direction:   "in",
			directionDB: "checked in",
		},
		{
			testName:    "check out item",
			direction:   "out",
			directionDB: "checked out",
		},
		{
			testName:  "invalid direction",
			direction: "outtt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"COUNT(1)"})
			rows.AddRow(1)

			if tc.direction != "in" && tc.direction != "out" {
				err := db.MoveItem("1234", tc.direction, 123)
				assert.Error(t, err)
			} else {
				mock.ExpectQuery(doesItemExist).
					WithArgs("1234").
					WillReturnRows(rows)
				mock.ExpectExec(updateItem).
					WithArgs(tc.directionDB, 123, "1234").
					WillReturnResult(sqlmock.NewResult(1234, 1))

				err := db.MoveItem("1234", tc.direction, 123)
				assert.NoError(t, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}

func TestMoveItemNotFound(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"COUNT(1)"})
	rows.AddRow(0)

	mock.ExpectQuery(doesItemExist).
		WithArgs("1234").
		WillReturnRows(rows)

	err := db.MoveItem("1234", "in", 123)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMoveItemInternalErr(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectQuery(doesItemExist).
		WithArgs("1234").
		WillReturnError(errors.New("sorry"))

	err := db.MoveItem("1234", "in", 123)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteItemInternalErr(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectExec(deleteItem).
		WithArgs(123, "1234").
		WillReturnError(errors.New("sorry"))

	err := db.DeleteItem("1234", 123)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteItemNoRowsAff(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectExec(deleteItem).
		WithArgs(123, "1234").
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := db.DeleteItem("1234", 123)
	assert.Equal(t, items.ItemNotFoundErr, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteItemSuccess(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectExec(deleteItem).
		WithArgs(123, "1234").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.DeleteItem("1234", 123)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchItems(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	type testCase struct {
		testName string
		addRows  func(rows *sqlmock.Rows)
		expected items.ItemDetailList
	}

	testCases := []testCase{
		{
			testName: "find 1 item",
			addRows: func(rows *sqlmock.Rows) {
				rows.AddRow("1", "foo", "fi", "bar", "fum", "bah", "humbug", 1, "checked in")
			},
			expected: items.ItemDetailList{
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
			},
		},
		{
			testName: "find multiple items",
			addRows: func(rows *sqlmock.Rows) {
				rows.AddRow("1", "foo", "fi", "bar", "fum", "bah", "humbug", 1, "checked in")
				rows.AddRow("2", "foo", "fi", "bar", "fum", "bah", "humbug", 1, "checked in")
			},
			expected: items.ItemDetailList{
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
			testName: "find no items",
			addRows:  func(rows *sqlmock.Rows) {},
			expected: items.ItemDetailList{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"ID", "NAME", "CATEGORY", "PICTURE_URL", "DETAILS", "LOCATION", "USERNAME", "QUANTITY", "STATUS"})
			tc.addRows(rows)

			mock.ExpectQuery(searchItems).
				WithArgs("foo", "foo").
				WillReturnRows(rows)

			r, err := db.SearchItems("foo")
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
			assert.Equal(t, tc.expected, r)
		})
	}
}

func TestAddItemAlreadyExists(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"COUNT(1)"})
	rows.AddRow(1)

	mock.ExpectQuery(doesItemExist).
		WithArgs("1234").
		WillReturnRows(rows)

	err := db.AddItem(items.ItemDetail{ID: "1234"}, false)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddItemAlreadyExistsError(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"COUNT(1)"})
	rows.AddRow(1)

	mock.ExpectQuery(doesItemExist).
		WithArgs("1234").
		WillReturnError(errors.New("some error"))

	err := db.AddItem(items.ItemDetail{ID: "1234"}, false)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddItemSuccess(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"COUNT(1)"})
	rows.AddRow(0)

	item := items.ItemDetail{
		ID:              "ID",
		Name:            "NAME",
		Category:        "CATEGORY",
		PictureURL:      "PICTURE_URL",
		Details:         "DETAILS",
		Location:        "LOCATION",
		LastPerformedBy: "1",
		Quantity:        1,
		Status:          "checked in",
	}

	mock.ExpectQuery(doesItemExist).
		WithArgs("ID").
		WillReturnRows(rows)
	mock.ExpectExec(addItem).
		WithArgs("ID", "NAME", "CATEGORY", "PICTURE_URL", "DETAILS", "LOCATION", "1", 1, "checked in").
		WillReturnResult(sqlmock.NewResult(123, 1))

	err := db.AddItem(item, false)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddItemOverwriteSuccess(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"COUNT(1)"})
	rows.AddRow(1)

	item := items.ItemDetail{
		ID:              "ID",
		Name:            "NAME",
		Category:        "CATEGORY",
		PictureURL:      "PICTURE_URL",
		Details:         "DETAILS",
		Location:        "LOCATION",
		LastPerformedBy: "1",
		Quantity:        1,
		Status:          "checked in",
	}

	mock.ExpectQuery(doesItemExist).
		WithArgs("ID").
		WillReturnRows(rows)
	mock.ExpectExec(addItemOverwrite).
		WithArgs("ID", "NAME", "CATEGORY", "PICTURE_URL", "DETAILS", "LOCATION", "1", 1, "ID").
		WillReturnResult(sqlmock.NewResult(123, 1))

	err := db.AddItem(item, true)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddItemOverwriteFailure(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"COUNT(1)"})
	rows.AddRow(0)

	item := items.ItemDetail{
		ID:              "ID",
		Name:            "NAME",
		Category:        "CATEGORY",
		PictureURL:      "PICTURE_URL",
		Details:         "DETAILS",
		Location:        "LOCATION",
		LastPerformedBy: "1",
		Quantity:        1,
		Status:          "checked in",
	}

	mock.ExpectQuery(doesItemExist).
		WithArgs("ID").
		WillReturnRows(rows)

	err := db.AddItem(item, true)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIsValidToken(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"ID", "ISSYSADMIN", "EMAIL", "TOKEN", "USERNAME"})
	rows.AddRow(0, 1, "foo@bar.com", "someToken", "someUser")

	mock.ExpectQuery(GetUser).
		WithArgs("foo").
		WillReturnRows(rows)

	user, err := db.GetUserByToken("foo")
	assert.True(t, user.Valid)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIsNotValidToken(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"ID", "ISSYSADMIN", "EMAIL", "TOKEN", "USERNAME"})

	mock.ExpectQuery(GetUser).
		WithArgs("foo").
		WillReturnRows(rows)

	user, err := db.GetUserByToken("foo")
	assert.False(t, user.Valid)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIsValidTokenFail(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectQuery(GetUser).
		WithArgs("foo").
		WillReturnError(errors.New("sorry"))

	_, err := db.GetUserByToken("foo")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserSuccess(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"ID", "ISSYSADMIN", "EMAIL", "TOKEN", "USERNAME"})
	rows.AddRow(123, 1, "foo@bar.com", "someToken", "someUser")

	mock.ExpectQuery(GetUser).
		WithArgs("someUser", "nU4eI71bcnBGqeO0t9tXvY1u5oQ=").
		WillReturnRows(rows)

	expectedUser := users.User{
		Valid:      true,
		ID:         123,
		Token:      "someToken",
		Username:   "someUser",
		IsSysAdmin: true,
		Email:      "foo@bar.com",
	}
	expectedUserJson, err := json.Marshal(expectedUser)
	assert.NoError(t, err)

	u, err := db.GetUser("someUser", "pass")
	assert.NoError(t, err)
	uJson, err := json.Marshal(u)

	assert.JSONEq(t, string(expectedUserJson), string(uJson))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserNotFound(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"ID", "ISSYSADMIN", "EMAIL", "TOKEN"})

	mock.ExpectQuery(GetUser).
		WithArgs("user", "nU4eI71bcnBGqeO0t9tXvY1u5oQ=").
		WillReturnRows(rows)

	expectedUser := users.User{
		Valid: false,
	}
	expectedUserJson, err := json.Marshal(expectedUser)
	assert.NoError(t, err)

	u, err := db.GetUser("user", "pass")
	assert.NoError(t, err)
	uJson, err := json.Marshal(u)

	assert.JSONEq(t, string(expectedUserJson), string(uJson))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserErr(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectQuery(GetUser).
		WithArgs("user", "nU4eI71bcnBGqeO0t9tXvY1u5oQ=").
		WillReturnError(errors.New("error"))

	_, err := db.GetUser("user", "pass")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsernameSuccess(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"ID", "ISSYSADMIN", "EMAIL", "TOKEN", "USERNAME"})
	rows.AddRow(123, 1, "foo@bar.com", "someToken", "someUser")

	mock.ExpectQuery(GetUser).
		WithArgs("someUser").
		WillReturnRows(rows)

	expectedUser := users.User{
		Valid:      true,
		ID:         123,
		Token:      "someToken",
		Username:   "someUser",
		IsSysAdmin: true,
		Email:      "foo@bar.com",
	}
	expectedUserJson, err := json.Marshal(expectedUser)
	assert.NoError(t, err)

	u, err := db.GetUserByUsername("someUser", 123)
	assert.NoError(t, err)
	uJson, err := json.Marshal(u)

	assert.JSONEq(t, string(expectedUserJson), string(uJson))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsernameNotFound(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"ID", "ISSYSADMIN", "EMAIL", "TOKEN"})

	mock.ExpectQuery(GetUser).
		WithArgs("user").
		WillReturnRows(rows)

	expectedUser := users.User{
		Valid: false,
	}
	expectedUserJson, err := json.Marshal(expectedUser)
	assert.NoError(t, err)

	u, err := db.GetUserByUsername("user", 123)
	assert.NoError(t, err)
	uJson, err := json.Marshal(u)

	assert.JSONEq(t, string(expectedUserJson), string(uJson))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsernameErr(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectQuery(GetUser).
		WithArgs("user").
		WillReturnError(errors.New("error"))

	_, err := db.GetUserByUsername("user", 123)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUsersSuccess(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	rows := sqlmock.NewRows([]string{"ID", "ISSYSADMIN", "EMAIL", "TOKEN", "USERNAME"})
	rows.AddRow(123, 1, "foo@bar.com", "someToken", "someUser")
	rows.AddRow(124, 0, "foo2@bar.com", "someToken2", "someUser2")

	mock.ExpectQuery(GetUser).
		WillReturnRows(rows)

	expectedUsers := users.MultipleUsers{
		{
			Valid:      true,
			ID:         123,
			Token:      "someToken",
			Username:   "someUser",
			IsSysAdmin: true,
			Email:      "foo@bar.com",
		},
		{
			Valid:      true,
			ID:         124,
			Token:      "someToken2",
			Username:   "someUser2",
			IsSysAdmin: false,
			Email:      "foo2@bar.com",
		},
	}
	expectedUserJson, err := json.Marshal(expectedUsers)
	assert.NoError(t, err)

	u, err := db.GetUsers()
	assert.NoError(t, err)
	uJson, err := json.Marshal(u)

	assert.JSONEq(t, string(expectedUserJson), string(uJson))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUsersErr(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectQuery(GetUser).
		WillReturnError(errors.New("error"))

	_, err := db.GetUsers()
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddUserSuccess(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectExec(addUser).
		WithArgs("user", "email", "W6ph5Mm5Pz8GgiULbPgzG37mj9g=", sqlmock.AnyArg(), true).
		WillReturnResult(sqlmock.NewResult(123, 1))

	err := db.AddUser("user", "email", "password", true, false)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUserInternalErr(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectExec(deleteUser).
		WithArgs(1234).
		WillReturnError(errors.New("sorry"))

	err := db.DeleteUser(1234, 123)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUserNoRowsAff(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectExec(deleteUser).
		WithArgs(1234).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := db.DeleteUser(1234, 123)
	assert.Equal(t, items.ItemNotFoundErr, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUserSuccess(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectExec(deleteUser).
		WithArgs(1234).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.DeleteUser(1234, 123)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
