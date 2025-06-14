package connection

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/http-server-starter-go/internal/handler"
)

type HTTPConnHandler interface {
	Handle(conn net.Conn)
	Read(conn net.Conn) (string, error)
	Write(conn net.Conn, response string) error
}

type HTTPConnHandlerImpl struct {
	connPool ConnPoolHandler
}

func NewConnHandler(cp ConnPoolHandler) HTTPConnHandler {
	return &HTTPConnHandlerImpl{connPool: cp}
}

func (h *HTTPConnHandlerImpl) Handle(conn net.Conn) {
	defer conn.Close()
	defer h.connPool.Remove(conn.RemoteAddr().String())
	h.connPool.Add(conn)
	for {
		rawRequest, err := h.Read(conn)
		if err != nil {
			if err.Error() == "failed to read from connection: EOF" {
				fmt.Println("Connection closed by client:", conn.RemoteAddr().String())
				return
			}
			fmt.Printf("Error reading from connection: %v\n", err)
			return
		}
		decoder := handler.NewRequestDecoder()
		request := decoder.Decode(rawRequest)
		if request == nil {
			println("Failed to decode request:", rawRequest)
			return
		}
		fmt.Println("Received request:", request.String())

		responseStr := handler.NewResponseEncoder().Encode(handler.NewHTTPResponse(
			handler.NewStatusLine(request.Request.HTTPVersion, 200, "OK"),
			map[string]string{},
			"",
		))

		if err := h.Write(conn, responseStr); err != nil {
			fmt.Printf("Error writing response: %v\n", err)
			return
		}
		// fmt.Println("Sent response:", responseStr)
	}
}

func (h *HTTPConnHandlerImpl) Read(conn net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read from connection: %v", err)
	}
	if n == 0 {
		return "", fmt.Errorf("failed to read from connection: EOF")
	}
	rawRequest := string(buffer[:n])
	return rawRequest, nil
}

func (h *HTTPConnHandlerImpl) Write(conn net.Conn, response string) error {
	_, err := conn.Write([]byte(response))
	if err != nil {
		return fmt.Errorf("failed to write response: %v", err)
	}
	return nil
}
