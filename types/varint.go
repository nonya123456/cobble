package types

import "io"

type VarInt int32

// ReadFrom implements io.ReaderFrom.
func (v *VarInt) ReadFrom(r io.Reader) (int64, error) {
	var n int64
	var result int32
	var shift uint
	for {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return n, err
		}
		n++
		result |= int32(b[0]&0b01111111) << shift
		if b[0]&0b10000000 == 0 {
			break
		}
		shift += 7
	}

	*v = VarInt(result)
	return n, nil
}

// WriteTo implements io.WriterTo.
func (v *VarInt) WriteTo(r io.Writer) (int64, error) {
	value := *v
	var p []byte
	for {
		temp := byte(value & 0b01111111)
		value >>= 7
		if value != 0 {
			temp |= 0b10000000
		}
		p = append(p, temp)
		if value == 0 {
			break
		}
	}

	n, err := r.Write(p)
	return int64(n), err
}

var (
	_ io.WriterTo   = (*VarInt)(nil)
	_ io.ReaderFrom = (*VarInt)(nil)
)
