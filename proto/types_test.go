package proto_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/proto"
)

func TestVarInt_Bytes(t *testing.T) {
	tests := []struct {
		name string
		v    proto.VarInt
		want []byte
	}{
		{
			name: "Zero",
			v:    0,
			want: []byte{0x00},
		},
		{
			name: "Small positive number",
			v:    1,
			want: []byte{0x01},
		},
		{
			name: "Medium positive number",
			v:    300,
			want: []byte{0xAC, 0x02},
		},
		{
			name: "Large positive number",
			v:    2147483647,
			want: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x07},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VarInt.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadVarInt(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    proto.VarInt
		wantErr bool
	}{
		{
			name:    "Zero",
			args:    args{bytes.NewReader([]byte{0x00})},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Small positive number",
			args:    args{bytes.NewReader([]byte{0x01})},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Medium positive number",
			args:    args{bytes.NewReader([]byte{0xAC, 0x02})},
			want:    300,
			wantErr: false,
		},
		{
			name:    "Large positive number",
			args:    args{bytes.NewReader([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x07})},
			want:    2147483647,
			wantErr: false,
		},
		{
			name:    "Truncated VarInt",
			args:    args{bytes.NewReader([]byte{0xFF})},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := proto.ReadVarInt(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadVarInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadVarInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Bytes(t *testing.T) {
	tests := []struct {
		name string
		s    proto.String
		want []byte
	}{
		{
			name: "Empty string",
			s:    "",
			want: []byte{0x00},
		},
		{
			name: "Short string",
			s:    "hello",
			want: append([]byte{0x05}, []byte("hello")...),
		},
		{
			name: "Multibyte characters",
			s:    "你好",
			want: append([]byte{0x06}, []byte{0xE4, 0xBD, 0xA0, 0xE5, 0xA5, 0xBD}...),
		},
		{
			name: "Long string",
			s:    "This is a longer test string.",
			want: append([]byte{0x1D}, []byte("This is a longer test string.")...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadString(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    proto.String
		wantErr bool
	}{
		{
			name:    "Empty string",
			args:    args{bytes.NewReader([]byte{0x00})},
			want:    "",
			wantErr: false,
		},
		{
			name:    "Short string",
			args:    args{bytes.NewReader(append([]byte{0x05}, []byte("hello")...))},
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "Multibyte characters",
			args:    args{bytes.NewReader(append([]byte{0x06}, []byte{0xE4, 0xBD, 0xA0, 0xE5, 0xA5, 0xBD}...))},
			want:    "你好",
			wantErr: false,
		},
		{
			name:    "Long string",
			args:    args{bytes.NewReader(append([]byte{0x1D}, []byte("This is a longer test string.")...))},
			want:    "This is a longer test string.",
			wantErr: false,
		},
		{
			name:    "Truncated string",
			args:    args{bytes.NewReader([]byte{0x04, 0x74, 0x65, 0x73})},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid VarInt length",
			args:    args{bytes.NewReader([]byte{0x80, 0x80, 0x80, 0x80})},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := proto.ReadString(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLong_Bytes(t *testing.T) {
	tests := []struct {
		name string
		l    proto.Long
		want []byte
	}{
		{
			name: "Zero",
			l:    0,
			want: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "Small positive number",
			l:    1,
			want: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			name: "Large positive number",
			l:    9223372036854775807,
			want: []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name: "Small negative number",
			l:    -1,
			want: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name: "Large negative number",
			l:    -9223372036854775808,
			want: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Long.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadLong(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    proto.Long
		wantErr bool
	}{
		{
			name:    "Zero",
			args:    args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Small positive number",
			args:    args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Large positive number",
			args:    args{bytes.NewReader([]byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})},
			want:    9223372036854775807,
			wantErr: false,
		},
		{
			name:    "Small negative number",
			args:    args{bytes.NewReader([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})},
			want:    -1,
			wantErr: false,
		},
		{
			name:    "Large negative number",
			args:    args{bytes.NewReader([]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})},
			want:    -9223372036854775808,
			wantErr: false,
		},
		{
			name:    "Truncated long (less than 8 bytes)",
			args:    args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00})},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Empty reader",
			args:    args{bytes.NewReader([]byte{})},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := proto.ReadLong(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLong() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadLong() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnsignedShort_Bytes(t *testing.T) {
	tests := []struct {
		name string
		u    proto.UnsignedShort
		want []byte
	}{
		{
			name: "Minimum value",
			u:    0,
			want: []byte{0x00, 0x00},
		},
		{
			name: "Maximum value",
			u:    65535,
			want: []byte{0xFF, 0xFF},
		},
		{
			name: "Typical port",
			u:    25565,
			want: []byte{0x63, 0xDD},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnsignedShort.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadUnsignedShort(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    proto.UnsignedShort
		wantErr bool
	}{
		{
			name:    "Minimum value",
			args:    args{bytes.NewReader([]byte{0x00, 0x00})},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Maximum value",
			args:    args{bytes.NewReader([]byte{0xFF, 0xFF})},
			want:    65535,
			wantErr: false,
		},
		{
			name:    "Typical port value",
			args:    args{bytes.NewReader([]byte{0x63, 0xDD})},
			want:    25565,
			wantErr: false,
		},
		{
			name:    "Truncated data (only 1 byte)",
			args:    args{bytes.NewReader([]byte{0x63})},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Empty reader",
			args:    args{bytes.NewReader([]byte{})},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := proto.ReadUnsignedShort(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUnsignedShort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUnsignedShort() = %v, want %v", got, tt.want)
			}
		})
	}
}
