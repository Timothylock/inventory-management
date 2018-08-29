package persistence

import (
	"errors"
	"fmt"

	"inventory-management/config"
	"inventory-management/items"

	"github.com/jmoiron/sqlx"
)

type MySQL struct {
	conn     *sqlx.DB
}

func NewMySQL(cfg *config.Config) (*MySQL, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		cfg.DbUser, cfg.DbPass, cfg.DbUrl, cfg.DbName)

	conn, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		return nil, errors.New("unable to open db connection")
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
		"GET count(1) FROM items WHERE ID = ?",
		ID,
	)
	if err != nil {
		return false, err
	}

	return count > 0, err
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