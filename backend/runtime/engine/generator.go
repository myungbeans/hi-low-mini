package engine

import (
	"math"
	"math/rand/v2"

	pb "hi_low_mini/gen/game_engine/v1"
	"hi_low_mini/runtime/models"
)

func randomNum() float32 {
	num := math.Max(float64(rand.IntN(10)), float64(1))
	return float32(num)
}

func randomOperator() string {
	idx := rand.IntN(len(models.OPERATORS))
	return models.OPERATORS[idx]
}

func GenerateNumCards(n int64) []*pb.Card {
	tmp := map[float32]float32{}
	cards := []*pb.Card{}

	for _ = range n {
		num := randomNum()
		// Reroll on dupes - still possible to get dupes, but more unlikely
		if _, ok := tmp[num]; ok {
			num = randomNum()
		}
		cards = append(cards, models.NewNumCard(num))
	}
	return cards
}

func GenerateOpCards(n int64) []*pb.Card {
	cards := []*pb.Card{}

	for _ = range n {
		cards = append(cards, models.NewOpCard(randomOperator()))
	}
	return cards
}

func NewSet() []*pb.Card {
	cards := []*pb.Card{}

	cards = append(cards, GenerateNumCards(countExpectedNums)...)
	cards = append(cards, GenerateOpCards(countExpectedOps)...)

	return cards
}
