package dbstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lammer90/shortener/internal/models"
	"github.com/lammer90/shortener/internal/storage"
)

type dbStorage struct {
	db *sql.DB
}

// New dbStorage конструктор.
func New(db *sql.DB) dbStorage {
	initDB(db)
	return dbStorage{db: db}
}

// Save  сохранить ссылку с параметрами: key, value, userID.
func (d dbStorage) Save(key, value, userID string) error {
	_, err := d.db.ExecContext(context.Background(), `
        INSERT INTO shorts
        (short_url, original_url, user_id)
        VALUES
        ($1, $2, $3);
    `, key, value, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			shortURL := findByOriginal(d.db, value)
			return storage.NewErrConflict(shortURL, err)
		}
	}
	return err
}

// SaveBatch  сохранить батч с ссылками
func (d dbStorage) SaveBatch(shorts []*models.BatchToSave) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.Background()
	stmt, err := d.db.PrepareContext(ctx, `
        INSERT INTO shorts
        (short_url, original_url, user_id)
        VALUES
        ($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, short := range shorts {
		_, err := stmt.ExecContext(ctx, short.ShortURL, short.OriginalURL, short.UserID)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

// Find  Найти оригинальную ссылку по сокращенной
func (d dbStorage) Find(key string) (string, bool, error) {
	row := d.db.QueryRowContext(context.Background(), `
        SELECT
            s.original_url,
        	s.id_deleted
        FROM shorts s
        WHERE
            s.short_url = $1
    `, key)

	type result struct {
		OriginalURL string
		IsDeleted   bool
	}
	var r result
	err := row.Scan(&r.OriginalURL, &r.IsDeleted)
	if err != nil {
		return "", false, err
	}
	return r.OriginalURL, !r.IsDeleted, nil
}

// Find  Найти оригинальную ссылки по владельцу
func (d dbStorage) FindByUserID(userID string) (map[string]string, error) {
	resultMap := make(map[string]string)
	rows, err := d.db.QueryContext(context.Background(), `
        SELECT
            s.short_url,
            s.original_url
        FROM shorts s
        WHERE
            s.user_id = $1
    `, userID)

	if err != nil {
		return nil, err
	}

	type result struct {
		ShortURL    string
		OriginalURL string
	}

	for rows.Next() {
		var r result
		err = rows.Scan(&r.ShortURL, &r.OriginalURL)
		if err != nil {
			return nil, err
		}
		resultMap[r.ShortURL] = r.OriginalURL
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return resultMap, nil
}

// Delete  Удалить ссылки
func (d dbStorage) Delete(keys []string, userID string) error {
	query := `UPDATE shorts SET id_deleted = true WHERE short_url IN ($1, $2, $3, $4) AND user_id = $5`
	params := params(keys)
	_, err := d.db.ExecContext(context.Background(), query, params[0], params[1], params[2], params[3], userID)
	return err
}

func initDB(db *sql.DB) {
	ctx := context.Background()
	db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS shorts (
            short_url varchar,
            original_url varchar,
            user_id varchar,
            id_deleted boolean default false
        )
    `)
	db.ExecContext(ctx, `
        CREATE UNIQUE INDEX IF NOT EXISTS shorts_original_url_idx ON shorts (original_url)
    `)
}

func findByOriginal(db *sql.DB, value string) string {
	row := db.QueryRowContext(context.Background(), `
        SELECT
            s.short_url
        FROM shorts s
        WHERE
            s.original_url = $1
    `, value)

	var shortURL string
	err := row.Scan(&shortURL)
	if err != nil {
		return ""
	}
	return shortURL
}

func params(keys []string) [5]string {
	var arr [5]string
	for i := 0; i < 4; i++ {
		if i+1 <= len(keys) {
			arr[i] = keys[i]
		} else {
			arr[i] = ""
		}
	}
	return arr
}
