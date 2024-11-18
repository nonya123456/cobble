package packets_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/proto/packets"
)

func TestHandshake_ReadFrom(t *testing.T) {
	type fields struct {
		ProtocolVersion int32
		ServerAddress   string
		ServerPort      uint16
		NextState       packets.State
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
		wantModified packets.Handshake
	}{
		{
			name:         "Valid handshake (status request)",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x04, 0x09, 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't', 0x1F, 0x90, 0x01})},
			want:         14,
			wantErr:      false,
			wantModified: packets.Handshake{ProtocolVersion: 4, ServerAddress: "localhost", ServerPort: 8080, NextState: 1},
		},
		{
			name:         "Valid handshake (login request)",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x04, 0x09, 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't', 0x1F, 0x90, 0x02})},
			want:         14,
			wantErr:      false,
			wantModified: packets.Handshake{ProtocolVersion: 4, ServerAddress: "localhost", ServerPort: 8080, NextState: 2},
		},
		{
			name:         "Invalid handshake (truncated data)",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0x04, 0x09})},
			want:         2,
			wantErr:      true,
			wantModified: packets.Handshake{},
		},
		{
			name:         "Invalid handshake (malformed VarInt)",
			fields:       fields{},
			args:         args{bytes.NewReader([]byte{0xFF})},
			want:         1,
			wantErr:      true,
			wantModified: packets.Handshake{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &packets.Handshake{
				ProtocolVersion: tt.fields.ProtocolVersion,
				ServerAddress:   tt.fields.ServerAddress,
				ServerPort:      tt.fields.ServerPort,
				NextState:       tt.fields.NextState,
			}
			got, err := h.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handshake.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Handshake.ReadFrom() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(*h, tt.wantModified) {
				t.Errorf("Long.WriteTo() = %v, want %v", *h, tt.wantModified)
			}
		})
	}
}

func TestHandshake_WriteTo(t *testing.T) {
	type fields struct {
		ProtocolVersion int32
		ServerAddress   string
		ServerPort      uint16
		NextState       packets.State
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantW   []byte
		wantErr bool
	}{
		{
			name:    "Valid handshake (status request)",
			fields:  fields{ProtocolVersion: 4, ServerAddress: "localhost", ServerPort: 8080, NextState: 1},
			want:    14,
			wantW:   []byte{0x04, 0x09, 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't', 0x1F, 0x90, 0x01},
			wantErr: false,
		},
		{
			name:    "Valid handshake (login request)",
			fields:  fields{ProtocolVersion: 754, ServerAddress: "example.com", ServerPort: 25565, NextState: 2},
			want:    17,
			wantW:   []byte{0xF2, 0x05, 0x0B, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm', 0x63, 0xDD, 0x02},
			wantErr: false,
		},
		{
			name:    "Empty server address",
			fields:  fields{ProtocolVersion: 1, ServerAddress: "", ServerPort: 0, NextState: 1},
			want:    5,
			wantW:   []byte{0x01, 0x00, 0x00, 0x00, 0x01},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &packets.Handshake{
				ProtocolVersion: tt.fields.ProtocolVersion,
				ServerAddress:   tt.fields.ServerAddress,
				ServerPort:      tt.fields.ServerPort,
				NextState:       tt.fields.NextState,
			}
			w := &bytes.Buffer{}
			got, err := h.WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handshake.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Handshake.WriteTo() = %v, want %v", got, tt.want)
			}
			if gotW := w.Bytes(); !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("Handshake.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
