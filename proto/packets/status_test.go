package packets_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/proto/packets"
)

func TestStatusRequest_ReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		s       *packets.StatusRequest
		args    args
		wantN   int64
		wantErr bool
	}{
		{
			name:    "Valid",
			s:       new(packets.StatusRequest),
			args:    args{bytes.NewReader([]byte{})},
			wantN:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &packets.StatusRequest{}
			gotN, err := s.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("StatusRequest.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("StatusRequest.ReadFrom() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestStatusRequest_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		s       *packets.StatusRequest
		wantN   int64
		wantW   []byte
		wantErr bool
	}{
		{
			name:    "Valid",
			s:       new(packets.StatusRequest),
			wantN:   0,
			wantW:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &packets.StatusRequest{}
			w := &bytes.Buffer{}
			gotN, err := s.WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("StatusRequest.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("StatusRequest.WriteTo() = %v, want %v", gotN, tt.wantN)
			}
			if gotW := w.Bytes(); !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("StatusRequest.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestStatusResponse_ReadFrom(t *testing.T) {
	type fields struct {
		JSONResponse string
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantN        int64
		wantErr      bool
		wantModified packets.StatusResponse
	}{
		{
			name:   "Valid JSON response",
			fields: fields{},
			args: args{bytes.NewReader(
				append([]byte{0x7D}, []byte(`{"version":{"name":"1.16.5","protocol":754},"players":{"max":100,"online":5},"description":{"text":"Welcome to the server!"}}`)...)),
			},
			wantN:        126,
			wantErr:      false,
			wantModified: packets.StatusResponse{JSONResponse: `{"version":{"name":"1.16.5","protocol":754},"players":{"max":100,"online":5},"description":{"text":"Welcome to the server!"}}`},
		},
		{
			name:         "Empty JSON response",
			fields:       fields{},
			args:         args{bytes.NewReader(append([]byte{0x02}, []byte(`{}`)...))},
			wantN:        3,
			wantErr:      false,
			wantModified: packets.StatusResponse{JSONResponse: `{}`},
		},
		{
			name:    "Invalid JSON response",
			fields:  fields{},
			args:    args{bytes.NewReader([]byte(`{invalid-json}`))},
			wantN:   14,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &packets.StatusResponse{
				JSONResponse: tt.fields.JSONResponse,
			}
			gotN, err := s.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("StatusResponse.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("StatusResponse.ReadFrom() = %v, want %v", gotN, tt.wantN)
			}
			if !reflect.DeepEqual(*s, tt.wantModified) {
				t.Errorf("StatusResponse.ReadFrom() s = %v, wantModified %v", *s, tt.wantModified)
			}
		})
	}
}

func TestStatusResponse_WriteTo(t *testing.T) {
	type fields struct {
		JSONResponse string
	}
	tests := []struct {
		name    string
		fields  fields
		wantN   int64
		wantW   []byte
		wantErr bool
	}{
		{
			name: "Valid JSON response",
			fields: fields{
				JSONResponse: `{"version":{"name":"1.16.5","protocol":754},"players":{"max":100,"online":5},"description":{"text":"Welcome to the server!"}}`,
			},
			wantN:   126,
			wantW:   append([]byte{0x7D}, []byte(`{"version":{"name":"1.16.5","protocol":754},"players":{"max":100,"online":5},"description":{"text":"Welcome to the server!"}}`)...),
			wantErr: false,
		},
		{
			name: "Empty JSON response",
			fields: fields{
				JSONResponse: `{}`,
			},
			wantN:   3,
			wantW:   append([]byte{0x02}, []byte(`{}`)...),
			wantErr: false,
		},
		{
			name: "Invalid JSON response",
			fields: fields{
				JSONResponse: "",
			},
			wantN:   1,
			wantW:   []byte{0x00},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &packets.StatusResponse{
				JSONResponse: tt.fields.JSONResponse,
			}
			w := &bytes.Buffer{}
			gotN, err := s.WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("StatusResponse.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("StatusResponse.WriteTo() = %v, want %v", gotN, tt.wantN)
			}
			if gotW := w.Bytes(); !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("StatusResponse.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPingRequest_ReadFrom(t *testing.T) {
	type fields struct {
		Payload int64
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantN        int64
		wantErr      bool
		wantModified packets.PingRequest
	}{
		{
			name:         "Valid payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})},
			wantN:        8,
			wantErr:      false,
			wantModified: packets.PingRequest{Payload: 1},
		},
		{
			name:         "Zero payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})},
			wantN:        8,
			wantErr:      false,
			wantModified: packets.PingRequest{Payload: 0},
		},
		{
			name:         "Negative payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})},
			wantN:        8,
			wantErr:      false,
			wantModified: packets.PingRequest{Payload: -1},
		},
		{
			name:         "Invalid data length",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00})},
			wantN:        4,
			wantErr:      true,
			wantModified: packets.PingRequest{},
		},
		{
			name:         "Empty data",
			fields:       fields{},
			args:         args{r: bytes.NewReader([]byte{})},
			wantN:        0,
			wantErr:      true,
			wantModified: packets.PingRequest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &packets.PingRequest{
				Payload: tt.fields.Payload,
			}
			gotN, err := p.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("PingRequest.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("PingRequest.ReadFrom() = %v, want %v", gotN, tt.wantN)
			}
			if !reflect.DeepEqual(*p, tt.wantModified) {
				t.Errorf("PingRequest.ReadFrom() p = %v, wantModified %v", *p, tt.wantModified)
			}
		})
	}
}

func TestPingRequest_WriteTo(t *testing.T) {
	type fields struct {
		Payload int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantN   int64
		wantW   []byte
		wantErr bool
	}{
		{
			name:    "Valid payload - positive",
			fields:  fields{Payload: 1},
			wantN:   8,
			wantW:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			wantErr: false,
		},
		{
			name:    "Valid payload - zero",
			fields:  fields{Payload: 0},
			wantN:   8,
			wantW:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			wantErr: false,
		},
		{
			name:    "Valid payload - negative",
			fields:  fields{Payload: -1},
			wantN:   8,
			wantW:   []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &packets.PingRequest{
				Payload: tt.fields.Payload,
			}
			w := &bytes.Buffer{}
			gotN, err := p.WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("PingRequest.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("PingRequest.WriteTo() = %v, want %v", gotN, tt.wantN)
			}
			if gotW := w.Bytes(); !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("PingRequest.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPingResponse_ReadFrom(t *testing.T) {
	type fields struct {
		Payload int64
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         int64
		wantErr      bool
		wantModified packets.PingResponse
	}{
		{
			name:         "Valid payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})},
			want:         8,
			wantErr:      false,
			wantModified: packets.PingResponse{Payload: 1},
		},
		{
			name:         "Zero payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})},
			want:         8,
			wantErr:      false,
			wantModified: packets.PingResponse{Payload: 0},
		},
		{
			name:         "Negative payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})},
			want:         8,
			wantErr:      false,
			wantModified: packets.PingResponse{Payload: -1},
		},
		{
			name:         "Invalid data length",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00})},
			want:         4,
			wantErr:      true,
			wantModified: packets.PingResponse{},
		},
		{
			name:         "Empty data",
			fields:       fields{},
			args:         args{r: bytes.NewReader([]byte{})},
			want:         0,
			wantErr:      true,
			wantModified: packets.PingResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &packets.PingResponse{
				Payload: tt.fields.Payload,
			}
			got, err := p.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("PingResponse.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PingResponse.ReadFrom() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(*p, tt.wantModified) {
				t.Errorf("PingResponse.ReadFrom() p = %v, wantModified %v", *p, tt.wantModified)
			}
		})
	}
}

func TestPingResponse_WriteTo(t *testing.T) {
	type fields struct {
		Payload int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantN   int64
		wantW   []byte
		wantErr bool
	}{
		{
			name:    "Valid payload - positive",
			fields:  fields{Payload: 1},
			wantN:   8,
			wantW:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			wantErr: false,
		},
		{
			name:    "Valid payload - zero",
			fields:  fields{Payload: 0},
			wantN:   8,
			wantW:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			wantErr: false,
		},
		{
			name:    "Valid payload - negative",
			fields:  fields{Payload: -1},
			wantN:   8,
			wantW:   []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &packets.PingResponse{
				Payload: tt.fields.Payload,
			}
			w := &bytes.Buffer{}
			got, err := p.WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("PingResponse.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantN {
				t.Errorf("PingResponse.WriteTo() = %v, want %v", got, tt.wantN)
			}
			if gotW := w.Bytes(); !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("PingResponse.WriteTo() got = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
