package pullbuddy

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	Addr string
}

func (server *Server) Start() error {
	httpServ := http.Server{
		Addr:    orDefaultAddr(server.Addr),
		Handler: chi.NewRouter(),
	}
	return httpServ.ListenAndServe()
}

const DefaultAddr = ":30666"

func orDefaultAddr(addr string) string {
	if addr != "" {
		return addr
	}
	return DefaultAddr
}
