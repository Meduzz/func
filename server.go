package funclib

import (
	"log"
	"net"
	"os"
)

type (
	Server struct {
		srv     net.Listener
		running bool
	}
)

// Standard listen func with standard go net.Listen params.
func Listen(network, addr string) (*Server, error) {
	listener, err := net.Listen(network, addr)

	log.Printf("Server listening to %s.\n", addr)

	return &Server{listener, true}, err
}

// Starts the server using the provided handler to handle new connections
func (s *Server) Run(handler func(net.Conn)) {
	for {
		if !s.running {
			break
		}

		conn, err := s.srv.Accept()

		if err != nil {
			if err != net.ErrClosed {
				break
			}

			log.Printf("Accepting connection threw error: %v\n", err)
			continue
		}

		go handler(conn)
	}
}

// Close the server terminating all connections?
func (s *Server) Close() error {
	s.running = false
	err := s.srv.Close()

	if err != nil {
		return err
	}

	if s.srv.Addr().Network() == "unix" {
		return os.Remove(s.srv.Addr().String())
	}

	return nil
}
