package zigzag_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/dnaka91/stef/internal/zigzag"
	"github.com/stretchr/testify/assert"
)

func TestI16(t *testing.T) {
	for _, value := range []int16{0, 1, math.MinInt16, math.MaxInt16} {
		assert.Equal(t, value, zigzag.DecodeI16(zigzag.EncodeI16(value)))
	}
}

func TestI32(t *testing.T) {
	for _, value := range []int32{0, 1, math.MinInt32, math.MaxInt32} {
		assert.Equal(t, value, zigzag.DecodeI32(zigzag.EncodeI32(value)))
	}
}

func TestI64(t *testing.T) {
	for _, value := range []int64{0, 1, math.MinInt64, math.MaxInt64} {
		assert.Equal(t, value, zigzag.DecodeI64(zigzag.EncodeI64(value)))
	}
}

func TestI128(t *testing.T) {
	for _, value := range []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
		new(big.Int).Exp(big.NewInt(-2), big.NewInt(127), nil),
		new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(127), nil), big.NewInt(1)),
	} {
		assert.Equal(t, value, zigzag.DecodeI128(zigzag.EncodeI128(value)))
	}
}
