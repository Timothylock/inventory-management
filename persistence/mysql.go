package persistence

import (
	"errors"
	"fmt"

	"github.com/Timothylock/inventory-management/config"
	"github.com/Timothylock/inventory-management/items"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQL struct {
	conn *sqlx.DB
}

func NewMySQL(cfg *config.Config) (*MySQL, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		cfg.DbUser, cfg.DbPass, cfg.DbUrl, cfg.DbName)

	conn, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		return nil, err
	}

	return &MySQL{
		conn: conn,
	}, nil
}

// doesIDExist returns whether the ID is found in the database
func (m *MySQL) doesIDExist(ID string) (bool, error) {
	var count int

	err := m.conn.Get(
		&count,
		"SELECT count(1) FROM items WHERE ID = ?",
		ID,
	)
	if err != nil {
		return false, err
	}

	return count > 0, err
}

func (m *MySQL) SearchItems(search string) (items.ItemDetailList, error) {
	dl := items.ItemDetailList{}
	err := m.conn.Select(
		&dl,
		`SELECT search.ID AS ID, NAME, CATEGORY, PICTURE_URL, DETAILS, LOCATION, USERNAME, QUANTITY, STATUS FROM
		(
		SELECT * FROM items WHERE MATCH (ID, NAME, CATEGORY, PICTURE_URL, DETAILS, LOCATION) AGAINST (? IN NATURAL LANGUAGE MODE)
		) AS search
		JOIN users ON search.LAST_PERFORMED_BY = users.ID`,
		search,
	)

	return dl, err
}

func (m *MySQL) MoveItem(ID, direction string) error {
	var status string
	if direction == "in" {
		status = "checked in"
	} else if direction == "out" {
		status = "checked out"
	} else {
		return errors.New("invalid direction")
	}

	exist, err := m.doesIDExist(ID)
	if err != nil {
		return err
	}
	if !exist {
		return items.ItemNotFoundErr
	}

	_, err = m.conn.Exec(
		"UPDATE items SET STATUS = ? WHERE ID = ?",
		status, ID,
	)
	return err
}

func (m *MySQL) DeleteItem(ID string) error {
	r, err := m.conn.Exec(
		"DELETE FROM items WHERE ID = ?",
		ID,
	)
	if err != nil {
		return err
	}

	ra, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if ra <= 0 {
		return items.ItemNotFoundErr
	}

	return err
}
