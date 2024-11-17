package proto

import (
	"reflect"
	"testing"
)

func TestHandshake_Packet(t *testing.T) {
	type fields struct {
		ProtocolVersion VarInt
		ServerAddress   String
		ServerPort      UnsignedShort
		NextState       State
	}
	tests := []struct {
		name   string
		fields fields
		want   Packet
	}{
		{
			name: "Valid handshake (status request)",
			fields: fields{
				ProtocolVersion: 4,
				ServerAddress:   "localhost",
				ServerPort:      8080,
				NextState:       1,
			},
			want: Packet{
				ID: 0x00,
				Data: []byte{
					0x04,
					0x09, 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't',
					0x1F, 0x90,
					0x01,
				},
			},
		},
		{
			name: "Valid handshake (login request)",
			fields: fields{
				ProtocolVersion: 754,
				ServerAddress:   "example.com",
				ServerPort:      25565,
				NextState:       2,
			},
			want: Packet{
				ID: 0x00,
				Data: []byte{
					0xF2, 0x05,
					0x0B, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
					0x63, 0xDD,
					0x02,
				},
			},
		},
		{
			name: "Empty server address",
			fields: fields{
				ProtocolVersion: 1,
				ServerAddress:   "",
				ServerPort:      0,
				NextState:       1,
			},
			want: Packet{
				ID: 0x00,
				Data: []byte{
					0x01,
					0x00,
					0x00, 0x00,
					0x01,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handshake{
				ProtocolVersion: tt.fields.ProtocolVersion,
				ServerAddress:   tt.fields.ServerAddress,
				ServerPort:      tt.fields.ServerPort,
				NextState:       tt.fields.NextState,
			}
			if got := h.Packet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handshake.Packet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalHandshake(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Handshake
		wantErr bool
	}{
		{
			name: "Valid handshake (status request)",
			args: args{
				data: []byte{
					0x04,
					0x09, 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't',
					0x1F, 0x90,
					0x01,
				},
			},
			want: Handshake{
				ProtocolVersion: 4,
				ServerAddress:   "localhost",
				ServerPort:      8080,
				NextState:       1,
			},
			wantErr: false,
		},
		{
			name: "Valid handshake (login request)",
			args: args{
				data: []byte{
					0x04,
					0x09, 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't',
					0x1F, 0x90,
					0x02,
				},
			},
			want: Handshake{
				ProtocolVersion: 4,
				ServerAddress:   "localhost",
				ServerPort:      8080,
				NextState:       2,
			},
			wantErr: false,
		},
		{
			name: "Invalid handshake (truncated data)",
			args: args{
				data: []byte{0x04, 0x09},
			},
			want:    Handshake{},
			wantErr: true,
		},
		{
			name: "Invalid handshake (malformed VarInt)",
			args: args{
				data: []byte{0xFF},
			},
			want:    Handshake{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalHandshake(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalHandshake() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalHandshake() = %v, want %v", got, tt.want)
			}
		})
	}
}
