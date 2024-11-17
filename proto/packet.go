package proto

import (
	"errors"
	"io"
)

var (
	ErrZeroDataLength      = errors.New("zero data length")
	ErrInvalidDataLength   = errors.New("invalid data length")
	ErrInvalidPacketLength = errors.New("invalid packet length")
)

type Packet struct {
	ID   VarInt
	Data []byte
}

func (p Packet) Bytes() []byte {
	idBytes := p.ID.Bytes()
	dataBytes := p.Data
	length := VarInt(len(idBytes) + len(dataBytes)).Bytes()
	return append(append(length, idBytes...), dataBytes...)
}

func ReadPacket(r io.Reader) (Packet, error) {
	length, err := ReadVarInt(r)
	if err != nil {
		return Packet{}, err
	}

	id, err := ReadVarInt(r)
	if err != nil {
		return Packet{}, err
	}

	dataLen := int(length) - len(id.Bytes())
	if dataLen < 0 {
		return Packet{}, ErrInvalidPacketLength
	}

	data := make([]byte, dataLen)
	if _, err = io.ReadFull(r, data); err != nil {
		return Packet{}, err
	}

	return Packet{
		ID:   id,
		Data: data,
	}, nil
}
