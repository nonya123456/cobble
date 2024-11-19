package handshaking

import (
	"io"

	"github.com/nonya123456/cobble/proto/stream"
	"github.com/nonya123456/cobble/proto/types"
)

const HandshakeID int32 = 0x00

type Handshake struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       int32
}

func (h *Handshake) ReadFrom(r io.Reader) (int64, error) {
	var protocolVersion types.VarInt
	var serverAddress types.String
	var serverPort types.UnsignedShort
	var nextState types.VarInt

	n, err := stream.ReadAll(r, &protocolVersion, &serverAddress, &serverPort, &nextState)
	if err != nil {
		return n, err
	}

	h.ProtocolVersion = int32(protocolVersion)
	h.ServerAddress = string(serverAddress)
	h.ServerPort = uint16(serverPort)
	h.NextState = int32(nextState)
	return n, nil
}

func (h *Handshake) WriteTo(w io.Writer) (int64, error) {
	protocolVersion := types.VarInt(h.ProtocolVersion)
	serverAddress := types.String(h.ServerAddress)
	serverPort := types.UnsignedShort(h.ServerPort)
	nextState := types.VarInt(h.NextState)
	return stream.WriteAll(w, &protocolVersion, &serverAddress, &serverPort, &nextState)
}
