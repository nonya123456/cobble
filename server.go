package cobble

import (
	"bytes"
	"io"
	"log"
	"net"

	"github.com/nonya123456/cobble/proto"
	"github.com/nonya123456/cobble/proto/packets"
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

	state := packets.StateHandshaking

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
		case packets.StateHandshaking:
			switch p.ID {
			case packets.HandshakeID:
				var handshake packets.Handshake
				if _, err := handshake.ReadFrom(r); err != nil {
					log.Printf("Failed to read handshake: %v\n", err)
					continue
				}

				if handshake.NextState == packets.StateStatus {
					state = packets.StateStatus
					log.Printf("Client %s entered Status state\n", conn.RemoteAddr())
				} else if handshake.NextState == packets.StateLogin {
					state = packets.StateLogin
					log.Printf("Client %s entered Login state\n", conn.RemoteAddr())
				} else {
					log.Printf("Invalid next state: %d from %s\n", handshake.NextState, conn.RemoteAddr())
					continue
				}
			}
		case packets.StateStatus:
			switch p.ID {
			case packets.StatusRequestID:
				var req packets.StatusRequest
				if _, err := req.ReadFrom(r); err != nil {
					log.Printf("Failed to read status request: %v\n", err)
					continue
				}

				res := packets.StatusResponse{JSONResponse: `{"version":{"name":"1.23.1","protocol": 768}}`}
				if err := proto.WritePacket(conn, packets.StatusResponseID, &res); err != nil {
					log.Printf("Failed to write status response: %v\n", err)
				}

			case packets.PingRequestID:
				var req packets.PingRequest
				if _, err := req.ReadFrom(r); err != nil {
					log.Printf("Failed to read ping request: %v\n", err)
					continue
				}

				payload := req.Payload
				res := packets.PingResponse{Payload: payload}
				if err := proto.WritePacket(conn, packets.PingResponseID, &res); err != nil {
					log.Printf("Failed to write ping response: %v\n", err)
				}

			default:
				log.Printf("Received unknown packet %v\n", p.ID)
			}
		}
	}
}
