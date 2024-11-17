package proto_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/proto"
)

func TestPacket_Bytes(t *testing.T) {
	type fields struct {
		ID   proto.VarInt
		Data []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Empty packet",
			fields: fields{
				ID:   0,
				Data: []byte{},
			},
			want: []byte{0x00, 0x00},
		},
		{
			name: "Packet with data",
			fields: fields{
				ID:   1,
				Data: []byte{0x01, 0x02, 0x03},
			},
			want: []byte{0x03, 0x01, 0x01, 0x02, 0x03},
		},
		{
			name: "Packet with large ID and data",
			fields: fields{
				ID:   300,
				Data: []byte{0x0A, 0x0B},
			},
			want: []byte{0x02, 0xAC, 0x02, 0x0A, 0x0B},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := proto.Packet{
				ID:   tt.fields.ID,
				Data: tt.fields.Data,
			}
			if got := p.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Packet.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadPacket(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    proto.Packet
		wantErr bool
	}{
		{
			name: "Empty packet",
			args: args{
				r: bytes.NewReader([]byte{0x00, 0x00}),
			},
			want: proto.Packet{
				ID:   0,
				Data: []byte{},
			},
			wantErr: false,
		},
		{
			name: "Packet with data",
			args: args{
				r: bytes.NewReader([]byte{0x03, 0x01, 0x01, 0x02, 0x03}),
			},
			want: proto.Packet{
				ID:   1,
				Data: []byte{0x01, 0x02, 0x03},
			},
			wantErr: false,
		},
		{
			name: "Packet with large ID and data",
			args: args{
				r: bytes.NewReader([]byte{0x02, 0xAC, 0x02, 0x0A, 0x0B}),
			},
			want: proto.Packet{
				ID:   300,
				Data: []byte{0x0A, 0x0B},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := proto.ReadPacket(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}
