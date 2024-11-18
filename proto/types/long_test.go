package types_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/proto/types"
)

func TestLong_ReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name         string
		l            *types.Long
		args         args
		want         int64
		wantErr      bool
		wantModified int64
	}{
		{
			name:         "Zero",
			l:            new(types.Long),
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})},
			want:         8,
			wantErr:      false,
			wantModified: 0,
		},
		{
			name:         "Small positive number",
			l:            new(types.Long),
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})},
			want:         8,
			wantErr:      false,
			wantModified: 1,
		},
		{
			name:         "Large positive number",
			l:            new(types.Long),
			args:         args{bytes.NewReader([]byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})},
			want:         8,
			wantErr:      false,
			wantModified: 9223372036854775807,
		},
		{
			name:         "Small negative number",
			l:            new(types.Long),
			args:         args{bytes.NewReader([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})},
			want:         8,
			wantErr:      false,
			wantModified: -1,
		},
		{
			name:         "Large negative number",
			l:            new(types.Long),
			args:         args{bytes.NewReader([]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})},
			want:         8,
			wantErr:      false,
			wantModified: -9223372036854775808,
		},
		{
			name:         "Truncated long (less than 8 bytes)",
			l:            new(types.Long),
			args:         args{bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00})},
			want:         4,
			wantErr:      true,
			wantModified: 0,
		},
		{
			name:         "Empty reader",
			l:            new(types.Long),
			args:         args{bytes.NewReader([]byte{})},
			want:         0,
			wantErr:      true,
			wantModified: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Long.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Long.ReadFrom() = %v, want %v", got, tt.want)
			}
			if int64(*tt.l) != tt.wantModified {
				t.Errorf("VarInt.ReadFrom() modified v = %v, want %v", *tt.l, tt.wantModified)
			}
		})
	}
}

func TestLong_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		l       *types.Long
		want    int64
		wantW   []byte
		wantErr bool
	}{
		{
			name:  "Zero",
			l:     newLong(0),
			want:  8,
			wantW: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:  "Small positive number",
			l:     newLong(1),
			want:  8,
			wantW: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			name:  "Large positive number",
			l:     newLong(9223372036854775807),
			want:  8,
			wantW: []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:  "Small negative number",
			l:     newLong(-1),
			want:  8,
			wantW: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:  "Large negative number",
			l:     newLong(-9223372036854775808),
			want:  8,
			wantW: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := tt.l.WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("Long.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Long.WriteTo() = %v, want %v", got, tt.want)
			}
			if gotW := w.Bytes(); !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("Long.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func newLong(i int64) *types.Long {
	l := types.Long(i)
	return &l
}
