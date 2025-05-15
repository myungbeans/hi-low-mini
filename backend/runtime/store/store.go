package store

import (
	pb "hi_low_mini/gen/game_engine/v1"
)

type GamesStore interface {
	Get(id string) (*pb.Game, error)
	Put(game *pb.Game) error
}
