package cobble

import "io"

type Packet struct {
	Length int32
	ID     int32
	Data   []byte
}

func ReadPacket(r io.Reader) (Packet, error) {
	length, err := readVarInt(r)
	if err != nil {
		return Packet{}, err
	}

	id, err := readVarInt(r)
	if err != nil {
		return Packet{}, err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return Packet{}, err
	}

	return Packet{
		Length: length,
		ID:     id,
		Data:   data,
	}, nil
}

func readVarInt(r io.Reader) (int32, error) {
	var result int32
	var shift uint
	for {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return 0, err
		}
		result |= int32(b[0]&0b01111111) << shift
		if b[0]&0b10000000 == 0 {
			break
		}
		shift += 7
	}
	return result, nil
}

func readString(r io.Reader) (string, error) {
	length, err := readVarInt(r)
	if err != nil {
		return "", err
	}
	data := make([]byte, length)
	if _, err := r.Read(data); err != nil {
		return "", err
	}
	return string(data), nil
}
