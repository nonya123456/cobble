package proto

import (
	"bytes"
	"errors"
	"io"

	"github.com/nonya123456/cobble/proto/types"
)

var (
	ErrInvalidPacketLength = errors.New("invalid packet length")
)

type Packet struct {
	ID   int32
	Data []byte
}

func ReadPacket(r io.Reader) (Packet, error) {
	var lengthProto types.VarInt
	if _, err := lengthProto.ReadFrom(r); err != nil {
		return Packet{}, err
	}
	length := int(lengthProto)

	var id types.VarInt
	idLength, err := id.ReadFrom(r)
	if err != nil {
		return Packet{}, err
	}

	dataLength := length - int(idLength)
	if dataLength < 0 {
		return Packet{}, ErrInvalidPacketLength
	}

	data := make([]byte, dataLength)
	if _, err := io.ReadFull(r, data); err != nil {
		return Packet{}, err
	}

	return Packet{
		ID:   int32(id),
		Data: data,
	}, nil
}

func WritePacket(w io.Writer, id int32, p io.WriterTo) error {
	buf := bytes.Buffer{}
	idProto := types.VarInt(id)
	if _, err := idProto.WriteTo(&buf); err != nil {
		return err
	}
	if _, err := p.WriteTo(&buf); err != nil {
		return err
	}

	body := buf.Bytes()
	length := types.VarInt(len(body))
	if _, err := length.WriteTo(w); err != nil {
		return err
	}

	if _, err := w.Write(body); err != nil {
		return err
	}

	return nil
}
