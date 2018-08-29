package persistence

import (
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"errors"
	"github.com/Timothylock/inventory-management/config"
)

const (
	updateItem = `UPDATE items`
	doesItemExist = `GET count\(1\) FROM items`
)

func newTestDB(t *testing.T) (*MySQL, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	return &MySQL{
		conn:     sqlx.NewDb(db, "sqlmock"),
	}, mock
}

func TestConnFailsNoneExistent(t *testing.T) {
	c := config.Config{
		DbUrl: "localhost",
		DbUser: "nobody",
		DbPass: "foo",
		DbName: "blah",
	}

	_, err := NewMySQL(&c)
	assert.Error(t,err)
}

func TestMoveItem(t *testing.T) {
	db, mock := newTestDB(t)
	defer db.conn.Close()

	type testCase struct {
		testName  string
		direction string
		directionDB string
	}

	testCases := []testCase{
		{
			testName: "check in item",
			direction: "in",
			directionDB: "checked in",
		},
		{
			testName: "check out item",
			direction: "out",
			directionDB: "checked out",
		},
		{
			testName: "invalid direction",
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
