package types

import (
	"encoding/binary"
	"io"
)

type UnsignedShort uint16

func (u *UnsignedShort) ReadFrom(r io.Reader) (int64, error) {
	buffer := make([]byte, 2)
	n, err := io.ReadFull(r, buffer)
	if err != nil {
		return int64(n), err
	}

	*u = UnsignedShort(binary.BigEndian.Uint16(buffer))
	return int64(n), nil
}

func (u *UnsignedShort) WriteTo(w io.Writer) (int64, error) {
	buffer := make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, uint16(*u))
	n, err := w.Write(buffer)
	return int64(n), err
}
