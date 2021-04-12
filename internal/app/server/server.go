package server

import (
	"net/http"

	"github.com/BroNaz/queue_broker/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config *Configer
	logger *logrus.Logger
	router *mux.Router
	store  *store.RuntimeStorege
}

func NewServer(config *Configer) *Server {
	return &Server{
		logger: logrus.New(),
		config: config,
		router: mux.NewRouter(),
		store:  store.NewRuntimeStorege(),
	}
}

func (s *Server) Start() error {
	if err := s.configLogger(); err != nil {
		return err
	}
	s.logger.Info("Start Server")

	s.configRouter()

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configRouter() {
	s.router.HandleFunc("/{key}", s.PushMessagesInAQueue()).Methods("PUT")
	s.router.HandleFunc("/{key}", s.PopFromTheQueue()).Methods("GET")

}

func (s *Server) configLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}
