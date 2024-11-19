package proto_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/proto"
)

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
			name:    "Simple packet",
			args:    args{bytes.NewReader(append([]byte{0x0F, 0x01}, []byte("Hello, Packet!")...))},
			want:    proto.Packet{ID: 1, Data: []byte("Hello, Packet!")},
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

func TestWritePacket(t *testing.T) {
	type args struct {
		id int32
		p  io.WriterTo
	}
	tests := []struct {
		name    string
		args    args
		wantW   []byte
		wantErr bool
	}{
		{
			name: "Simple packet",
			args: args{
				id: 1,
				p:  bytes.NewBuffer([]byte("Hello, Packet!")),
			},
			wantW:   append([]byte{0x0F, 0x01}, []byte("Hello, Packet!")...),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := proto.WritePacket(w, tt.args.id, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("WritePacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.Bytes(); !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("WritePacket() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
