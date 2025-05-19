package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunOperationTable(t *testing.T) {
	var tests = []struct {
		a, b, operator string
		want           float32
		err            error
	}{
		{"1", "2", "+", float32(3), nil},
		{"2", "1", "-", float32(1), nil},
		{"2", "1", "*", float32(2), nil},
		{"10", "2", "/", float32(5), nil},
		{"10", "2", "^", float32(100), nil},
		{"10", "0", "/", 0.0, fmt.Errorf("operation %v %v %s is not valid. Division by 0 is not allowed", 10, 0, "/")},
		{"1", "2", "_", 0.0, fmt.Errorf("operator %v not recognized", "_")},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v%v%v", tt.a, tt.operator, tt.b)

		t.Run(testname, func(t *testing.T) {
			res, err := runOperation(tt.a, tt.b, tt.operator)
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
				assert.Equal(t, tt.want, res)
				return
			}
			assert.Equal(t, tt.want, res)
			return
		})
	}
}

func BenchmarkIntMin(b *testing.B) {
	for b.Loop() {
		runOperation("10", "10", "^")
	}
}
