package persistence

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"github.com/Timothylock/inventory-management/config"
	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/users"

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
		SELECT * FROM items WHERE (MATCH (ID, NAME, CATEGORY, DETAILS, LOCATION) AGAINST (? IN NATURAL LANGUAGE MODE) AND DELETED=0) OR (ID = ? AND DELETED=0)
		) AS search
		JOIN users ON search.LAST_PERFORMED_BY = users.ID`,
		search, search,
	)

	return dl, err
}

func (m *MySQL) MoveItem(ID, direction string, userID int) error {
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
		"UPDATE items SET STATUS = ?, LAST_PERFORMED_BY = ? WHERE ID = ?",
		status, userID, ID,
	)

	if err == nil {
		m.addLog(userID, ID, status, "")
	}

	return err
}

func (m *MySQL) AddItem(obj items.ItemDetail, overwrite bool) error {
	exist, err := m.doesIDExist(obj.ID)
	if err != nil {
		return err
	}
	if exist && !overwrite {
		return items.ItemAlreadyExistsErr
	}
	if !exist && overwrite {
		return items.ItemNotFoundErr
	}

	if !overwrite {
		_, err = m.conn.Exec(
			`INSERT INTO items (ID, NAME, CATEGORY, PICTURE_URL, DETAILS, LOCATION, LAST_PERFORMED_BY, QUANTITY, STATUS)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) `,
			obj.ID, obj.Name, obj.Category, obj.PictureURL, obj.Details, obj.Location, obj.LastPerformedBy, obj.Quantity, "checked in")
	} else if overwrite {
		_, err = m.conn.Exec(
			`UPDATE items SET ID = ?, NAME = ?, CATEGORY = ?, PICTURE_URL = ?, DETAILS = ?, LOCATION = ?, LAST_PERFORMED_BY = ?, QUANTITY = ? WHERE ID = ?`,
			obj.ID, obj.Name, obj.Category, obj.PictureURL, obj.Details, obj.Location, obj.LastPerformedBy, obj.Quantity, obj.ID)
	}

	if err == nil {
		uid, err := strconv.Atoi(obj.LastPerformedBy)
		if err == nil {
			m.addLog(uid, obj.ID, "add", fmt.Sprintf("overwrite/skip exist check flag was recieved as %t", overwrite))
		}
	}

	return err
}

func (m *MySQL) DeleteItem(ID string, userID int) error {
	r, err := m.conn.Exec(`UPDATE items SET DELETED=1, LAST_PERFORMED_BY=? WHERE ID=?`,
		userID, ID,
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

	m.addLog(userID, ID, "delete", "")

	return err
}

func (m *MySQL) addLog(uID int, objID, action, details string) error {
	_, err := m.conn.Exec(`INSERT INTO logs (USERID, OBJECTID, ACTION, DETAILS, DATE) VALUES
	(?, ?, ?, ?, NOW())`, uID, objID, action, details)
	return err
}

type MultiUserDB []UserDB
type UserDB struct {
	ID         int    `db:"ID"`
	IsSysAdmin int    `db:"ISSYSADMIN"`
	Email      string `db:"EMAIL"`
	Token      string `db:"TOKEN"`
	Username   string `db:"USERNAME"`
}

// GetUser gets the given user if possible
func (m *MySQL) GetUser(username, password string) (users.User, error) {
	var user users.User
	user.Valid = false

	hasher := sha1.New()
	hasher.Write([]byte(password))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	var userdb UserDB
	err := m.conn.Get(
		&userdb,
		"SELECT ID, ISSYSADMIN, EMAIL, TOKEN, USERNAME FROM users WHERE USERNAME = ? AND PASSWORD = ? AND ACTIVE = 1",
		username, sha,
	)
	if err == sql.ErrNoRows {
		return user, nil
	} else if err != nil {
		return user, err
	}

	user.Valid = true
	user.ID = userdb.ID
	user.IsSysAdmin = userdb.IsSysAdmin == 1
	user.Email = userdb.Email
	user.Token = userdb.Token
	user.Username = userdb.Username

	return user, err
}

// GetUserByToken returns the user by token
func (m *MySQL) GetUserByToken(token string) (users.User, error) {
	var user users.User
	user.Valid = false
	var userdb UserDB
	err := m.conn.Get(
		&userdb,
		"SELECT ID, ISSYSADMIN, EMAIL, TOKEN, USERNAME FROM users WHERE TOKEN = ? AND ACTIVE = 1",
		token,
	)
	if err == sql.ErrNoRows {
		return user, nil
	} else if err != nil {
		return user, err
	}

	user.Valid = true
	user.ID = userdb.ID
	user.IsSysAdmin = userdb.IsSysAdmin == 1
	user.Email = userdb.Email
	user.Token = userdb.Token
	user.Username = userdb.Username

	return user, err
}

// AddUser adds a new user or updates and existing one
func (m *MySQL) AddUser(username, email, password string, isSysAdmin bool) error {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	token := generateToken()

	_, err := m.conn.Exec(
		`INSERT INTO users (USERNAME, EMAIL, PASSWORD, TOKEN, ISSYSADMIN) VALUES (?, ?, ?, ?, ?)`,
		username, email, sha, token, isSysAdmin,
	)

	return err
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// GetUsers gets all the active users
func (m *MySQL) GetUsers() (users.MultipleUsers, error) {
	dl := MultiUserDB{}
	err := m.conn.Select(
		&dl,
		"SELECT ID, ISSYSADMIN, EMAIL, TOKEN, USERNAME FROM users WHERE ACTIVE = 1",
	)

	ret := users.MultipleUsers{}

	for _, u := range dl {
		ret = append(ret, users.User{
			Valid:      true,
			ID:         u.ID,
			IsSysAdmin: u.IsSysAdmin == 1,
			Email:      u.Email,
			Token:      u.Token,
			Username:   u.Username,
		})

	}

	return ret, err
}
