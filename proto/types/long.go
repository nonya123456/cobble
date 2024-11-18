package types

import (
	"encoding/binary"
	"io"
)

type Long int64

func (l *Long) ReadFrom(r io.Reader) (int64, error) {
	buffer := make([]byte, 8)
	n, err := io.ReadFull(r, buffer)
	if err != nil {
		return int64(n), err
	}

	*l = Long(binary.BigEndian.Uint64(buffer))
	return int64(n), nil
}

func (l *Long) WriteTo(w io.Writer) (int64, error) {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, uint64(*l))
	n, err := w.Write(buffer)
	return int64(n), err
}
