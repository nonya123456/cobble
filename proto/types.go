package proto

import "io"

type VarInt int32

func (v VarInt) Bytes() []byte {
	var buffer []byte
	for {
		temp := byte(v & 0b01111111)
		v >>= 7
		if v != 0 {
			temp |= 0b10000000
		}
		buffer = append(buffer, temp)
		if v == 0 {
			break
		}
	}

	return buffer
}

func ReadVarInt(r io.Reader) (VarInt, error) {
	var result VarInt
	var shift uint
	for {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return 0, err
		}
		result |= VarInt(b[0]&0b01111111) << shift
		if b[0]&0b10000000 == 0 {
			break
		}
		shift += 7
	}

	return result, nil
}

type String string

func (s String) Bytes() []byte {
	data := []byte(s)
	length := VarInt(len(data)).Bytes()
	return append(length, data...)
}

func ReadString(r io.Reader) (String, error) {
	length, err := ReadVarInt(r)
	if err != nil {
		return "", err
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return "", err
	}

	return String(data), nil
}