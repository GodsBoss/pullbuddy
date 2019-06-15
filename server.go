package pullbuddy

import (
	"net/http"
)

type Server struct {
	Addr string
}

func (server *Server) Start() error {
	sch := newScheduler()
	go sch.run()
	httpServ := http.Server{
		Addr:    orDefaultAddr(server.Addr),
		Handler: newHandler(sch),
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
