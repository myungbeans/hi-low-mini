package models

import (
	pb "hi_low_mini/gen/game_engine/v1"
	"hi_low_mini/runtime/utils"
)

// NewNumCard creates a pb.Card from a float32 value
func NewNumCard(val float32) *pb.Card {
	return &pb.Card{
		Type:  pb.CardType_CARD_TYPE_NUMBER,
		Value: utils.FloatToStr(val),
	}
}

// NewOpCard creates a pb.Card from an operator str value
func NewOpCard(op string) *pb.Card {
	return &pb.Card{
		Type:  pb.CardType_CARD_TYPE_OPERATOR,
		Value: op,
	}
}
