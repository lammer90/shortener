package ping

import (
	"net/http"
)

type Provider interface {
	Ping() error
}

type Ping struct {
	provider Provider
}

func New(provider Provider) Ping {
	return Ping{provider: provider}
}

func (p Ping) Ping(res http.ResponseWriter, req *http.Request) {
	if err := p.provider.Ping(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
