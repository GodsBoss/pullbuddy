package pullbuddy

import (
	"net/http"

	dc "github.com/moby/moby/client"
)

type Server struct {
	Addr string
}

func (server *Server) Start() error {
	dockerClient, err := dc.NewClientWithOpts()
	if err != nil {
		return err
	}
	sch := newScheduler()
	sch.puller = newDockerImagePuller(dockerClient)
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
