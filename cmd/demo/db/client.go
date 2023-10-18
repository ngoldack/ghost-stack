package db

import "github.com/surrealdb/surrealdb.go"

type DB struct {
	*surrealdb.DB
}

func NewDB() (*DB, error) {
	db, err := surrealdb.New("ws://localhost:8000/rpc")
	if err != nil {
		return nil, err
	}

	if _, err = db.Signin(map[string]interface{}{
		"user": "admin",
		"pass": "admin",
	}); err != nil {
		return nil, err
	}

	if _, err = db.Use("ghost", "ghost"); err != nil {
		return nil, err
	}

	return &DB{
		DB: db,
	}, nil
}
