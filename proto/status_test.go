package proto_test

import (
	"reflect"
	"testing"

	"github.com/nonya123456/cobble/proto"
)

func TestPingRequest_Packet(t *testing.T) {
	type fields struct {
		Payload proto.Long
	}
	tests := []struct {
		name   string
		fields fields
		want   proto.Packet
	}{
		{
			name: "Valid packet",
			fields: fields{
				Payload: 0,
			},
			want: proto.Packet{
				ID:   0x01,
				Data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := proto.PingRequest{
				Payload: tt.fields.Payload,
			}
			if got := p.Packet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PingRequest.Packet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnMarshalPingRequest(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    proto.PingRequest
		wantErr bool
	}{
		{
			name: "Valid payload",
			args: args{
				data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			},
			want: proto.PingRequest{
				Payload: proto.Long(1),
			},
			wantErr: false,
		},
		{
			name: "Zero payload",
			args: args{
				data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
			want: proto.PingRequest{
				Payload: proto.Long(0),
			},
			wantErr: false,
		},
		{
			name: "Negative payload",
			args: args{
				data: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			},
			want: proto.PingRequest{
				Payload: proto.Long(-1),
			},
			wantErr: false,
		},
		{
			name: "Invalid data length",
			args: args{
				data: []byte{0x00, 0x00, 0x00, 0x00},
			},
			want:    proto.PingRequest{},
			wantErr: true,
		},
		{
			name: "Empty data",
			args: args{
				data: []byte{},
			},
			want:    proto.PingRequest{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := proto.UnMarshalPingRequest(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnMarshalPingRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnMarshalPingRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPingResponse_Packet(t *testing.T) {
	type fields struct {
		Payload proto.Long
	}
	tests := []struct {
		name   string
		fields fields
		want   proto.Packet
	}{
		{
			name: "Valid packet",
			fields: fields{
				Payload: 0,
			},
			want: proto.Packet{
				ID:   0x01,
				Data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := proto.PingResponse{
				Payload: tt.fields.Payload,
			}
			if got := p.Packet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PingResponse.Packet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusRequest_Packet(t *testing.T) {
	tests := []struct {
		name string
		s    proto.StatusRequest
		want proto.Packet
	}{
		{
			name: "Valid status request",
			s:    proto.StatusRequest{},
			want: proto.Packet{
				ID:   0x00,
				Data: []byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := proto.StatusRequest{}
			if got := s.Packet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusRequest.Packet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusResponse_Packet(t *testing.T) {
	type fields struct {
		JSONResponse proto.String
	}
	tests := []struct {
		name   string
		fields fields
		want   proto.Packet
	}{
		{
			name: "Valid JSON response",
			fields: fields{
				JSONResponse: `{"version":{"name":"1.16.5","protocol":754},"players":{"max":100,"online":5},"description":{"text":"Welcome to the server!"}}`,
			},
			want: proto.Packet{
				ID: 0x00,
				Data: []byte{
					0x7D,
					0x7B, 0x22, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6F, 0x6E, 0x22, 0x3A, 0x7B, 0x22, 0x6E, 0x61, 0x6D,
					0x65, 0x22, 0x3A, 0x22, 0x31, 0x2E, 0x31, 0x36, 0x2E, 0x35, 0x22, 0x2C, 0x22, 0x70, 0x72, 0x6F,
					0x74, 0x6F, 0x63, 0x6F, 0x6C, 0x22, 0x3A, 0x37, 0x35, 0x34, 0x7D, 0x2C, 0x22, 0x70, 0x6C, 0x61,
					0x79, 0x65, 0x72, 0x73, 0x22, 0x3A, 0x7B, 0x22, 0x6D, 0x61, 0x78, 0x22, 0x3A, 0x31, 0x30, 0x30,
					0x2C, 0x22, 0x6F, 0x6E, 0x6C, 0x69, 0x6E, 0x65, 0x22, 0x3A, 0x35, 0x7D, 0x2C, 0x22, 0x64, 0x65,
					0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6F, 0x6E, 0x22, 0x3A, 0x7B, 0x22, 0x74, 0x65, 0x78,
					0x74, 0x22, 0x3A, 0x22, 0x57, 0x65, 0x6C, 0x63, 0x6F, 0x6D, 0x65, 0x20, 0x74, 0x6F, 0x20, 0x74,
					0x68, 0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x21, 0x22, 0x7D, 0x7D,
				},
			},
		},
		{
			name: "Empty JSON response",
			fields: fields{
				JSONResponse: `{}`,
			},
			want: proto.Packet{
				ID: 0x00,
				Data: []byte{
					0x02,
					0x7B, 0x7D,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := proto.StatusResponse{
				JSONResponse: tt.fields.JSONResponse,
			}
			if got := s.Packet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusResponse.Packet() = %v, want %v", got, tt.want)
			}
		})
	}
}
