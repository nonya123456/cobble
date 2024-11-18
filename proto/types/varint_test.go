package types_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/proto/types"
)

func TestVarInt_ReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name         string
		v            *types.VarInt
		args         args
		want         int64
		wantErr      bool
		wantModified int32
	}{
		{
			name:         "Zero",
			v:            new(types.VarInt),
			args:         args{bytes.NewReader([]byte{0x00})},
			want:         1,
			wantErr:      false,
			wantModified: 0,
		},
		{
			name:         "Small positive number",
			v:            new(types.VarInt),
			args:         args{bytes.NewReader([]byte{0x01})},
			want:         1,
			wantErr:      false,
			wantModified: 1,
		},
		{
			name:         "Medium positive number",
			v:            new(types.VarInt),
			args:         args{bytes.NewReader([]byte{0xAC, 0x02})},
			want:         2,
			wantErr:      false,
			wantModified: 300,
		},
		{
			name:         "Large positive number",
			v:            new(types.VarInt),
			args:         args{bytes.NewReader([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x07})},
			want:         5,
			wantErr:      false,
			wantModified: 2147483647,
		},
		{
			name:         "Truncated VarInt",
			v:            new(types.VarInt),
			args:         args{bytes.NewReader([]byte{0xFF})},
			want:         1,
			wantErr:      true,
			wantModified: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("VarInt.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VarInt.ReadFrom() = %v, want %v", got, tt.want)
			}
			if int32(*tt.v) != tt.wantModified {
				t.Errorf("VarInt.ReadFrom() modified v = %v, want %v", *tt.v, tt.wantModified)
			}
		})
	}
}

func TestVarInt_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		v       *types.VarInt
		want    int64
		wantR   []byte
		wantErr bool
	}{
		{
			name:    "Zero",
			v:       newVarInt(0),
			want:    1,
			wantR:   []byte{0x00},
			wantErr: false,
		},
		{
			name:    "Small positive number",
			v:       newVarInt(1),
			want:    1,
			wantR:   []byte{0x01},
			wantErr: false,
		},
		{
			name:    "Medium positive number",
			v:       newVarInt(300),
			want:    2,
			wantR:   []byte{0xAC, 0x02},
			wantErr: false,
		},
		{
			name:    "Large positive number",
			v:       newVarInt(2147483647),
			want:    5,
			wantR:   []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x07},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &bytes.Buffer{}
			got, err := tt.v.WriteTo(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("VarInt.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VarInt.WriteTo() = %v, want %v", got, tt.want)
			}
			if gotR := r.Bytes(); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("VarInt.WriteTo() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func newVarInt(i int32) *types.VarInt {
	v := types.VarInt(i)
	return &v
}
