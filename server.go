package cobble

import (
	"io"
	"log"
	"net"

	"github.com/nonya123456/cobble/proto"
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

	state := proto.StateHandshaking

	for {
		p, err := proto.ReadPacket(conn)
		if err != nil {
			if err == io.EOF || err.Error() == "unexpected EOF" {
				log.Printf("Client %s disconnected\n", conn.RemoteAddr())
				return
			}

			log.Printf("Error reading packet from %s: %v\n", conn.RemoteAddr(), err)
		}

		switch state {
		case proto.StateHandshaking:
			switch p.ID {
			case proto.HandshakeID:
				handshake, err := proto.UnmarshalHandshake(p.Data)
				if err != nil {
					log.Printf("Failed to unmarshal handshake: %v\n", err)
					continue
				}

				if handshake.NextState == proto.StateStatus {
					state = proto.StateStatus
					log.Printf("Client %s entered Status state\n", conn.RemoteAddr())
				} else if handshake.NextState == proto.StateLogin {
					state = proto.StateLogin
					log.Printf("Client %s entered Login state\n", conn.RemoteAddr())
				} else {
					log.Printf("Invalid next state: %d from %s\n", handshake.NextState, conn.RemoteAddr())
					continue
				}
			}
		case proto.StateStatus:
			switch p.ID {
			case proto.StatusRequestID:
				res := proto.StatusResponse{JSONResponse: proto.String(`{"version":{"name":"1.23.1","protocol": 768}}`)}
				_, err := conn.Write(res.Packet().Bytes())
				if err != nil {
					log.Printf("Failed to write status response: %v\n", err)
				}

			case proto.PingRequestID:
				req, err := proto.UnMarshalPingRequest(p.Data)
				if err != nil {
					log.Printf("Failed to unmarshal ping request: %v\n", err)
					continue
				}

				res := proto.PingResponse{Payload: req.Payload}
				conn.Write(res.Packet().Bytes())
			default:
				log.Printf("Received unknown packet %v\n", p.ID)
			}
		}
	}
}
