package cobble

import (
	"log"
	"net"
)

type Server struct {
	Addr string
}

func (s Server) Run() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer l.Close()

	log.Printf("Server listening on %s\n", s.Addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s Server) handle(conn net.Conn) {
	defer conn.Close()
	log.Printf("New connection from %s\n", conn.RemoteAddr())
}
