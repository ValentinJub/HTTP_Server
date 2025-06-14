package connection

import (
	"log"
	"net"
)

// ConnPoolHandler is an interface that defines methods for managing a pool of connections.
type ConnPoolHandler interface {
	Add(net.Conn)
	Remove(key string) error
}

type HTTPConnPoolHandler struct {
	pool map[string]net.Conn
}

func NewConnPoolHandler() ConnPoolHandler {
	return &HTTPConnPoolHandler{
		pool: make(map[string]net.Conn),
	}
}

func (h *HTTPConnPoolHandler) Add(conn net.Conn) {
	key := conn.RemoteAddr().String()
	h.pool[key] = conn
	log.Printf("Connection added to pool: %s", key)
}

func (h *HTTPConnPoolHandler) Remove(key string) error {
	_, exists := h.pool[key]
	if !exists {
		return net.UnknownNetworkError("connection not found in pool")
	}

	delete(h.pool, key)

	log.Printf("Connection removed from pool: %s", key)
	return nil
}
