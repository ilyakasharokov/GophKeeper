package service

import (
	"gophkeeper/internal/common/models"
	proto "gophkeeper/pkg/grpc/proto"
	"reflect"
	"testing"
)

func TestNotesToProto(t *testing.T) {
	type args struct {
		notes []models.Note
	}
	tests := []struct {
		name string
		args args
		want []*proto.Note
	}{
		{name: "ok", args: args{notes: nil}, want: []*proto.Note{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotesToProto(tt.args.notes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotesToProto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProtoNotesToModels(t *testing.T) {
	type args struct {
		pN []*proto.Note
	}
	tests := []struct {
		name string
		args args
		want []models.Note
	}{
		{name: "ok", args: args{pN: nil}, want: []models.Note{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProtoNotesToModels(tt.args.pN); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProtoNotesToModels() = %v, want %v", got, tt.want)
			}
		})
	}
}
