package types_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/proto/types"
)

func TestUnsignedShort_ReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name         string
		u            *types.UnsignedShort
		args         args
		want         int64
		wantErr      bool
		wantModified uint16
	}{
		{
			name:         "Minimum value",
			u:            new(types.UnsignedShort),
			args:         args{bytes.NewReader([]byte{0x00, 0x00})},
			want:         2,
			wantErr:      false,
			wantModified: 0,
		},
		{
			name:         "Maximum value",
			u:            new(types.UnsignedShort),
			args:         args{bytes.NewReader([]byte{0xFF, 0xFF})},
			want:         2,
			wantErr:      false,
			wantModified: 65535,
		},
		{
			name:         "Typical port value",
			u:            new(types.UnsignedShort),
			args:         args{bytes.NewReader([]byte{0x63, 0xDD})},
			want:         2,
			wantErr:      false,
			wantModified: 25565,
		},
		{
			name:         "Truncated data (only 1 byte)",
			u:            new(types.UnsignedShort),
			args:         args{bytes.NewReader([]byte{0x63})},
			want:         1,
			wantErr:      true,
			wantModified: 0,
		},
		{
			name:         "Empty reader",
			u:            new(types.UnsignedShort),
			args:         args{bytes.NewReader([]byte{})},
			want:         0,
			wantErr:      true,
			wantModified: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnsignedShort.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UnsignedShort.ReadFrom() = %v, want %v", got, tt.want)
			}
			if uint16(*tt.u) != tt.wantModified {
				t.Errorf("VarInt.ReadFrom() modified v = %v, want %v", *tt.u, tt.wantModified)
			}
		})
	}
}

func TestUnsignedShort_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		u       *types.UnsignedShort
		want    int64
		wantW   []byte
		wantErr bool
	}{
		{
			name:    "Minimum value",
			u:       newUnsignedShort(0),
			want:    2,
			wantW:   []byte{0x00, 0x00},
			wantErr: false,
		},
		{
			name:    "Maximum value",
			u:       newUnsignedShort(65535),
			want:    2,
			wantW:   []byte{0xFF, 0xFF},
			wantErr: false,
		},
		{
			name:    "Typical port",
			u:       newUnsignedShort(25565),
			want:    2,
			wantW:   []byte{0x63, 0xDD},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := tt.u.WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnsignedShort.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UnsignedShort.WriteTo() = %v, want %v", got, tt.want)
			}
			if gotW := w.Bytes(); !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("UnsignedShort.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func newUnsignedShort(i uint16) *types.UnsignedShort {
	u := types.UnsignedShort(i)
	return &u
}
