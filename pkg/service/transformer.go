package service

import (
	proto "gophkeeper/pkg/grpc/proto"
	"gophkeeper/pkg/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProtoNotesToModels(pN []*proto.Note) []models.Note {
	notes := make([]models.Note, 0)
	for _, gd := range pN {
		n := models.Note{
			ID:        gd.Id,
			Title:     gd.Title,
			Body:      gd.Body,
			Deleted:   gd.Deleted,
			CreatedAt: gd.CreatedAt.AsTime(),
			UpdatedAt: gd.UpdatedAt.AsTime(),
			DeletedAt: gd.DeletedAt.AsTime(),
		}
		notes = append(notes, n)
	}
	return notes
}

func NotesToProto(notes []models.Note) []*proto.Note {
	pNotes := make([]*proto.Note, 0)
	for _, d := range notes {
		pN := &proto.Note{
			Id:        d.ID,
			Title:     d.Title,
			Body:      d.Body,
			Deleted:   d.Deleted,
			CreatedAt: timestamppb.New(d.CreatedAt),
			UpdatedAt: timestamppb.New(d.UpdatedAt),
			DeletedAt: timestamppb.New(d.DeletedAt),
		}
		pNotes = append(pNotes, pN)
	}
	return pNotes
}
