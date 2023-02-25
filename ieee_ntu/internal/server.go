package internal

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Server struct {
	r  *gin.Engine
	db *sql.DB
}

func (s *Server) attachHandler() {
	s.r.GET("/query", s.queryHandler)
	s.r.GET("/tx", s.txHandler)
}

func (s *Server) Serve() {
	s.r.Run(":8080")
}

func NewServer(db *sql.DB) *Server {
	s := new(Server)
	s.r = gin.Default()
	s.db = db
	s.attachHandler()
	return s
}
