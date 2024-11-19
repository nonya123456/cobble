package status

import (
	"io"

	"github.com/nonya123456/cobble/proto/stream"
	"github.com/nonya123456/cobble/proto/types"
)

const (
	PingRequestID  = 0x01
	PingResponseID = 0x01
)

type PingRequest struct {
	Payload int64
}

func (p *PingRequest) ReadFrom(r io.Reader) (int64, error) {
	var payload types.Long
	n, err := stream.ReadAll(r, &payload)
	if err != nil {
		return n, err
	}

	p.Payload = int64(payload)
	return n, nil
}

func (p *PingRequest) WriteTo(w io.Writer) (int64, error) {
	payload := types.Long(p.Payload)
	return stream.WriteAll(w, &payload)
}

type PingResponse struct {
	Payload int64
}

func (p *PingResponse) ReadFrom(r io.Reader) (int64, error) {
	var payload types.Long
	n, err := stream.ReadAll(r, &payload)
	if err != nil {
		return n, err
	}

	p.Payload = int64(payload)
	return n, nil
}

func (p *PingResponse) WriteTo(w io.Writer) (int64, error) {
	payload := types.Long(p.Payload)
	return stream.WriteAll(w, &payload)
}
