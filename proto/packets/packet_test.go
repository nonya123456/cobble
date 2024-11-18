package packets

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

type MockPacket struct {
	IDValue  int32
	Data     []byte
	ReadData []byte
}

func (p *MockPacket) ID() int32 {
	return p.IDValue
}

func (p *MockPacket) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(p.Data)
	return int64(n), err
}

func (p *MockPacket) ReadFrom(r io.Reader) (int64, error) {
	n, err := r.Read(p.ReadData)
	return int64(n), err
}

func TestWritePacket(t *testing.T) {
	type args struct {
		p Packet
	}
	tests := []struct {
		name    string
		args    args
		wantW   []byte
		wantErr bool
	}{
		{
			name:    "Simple packet",
			args:    args{&MockPacket{IDValue: 1, Data: []byte("Hello, Packet!")}},
			wantW:   append([]byte{0x0F, 0x01}, []byte("Hello, Packet!")...),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := WritePacket(w, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("WritePacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.Bytes(); !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("WritePacket() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestReadPacketID(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		{
			name:    "Simple packet",
			args:    args{bytes.NewReader(append([]byte{0x0F, 0x01}, []byte("Hello, Packet!")...))},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadPacketID(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPacketID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadPacketID() = %v, want %v", got, tt.want)
			}
		})
	}
}
