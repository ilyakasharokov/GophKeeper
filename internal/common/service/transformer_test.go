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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProtoNotesToModels(tt.args.pN); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProtoNotesToModels() = %v, want %v", got, tt.want)
			}
		})
	}
}
