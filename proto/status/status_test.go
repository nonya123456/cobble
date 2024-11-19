package status_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/proto/status"
)

func TestStatusRequest_ReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		s       *status.StatusRequest
		args    args
		wantN   int64
		wantErr bool
	}{
		{
			name:    "Valid",
			s:       new(status.StatusRequest),
			args:    args{bytes.NewReader([]byte{})},
			wantN:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &status.StatusRequest{}
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
		s       *status.StatusRequest
		wantN   int64
		wantW   []byte
		wantErr bool
	}{
		{
			name:    "Valid",
			s:       new(status.StatusRequest),
			wantN:   0,
			wantW:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &status.StatusRequest{}
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
		wantModified status.StatusResponse
	}{
		{
			name:   "Valid JSON response",
			fields: fields{},
			args: args{bytes.NewReader(
				append([]byte{0x7D}, []byte(`{"version":{"name":"1.16.5","protocol":754},"players":{"max":100,"online":5},"description":{"text":"Welcome to the server!"}}`)...)),
			},
			wantN:        126,
			wantErr:      false,
			wantModified: status.StatusResponse{JSONResponse: `{"version":{"name":"1.16.5","protocol":754},"players":{"max":100,"online":5},"description":{"text":"Welcome to the server!"}}`},
		},
		{
			name:         "Empty JSON response",
			fields:       fields{},
			args:         args{bytes.NewReader(append([]byte{0x02}, []byte(`{}`)...))},
			wantN:        3,
			wantErr:      false,
			wantModified: status.StatusResponse{JSONResponse: `{}`},
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
			s := &status.StatusResponse{
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
			s := &status.StatusResponse{
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
