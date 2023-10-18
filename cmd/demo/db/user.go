package db

import (
	"errors"

	"github.com/surrealdb/surrealdb.go"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (db *DB) CreateUser(username, password string) (*User, error) {
	u := map[string]interface{}{
		"username": username,
		"password": password,
	}

	// Unmarshal data
	users, err := surrealdb.SmartUnmarshal[[]User](db.DB.Query("CREATE user SET username = $username, password = crypto::bcrypt::generate($password)", u))
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("no users returned")
	}

	return &users[0], nil
}

func (db *DB) GetUser(username, password string) (*User, error) {
	resp, err := surrealdb.SmartUnmarshal[[]User](db.DB.Query("SELECT id, username FROM user WHERE username = $username AND crypto::bcrypt::compare(password, $password) = true", map[string]interface{}{
		"username": username,
		"password": password,
	}))
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("no users returned")
	}

	return &resp[0], nil
}
