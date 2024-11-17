package proto

import (
	"bytes"
)

type State VarInt

func (s State) Bytes() []byte {
	return VarInt(s).Bytes()
}

const (
	StateHandshaking State = 0
	StateStatus      State = 1
	StateLogin       State = 2
	StatePlay        State = 3
)

const HandshakeID = 0x00

type Handshake struct {
	ProtocolVersion VarInt
	ServerAddress   String
	ServerPort      UnsignedShort
	NextState       State
}

func (h Handshake) Packet() Packet {
	buf := &bytes.Buffer{}
	buf.Write(h.ProtocolVersion.Bytes())
	buf.Write(h.ServerAddress.Bytes())
	buf.Write(h.ServerPort.Bytes())
	buf.Write(h.NextState.Bytes())

	return Packet{
		ID:   HandshakeID,
		Data: buf.Bytes(),
	}
}

func UnmarshalHandshake(data []byte) (Handshake, error) {
	if len(data) == 0 {
		return Handshake{}, ErrZeroDataLength
	}

	r := bytes.NewReader(data)
	protocolVersion, err := ReadVarInt(r)
	if err != nil {
		return Handshake{}, err
	}

	serverAddress, err := ReadString(r)
	if err != nil {
		return Handshake{}, err
	}

	serverPort, err := ReadUnsignedShort(r)
	if err != nil {
		return Handshake{}, err
	}

	var nextState VarInt
	nextState, err = ReadVarInt(r)
	if err != nil {
		return Handshake{}, err
	}

	return Handshake{
		ProtocolVersion: protocolVersion,
		ServerAddress:   serverAddress,
		ServerPort:      serverPort,
		NextState:       State(nextState),
	}, nil
}
