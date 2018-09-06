package persistence

import (
	"testing"

	"errors"

	"github.com/Timothylock/inventory-management/config"
	"github.com/Timothylock/inventory-management/items"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	updateItem       = `UPDATE items.+`
	doesItemExist    = `SELECT count\(1\) FROM items.+`
	deleteItem       = `UPDATE items SET DELETED=1.+`
	addItem          = `INSERT INTO items`
	addItemOverwrite = `UPDATE items.+`
	searchItems      = `SELECT search.ID AS ID, NAME, CATEGORY, PICTURE_URL, DETAILS, LOCATION, USERNAME, QUANTITY, STATUS FROM.+`
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
				err := db.MoveItem("1234", tc.direction)
				assert.Error(t, err)
			} else {
				mock.ExpectQuery(doesItemExist).
					WithArgs("1234").
					WillReturnRows(rows)
				mock.ExpectExec(updateItem).
					WithArgs(tc.directionDB, "1234").
					WillReturnResult(sqlmock.NewResult(1234, 1))

				err := db.MoveItem("1234", tc.direction)
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

	err := db.MoveItem("1234", "in")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMoveItemInternalErr(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectQuery(doesItemExist).
		WithArgs("1234").
		WillReturnError(errors.New("sorry"))

	err := db.MoveItem("1234", "in")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteItemInternalErr(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectExec(deleteItem).
		WithArgs("1234").
		WillReturnError(errors.New("sorry"))

	err := db.DeleteItem("1234")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteItemNoRowsAff(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectExec(deleteItem).
		WithArgs("1234").
		WillReturnResult(sqlmock.NewResult(1, 0))

	err := db.DeleteItem("1234")
	assert.Equal(t, items.ItemNotFoundErr, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteItemSuccess(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	mock.ExpectExec(deleteItem).
		WithArgs("1234").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.DeleteItem("1234")
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
				WithArgs("foo").
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
		LastPerformedBy: "USERNAME",
		Quantity:        1,
		Status:          "checked in",
	}

	mock.ExpectQuery(doesItemExist).
		WithArgs("ID").
		WillReturnRows(rows)
	mock.ExpectExec(addItem).
		WithArgs("ID", "NAME", "CATEGORY", "PICTURE_URL", "DETAILS", "LOCATION", "USERNAME", 1, "checked in").
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
		LastPerformedBy: "USERNAME",
		Quantity:        1,
		Status:          "checked in",
	}

	mock.ExpectQuery(doesItemExist).
		WithArgs("ID").
		WillReturnRows(rows)
	mock.ExpectExec(addItemOverwrite).
		WithArgs("ID", "NAME", "CATEGORY", "PICTURE_URL", "DETAILS", "LOCATION", "USERNAME", 1, "ID").
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
		LastPerformedBy: "USERNAME",
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
