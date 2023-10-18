package db

import (
	"errors"
	"time"

	"github.com/surrealdb/surrealdb.go"
)

type Session struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Expiration time.Time `json:"expiration"`
}

func (db *DB) ValidateSession(id string) (bool, error) {
	resp, err := surrealdb.SmartUnmarshal[[]Session](db.Query("SELECT id, user_id, expiration FROM session WHERE id = $id AND ", map[string]interface{}{
		"id": id,
	}))

	if err != nil {
		return false, err
	}

	return len(resp) > 1, nil
}

func (db *DB) GetSession(userID string) (*Session, error) {
	resp, err := surrealdb.SmartUnmarshal[[]Session](db.Query("SELECT id, user_id, expiration FROM session WHERE user_id = $user_id AND time::now() < expiration", map[string]interface{}{
		"user_id": userID,
	}))
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("no sessions returned")
	}

	return &resp[0], nil
}

func (db *DB) GetOrCreateSession(userID string) (*Session, error) {
	if session, err := db.GetSession(userID); err == nil {
		return session, nil
	}

	s := map[string]interface{}{
		"user_id": userID,
	}

	// Unmarshal data
	resp, err := surrealdb.SmartUnmarshal[[]Session](db.Query("CREATE session SET user_id = $user_id, expiration = time::now() + 1h", s))
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("no sessions returned")
	}

	return &resp[0], nil
}
