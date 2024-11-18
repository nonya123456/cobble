package packets

import (
	"errors"
	"io"

	"github.com/nonya123456/cobble/proto/types"
)

const HandshakeID int32 = 0x00

type Handshake struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       State
}

type State int32

const (
	StateHandshaking State = iota
	StateStatus
	StateLogin
	StateTransfer
)

var allStates = map[State]struct{}{StateHandshaking: {}, StateStatus: {}, StateLogin: {}, StateTransfer: {}}

var (
	ErrInvalidState = errors.New("invalid state")
)

func (h *Handshake) ReadFrom(r io.Reader) (int64, error) {
	var protocolVersion types.VarInt
	var serverAddress types.String
	var serverPort types.UnsignedShort
	var nextState types.VarInt

	n, err := readAll(r, &protocolVersion, &serverAddress, &serverPort, &nextState)
	if err != nil {
		return n, err
	}

	if _, ok := allStates[State(nextState)]; !ok {
		return n, ErrInvalidState
	}

	h.ProtocolVersion = int32(protocolVersion)
	h.ServerAddress = string(serverAddress)
	h.ServerPort = uint16(serverPort)
	h.NextState = State(nextState)
	return n, nil
}

func (h *Handshake) WriteTo(w io.Writer) (int64, error) {
	protocolVersion := types.VarInt(h.ProtocolVersion)
	serverAddress := types.String(h.ServerAddress)
	serverPort := types.UnsignedShort(h.ServerPort)
	nextState := types.VarInt(h.NextState)
	return writeAll(w, &protocolVersion, &serverAddress, &serverPort, &nextState)
}
