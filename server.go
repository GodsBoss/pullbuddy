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
		Addr:    orDefaultAddr(server.Addr, DefaultServerAddr),
		Handler: newHandler(sch),
	}
	return httpServ.ListenAndServe()
}

const DefaultServerAddr = ":30666"

func orDefaultAddr(addr, defaultAddr string) string {
	if addr != "" {
		return addr
	}
	return defaultAddr
}
