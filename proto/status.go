package proto

import (
	"encoding/binary"
)

const (
	StatusRequestID = 0x00
	PingRequestID   = 0x01
)

const (
	StatusResponseID = 0x00
	PingResponseID   = 0x01
)

type StatusRequest struct{}

func (s StatusRequest) Packet() Packet {
	return Packet{
		ID:   StatusRequestID,
		Data: []byte{},
	}
}

type StatusResponse struct {
	JSONResponse String
}

func (s StatusResponse) Packet() Packet {
	return Packet{
		ID:   StatusResponseID,
		Data: s.JSONResponse.Bytes(),
	}
}

type PingRequest struct {
	Payload Long
}

func (p PingRequest) Packet() Packet {
	return Packet{
		ID:   PingRequestID,
		Data: p.Payload.Bytes(),
	}
}

func UnMarshalPingRequest(data []byte) (PingRequest, error) {
	if len(data) != 8 {
		return PingRequest{}, ErrInvalidDataLength
	}

	payload := Long(binary.BigEndian.Uint64(data))
	return PingRequest{Payload: payload}, nil
}

type PingResponse struct {
	Payload Long
}

func (p PingResponse) Packet() Packet {
	return Packet{
		ID:   PingResponseID,
		Data: p.Payload.Bytes(),
	}
}
