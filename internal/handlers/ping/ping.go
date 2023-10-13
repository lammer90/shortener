package ping

import (
	"database/sql"
	"net/http"
)

type Ping struct {
	db *sql.DB
}

func New(db *sql.DB) Ping {
	return Ping{db: db}
}

func (p Ping) Ping() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := p.db.Ping(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
}
