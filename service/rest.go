package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/ipfs-cluster-monitor/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// RestServer ...
type RestServer struct {
	*gin.Engine
	config *config.Configure
	server *http.Server
	Port   string
}

// NewRestServer ...
func NewRestServer(cfg *config.Configure) *RestServer {
	s := &RestServer{
		Engine: gin.Default(),
		config: cfg,
		Port:   config.DefaultString("", ":7783"),
	}
	return s
}

// Start ...
func (s *RestServer) Start() {
	if !s.config.REST.Enable {
		return
	}

	Router(s.Engine)
	s.server = &http.Server{
		Addr:    s.Port,
		Handler: s.Engine,
	}
	go func() {
		log.Printf("Listening and serving HTTP on %s\n", s.Port)
		if err := s.server.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
}

// Stop ...
func (s *RestServer) Stop() {
	if err := s.server.Shutdown(nil); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

// CheckPrefix ...
func CheckPrefix(url string) string {
	if strings.Index(url, "http") != 0 {
		return "http://" + url
	}
	return url
}
