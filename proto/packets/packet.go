package packets

import (
	"bytes"
	"io"

	"github.com/nonya123456/cobble/proto/types"
)

type Packet interface {
	ID() int32
	io.ReaderFrom
	io.WriterTo
}

func ReadPacketID(r io.Reader) (int32, error) {
	var length types.VarInt
	if _, err := length.ReadFrom(r); err != nil {
		return 0, err
	}

	var id types.VarInt
	if _, err := id.ReadFrom(r); err != nil {
		return 0, err
	}

	return int32(id), nil
}

func WritePacket(w io.Writer, p Packet) error {
	buf := bytes.Buffer{}
	id := types.VarInt(p.ID())
	if _, err := id.WriteTo(&buf); err != nil {
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
