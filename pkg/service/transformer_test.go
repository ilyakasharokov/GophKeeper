package service

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	proto "gophkeeper/pkg/grpc/proto"
	"gophkeeper/pkg/models"
	"reflect"
	"testing"
	"time"
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
		{name: "empty", args: args{notes: nil}, want: []*proto.Note{}},
		{name: "not empty", args: args{notes: []models.Note{{ID: "test"}}},
			want: []*proto.Note{
				{
					Id: "test", CreatedAt: timestamppb.New(time.Time{}),
					UpdatedAt: timestamppb.New(time.Time{}),
					DeletedAt: timestamppb.New(time.Time{}),
				},
			},
		},
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
	pn := proto.Note{Id: "test"}
	tests := []struct {
		name string
		args args
		want []models.Note
	}{
		{name: "ok", args: args{pN: nil}, want: []models.Note{}},
		{
			name: "not empty", args: args{
				pN: []*proto.Note{&pn},
			},
			want: []models.Note{{ID: "test", CreatedAt: (&timestamppb.Timestamp{}).AsTime(), DeletedAt: (&timestamppb.Timestamp{}).AsTime(), UpdatedAt: (&timestamppb.Timestamp{}).AsTime()}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProtoNotesToModels(tt.args.pN); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProtoNotesToModels() = %v, want %v", got, tt.want)
			}
		})
	}
}
