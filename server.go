package cobble

import (
	"bytes"
	"io"
	"log"
	"net"

	"github.com/nonya123456/cobble/proto"
	"github.com/nonya123456/cobble/proto/handshaking"
	"github.com/nonya123456/cobble/proto/status"
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

	state := 0
	for {
		p, err := proto.ReadPacket(conn)
		if err != nil {
			if err == io.EOF || err.Error() == "unexpected EOF" {
				log.Printf("Client %s disconnected\n", conn.RemoteAddr())
				return
			}

			log.Printf("Error reading packet from %s: %v\n", conn.RemoteAddr(), err)
		}

		r := bytes.NewReader(p.Data)

		switch state {
		case 0:
			switch p.ID {
			case handshaking.HandshakeID:
				var handshake handshaking.Handshake
				if _, err := handshake.ReadFrom(r); err != nil {
					log.Printf("Failed to read handshake: %v\n", err)
					continue
				}

				state = int(handshake.NextState)
			default:
				log.Printf("Received unknown packet %v\n", p.ID)
			}
		case 1:
			switch p.ID {
			case status.StatusRequestID:
				var req status.StatusRequest
				if _, err := req.ReadFrom(r); err != nil {
					log.Printf("Failed to read status request: %v\n", err)
					continue
				}

				res := status.StatusResponse{JSONResponse: `{"version":{"name":"1.23.1","protocol": 768}}`}
				if err := proto.WritePacket(conn, status.StatusResponseID, &res); err != nil {
					log.Printf("Failed to write status response: %v\n", err)
				}

			case status.PingRequestID:
				var req status.PingRequest
				if _, err := req.ReadFrom(r); err != nil {
					log.Printf("Failed to read ping request: %v\n", err)
					continue
				}

				payload := req.Payload
				res := status.PingResponse{Payload: payload}
				if err := proto.WritePacket(conn, status.PingResponseID, &res); err != nil {
					log.Printf("Failed to write ping response: %v\n", err)
				}
			default:
				log.Printf("Received unknown packet %v\n", p.ID)
			}

		default:
			log.Printf("Unimplemented state\n")
		}
	}
}
