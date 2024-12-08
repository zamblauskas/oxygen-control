package api

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router         *gin.Engine
	triggerHandler *TriggerHandler
}

func NewServer(triggerHandler *TriggerHandler) *Server {
	return &Server{
		router:         gin.Default(),
		triggerHandler: triggerHandler,
	}
}

func (s *Server) SetupRoutes() {
	s.router.POST("/triggers", s.triggerHandler.AddTrigger)
	s.router.GET("/triggers", s.triggerHandler.GetTriggers)
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
