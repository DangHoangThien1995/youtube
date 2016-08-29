package server

import (
	"net/http"
)

func Start(cfg Config) {
	s := setup(cfg)
	listenAddr := cfg.Server.Addr + ":" + cfg.Server.Port
	l.Println("API_server is listening on", listenAddr)
	http.ListenAndServe(listenAddr, s.Handler)
}
