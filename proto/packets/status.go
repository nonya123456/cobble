package packets

import (
	"io"

	"github.com/nonya123456/cobble/proto/types"
)

const (
	StatusRequestID  = 0x00
	StatusResponseID = 0x00
	PingRequestID    = 0x01
	PingResponseID   = 0x01
)

type StatusRequest struct{}

func (s *StatusRequest) ReadFrom(r io.Reader) (int64, error) {
	return 0, nil
}

func (s *StatusRequest) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

type StatusResponse struct {
	JSONResponse string
}

func (s *StatusResponse) ReadFrom(r io.Reader) (int64, error) {
	var jsonResponse types.String
	n, err := readAll(r, &jsonResponse)
	if err != nil {
		return n, err
	}

	s.JSONResponse = string(jsonResponse)
	return n, nil
}

func (s *StatusResponse) WriteTo(w io.Writer) (int64, error) {
	jsonResponse := types.String(s.JSONResponse)
	return writeAll(w, &jsonResponse)
}

type PingRequest struct {
	Payload int64
}

func (p *PingRequest) ReadFrom(r io.Reader) (int64, error) {
	var payload types.Long
	n, err := readAll(r, &payload)
	if err != nil {
		return n, err
	}

	p.Payload = int64(payload)
	return n, nil
}

func (p *PingRequest) WriteTo(w io.Writer) (int64, error) {
	payload := types.Long(p.Payload)
	return writeAll(w, &payload)
}

type PingResponse struct {
	Payload int64
}

func (p *PingResponse) ReadFrom(r io.Reader) (int64, error) {
	var payload types.Long
	n, err := readAll(r, &payload)
	if err != nil {
		return n, err
	}

	p.Payload = int64(payload)
	return n, nil
}

func (p *PingResponse) WriteTo(w io.Writer) (int64, error) {
	payload := types.Long(p.Payload)
	return writeAll(w, &payload)
}
