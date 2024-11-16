package cobble

import (
	"io"
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

	for {
		p, err := ReadPacket(conn)
		if err != nil {
			if err == io.EOF || err.Error() == "unexpected EOF" {
				log.Printf("Client %s disconnected\n", conn.RemoteAddr())
				return
			}

			log.Printf("Error reading packet from %s: %v\n", conn.RemoteAddr(), err)
		}

		log.Printf("Received packet %v\n", p)
	}
}
