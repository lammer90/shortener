package dbstorage

import (
	"context"
	"database/sql"
	"github.com/lammer90/shortener/internal/models"
)

type dbStorage struct {
	db *sql.DB
}

func New(db *sql.DB) dbStorage {
	initDB(db)
	return dbStorage{db: db}
}

func (d dbStorage) Save(key, value string) error {
	_, err := d.db.ExecContext(context.Background(), `
        INSERT INTO shorts
        (short_url, original_url)
        VALUES
        ($1, $2);
    `, key, value)
	return err
}

func (d dbStorage) SaveBatch(shorts []*models.BatchToSave) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.Background()
	stmt, err := d.db.PrepareContext(ctx, `
        INSERT INTO shorts
        (short_url, original_url)
        VALUES
        ($1, $2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, short := range shorts {
		_, err := stmt.ExecContext(ctx, short.ShortURL, short.OriginalURL)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func (d dbStorage) Find(key string) (string, bool, error) {
	row := d.db.QueryRowContext(context.Background(), `
        SELECT
            s.original_url
        FROM shorts s
        WHERE
            s.short_url = $1
    `, key)

	var value string
	err := row.Scan(&value)
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func initDB(db *sql.DB) {
	db.ExecContext(context.Background(), `
        CREATE TABLE IF NOT EXISTS shorts (
            short_url varchar PRIMARY KEY,
            original_url varchar
        )
    `)
}
