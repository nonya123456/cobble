package proto

import "io"

type Packet struct {
	ID   VarInt
	Data []byte
}

func (p Packet) Bytes() []byte {
	idBytes := p.ID.Bytes()
	dataBytes := p.Data
	length := VarInt(len(dataBytes)).Bytes()
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

	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return Packet{}, err
	}

	return Packet{
		ID:   id,
		Data: data,
	}, nil
}
