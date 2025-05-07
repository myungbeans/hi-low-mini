package models

import (
	pb "hi_low_mini/gen/game_engine/v1"
	"strconv"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newID(now *timestamppb.Timestamp) string {
	ts := strconv.Itoa(int(now.Nanos))
	return ts + "-" + uuid.New().String()
}

func NewGame(cards []*pb.Card) *pb.Game {
	now := timestamppb.Now()
	return &pb.Game{
		Id:        newID(now),
		Timestamp: now,
		Pool: &pb.Pool{
			Cards: cards,
		},
	}
}
