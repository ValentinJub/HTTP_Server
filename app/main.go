package main

import (
	"fmt"
	"net"
	"os"
)

const (
	PORT    = 4221
	ADDRESS = "localhost"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ADDRESS, PORT))
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	_, err = l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}
