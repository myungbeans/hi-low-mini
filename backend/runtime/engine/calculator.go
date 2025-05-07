package engine

import (
	"fmt"
	"math"

	pb "hi_low_mini/gen/game_engine/v1"
	"hi_low_mini/runtime/models"
	"hi_low_mini/runtime/utils"
)

// Calculate reduces `cards` by running operations on (a, b), storing the result of each step in `scores`
func Calculate(cards []*pb.Card, scores []float32) ([]float32, error) {
	if len(scores) == countExpectedOps {
		return scores, nil
	}

	res, err := runOperation(cards[0].GetValue(), cards[2].GetValue(), cards[1].GetValue())
	if err != nil {
		return nil, err
	}

	scores = append(scores, res)
	// result becomes `a` for the next operation
	cards = append([]*pb.Card{models.NewNumCard(res)}, cards[3:]...)

	return Calculate(cards, scores)
}

// runOperation performs the operation (`a` `operator` `b`) for supported operators
func runOperation(strA, strB, operator string) (float32, error) {
	if len(strA+operator+strB) < 3 {
		return 0.0, fmt.Errorf("missing arguments in operation \"%v %v %v\"", strA, operator, strB)
	}

	a, err := utils.StrToFloat(strA)
	if err != nil {
		return 0.0, err
	}
	b, err := utils.StrToFloat(strB)
	if err != nil {
		return 0.0, err
	}

	switch operator {
	case models.OP_ADD:
		return a + b, nil
	case models.OP_SUBTRACT:
		return a - b, nil
	case models.OP_MULTIPLY:
		return a * b, nil
	case models.OP_DIVIDE:
		if b == float32(0) {
			return 0.0, fmt.Errorf("operation %v %v %s is not valid. Division by 0 is not allowed", a, b, operator)
		}
		return a / b, nil
	case models.OP_POWER:
		powBase := math.Pow(float64(a), float64(b))
		return float32(powBase), nil
	default:
		return 0.0, fmt.Errorf("operator %v not recognized", operator)
	}
}
