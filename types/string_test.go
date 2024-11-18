package types_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/types"
)

func TestString_ReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name         string
		s            *types.String
		args         args
		want         int64
		wantErr      bool
		wantModified string
	}{
		{
			name:         "Empty string",
			s:            new(types.String),
			args:         args{bytes.NewReader([]byte{0x00})},
			want:         1,
			wantErr:      false,
			wantModified: "",
		},
		{
			name:         "Short string",
			s:            new(types.String),
			args:         args{bytes.NewReader(append([]byte{0x05}, []byte("hello")...))},
			want:         6,
			wantErr:      false,
			wantModified: "hello",
		},
		{
			name:         "Multibyte characters",
			s:            new(types.String),
			args:         args{bytes.NewReader(append([]byte{0x06}, []byte{0xE4, 0xBD, 0xA0, 0xE5, 0xA5, 0xBD}...))},
			want:         7,
			wantErr:      false,
			wantModified: "你好",
		},
		{
			name:         "Long string",
			s:            new(types.String),
			args:         args{bytes.NewReader(append([]byte{0x1D}, []byte("This is a longer test string.")...))},
			want:         30,
			wantErr:      false,
			wantModified: "This is a longer test string.",
		},
		{
			name:         "Truncated string",
			s:            new(types.String),
			args:         args{bytes.NewReader([]byte{0x04, 0x74, 0x65, 0x73})},
			want:         4,
			wantErr:      true,
			wantModified: "",
		},
		{
			name:         "Invalid VarInt length",
			s:            new(types.String),
			args:         args{bytes.NewReader([]byte{0x80, 0x80, 0x80, 0x80})},
			want:         4,
			wantErr:      true,
			wantModified: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("String.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("String.ReadFrom() = %v, want %v", got, tt.want)
			}
			if string(*tt.s) != tt.wantModified {
				t.Errorf("VarInt.ReadFrom() modified v = %v, want %v", *tt.s, tt.wantModified)
			}
		})
	}
}

func TestString_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		s       *types.String
		want    int64
		wantW   []byte
		wantErr bool
	}{
		{
			name:    "Empty string",
			s:       newString(""),
			want:    1,
			wantW:   []byte{0x00},
			wantErr: false,
		},
		{
			name:    "Short string",
			s:       newString("hello"),
			want:    6,
			wantW:   append([]byte{0x05}, []byte("hello")...),
			wantErr: false,
		},
		{
			name:    "Multibyte characters",
			s:       newString("你好"),
			want:    7,
			wantW:   append([]byte{0x06}, []byte{0xE4, 0xBD, 0xA0, 0xE5, 0xA5, 0xBD}...),
			wantErr: false,
		},
		{
			name:    "Long string",
			s:       newString("This is a longer test string."),
			want:    30,
			wantW:   append([]byte{0x1D}, []byte("This is a longer test string.")...),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := tt.s.WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("String.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("String.WriteTo() = %v, want %v", got, tt.want)
			}
			if gotW := w.Bytes(); !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("String.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func newString(s string) *types.String {
	ts := types.String(s)
	return &ts
}
