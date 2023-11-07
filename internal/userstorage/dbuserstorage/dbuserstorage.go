package dbuserstorage

import (
	"context"
	"database/sql"
)

type dbUserStorage struct {
	db *sql.DB
}

func New(db *sql.DB) dbUserStorage {
	initDB(db)
	return dbUserStorage{db: db}
}

func (d dbUserStorage) Save(value string) error {
	_, err := d.db.ExecContext(context.Background(), `
        INSERT INTO users
        (user_id)
        VALUES
        ($1);
    `, value)
	return err
}

func (d dbUserStorage) Find(key string) (string, bool, error) {
	row := d.db.QueryRowContext(context.Background(), `
        SELECT
            u.user_id
        FROM users u
        WHERE
            u.user_id = $1
    `, key)

	var value string
	err := row.Scan(&value)
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func initDB(db *sql.DB) {
	ctx := context.Background()
	db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS users (
            user_id varchar Primary Key
        )
    `)
}
