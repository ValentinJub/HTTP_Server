package main

import (
	"os"

	conn "github.com/codecrafters-io/http-server-starter-go/internal/connection"
	"github.com/codecrafters-io/http-server-starter-go/internal/router"
)

const (
	PORT         = "4221"
	ADDRESS      = "localhost"
	FULL_ADDRESS = ADDRESS + ":" + PORT
)

func main() {
	server := router.NewHTTPServer(ADDRESS, PORT, conn.NewConnPoolHandler())
	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
