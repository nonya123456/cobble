package status

import (
	"io"

	"github.com/nonya123456/cobble/proto/stream"
	"github.com/nonya123456/cobble/proto/types"
)

const (
	StatusRequestID  = 0x00
	StatusResponseID = 0x00
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
	n, err := stream.ReadAll(r, &jsonResponse)
	if err != nil {
		return n, err
	}

	s.JSONResponse = string(jsonResponse)
	return n, nil
}

func (s *StatusResponse) WriteTo(w io.Writer) (int64, error) {
	jsonResponse := types.String(s.JSONResponse)
	return stream.WriteAll(w, &jsonResponse)
}
