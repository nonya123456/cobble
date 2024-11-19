package types

import "io"

type String string

func (s *String) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64
	var length VarInt
	n1, err := length.ReadFrom(r)
	totalRead += n1
	if err != nil {
		return totalRead, err
	}

	data := make([]byte, length)
	n2, err := io.ReadFull(r, data)
	totalRead += int64(n2)
	if err != nil {
		return totalRead, err
	}

	*s = String(data)
	return totalRead, nil
}

func (s *String) WriteTo(w io.Writer) (int64, error) {
	var totalWrite int64
	data := []byte(*s)
	length := VarInt(len(data))

	n1, err := length.WriteTo(w)
	totalWrite += n1
	if err != nil {
		return totalWrite, err
	}

	n2, err := w.Write(data)
	totalWrite += int64(n2)
	if err != nil {
		return totalWrite, err
	}

	return totalWrite, nil
}
