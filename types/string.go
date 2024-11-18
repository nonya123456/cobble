package types

import "io"

type String string

// ReadFrom implements io.ReaderFrom.
func (s *String) ReadFrom(r io.Reader) (int64, error) {
	var length VarInt
	n1, err := length.ReadFrom(r)
	if err != nil {
		return n1, err
	}

	data := make([]byte, length)
	n2, err := io.ReadFull(r, data)
	if err != nil {
		return n1 + int64(n2), err
	}

	*s = String(data)
	return n1 + int64(n2), nil
}

// WriteTo implements io.WriterTo.
func (s *String) WriteTo(w io.Writer) (int64, error) {
	data := []byte(*s)
	length := VarInt(len(data))

	n1, err := length.WriteTo(w)
	if err != nil {
		return n1, err
	}

	n2, err := w.Write(data)
	if err != nil {
		return n1 + int64(n2), err
	}

	return n1 + int64(n2), nil
}

var (
	_ io.WriterTo   = (*String)(nil)
	_ io.ReaderFrom = (*String)(nil)
)
