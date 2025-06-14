package router

import (
	"log"
	"net"

	conn "github.com/codecrafters-io/http-server-starter-go/internal/connection"
)

type Server interface {
	ListenAndServe() error
}

type HTTPServer struct {
	address  string
	port     string
	connPool conn.ConnPoolHandler
}

func NewHTTPServer(address, port string, connPool conn.ConnPoolHandler) *HTTPServer {
	return &HTTPServer{
		address:  address,
		port:     port,
		connPool: connPool,
	}
}

func (s *HTTPServer) ListenAndServe() error {
	FQDN := s.address + ":" + s.port
	listener, err := net.Listen("tcp", FQDN)
	log.Printf("Listening on %s", FQDN)
	if err != nil {
		log.Printf("Failed to start server on %s: %v", FQDN, err)
		return err
	}
	for {
		c, err := listener.Accept()
		if err != nil {
			return err
		}
		go conn.NewConnHandler(s.connPool).Handle(c)
	}
}
