package packets_test

import (
	"testing"

	"github.com/nonya123456/cobble/proto/packets"
)

func TestState_IsValid(t *testing.T) {
	tests := []struct {
		name string
		s    packets.State
		want bool
	}{
		{
			name: "Handshaking",
			s:    0,
			want: true,
		},
		{
			name: "Status",
			s:    1,
			want: true,
		},
		{
			name: "Login",
			s:    2,
			want: true,
		},
		{
			name: "Transfer",
			s:    3,
			want: true,
		},
		{
			name: "Invalid",
			s:    4,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsValid(); got != tt.want {
				t.Errorf("State.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
