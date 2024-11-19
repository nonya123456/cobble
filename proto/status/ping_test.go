package status_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/proto/status"
)

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
		wantModified status.PingRequest
	}{
		{
			name:         "Valid payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})},
			wantN:        8,
			wantErr:      false,
			wantModified: status.PingRequest{Payload: 1},
		},
		{
			name:         "Zero payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})},
			wantN:        8,
			wantErr:      false,
			wantModified: status.PingRequest{Payload: 0},
		},
		{
			name:         "Negative payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})},
			wantN:        8,
			wantErr:      false,
			wantModified: status.PingRequest{Payload: -1},
		},
		{
			name:         "Invalid data length",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00})},
			wantN:        4,
			wantErr:      true,
			wantModified: status.PingRequest{},
		},
		{
			name:         "Empty data",
			fields:       fields{},
			args:         args{r: bytes.NewReader([]byte{})},
			wantN:        0,
			wantErr:      true,
			wantModified: status.PingRequest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &status.PingRequest{
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
			p := &status.PingRequest{
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
		wantModified status.PingResponse
	}{
		{
			name:         "Valid payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})},
			want:         8,
			wantErr:      false,
			wantModified: status.PingResponse{Payload: 1},
		},
		{
			name:         "Zero payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})},
			want:         8,
			wantErr:      false,
			wantModified: status.PingResponse{Payload: 0},
		},
		{
			name:         "Negative payload",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})},
			want:         8,
			wantErr:      false,
			wantModified: status.PingResponse{Payload: -1},
		},
		{
			name:         "Invalid data length",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00})},
			want:         4,
			wantErr:      true,
			wantModified: status.PingResponse{},
		},
		{
			name:         "Empty data",
			fields:       fields{},
			args:         args{r: bytes.NewReader([]byte{})},
			want:         0,
			wantErr:      true,
			wantModified: status.PingResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &status.PingResponse{
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
			p := &status.PingResponse{
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
