package ping

import (
	"net/http"
)

// Provider ping
type Provider interface {
	Ping() error
}

// Ping модель
type Ping struct {
	provider Provider
}

// New Ping констуктор
func New(provider Provider) Ping {
	return Ping{provider: provider}
}

// Ping проверить доступность бд
func (p Ping) Ping(res http.ResponseWriter, req *http.Request) {
	if err := p.provider.Ping(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
