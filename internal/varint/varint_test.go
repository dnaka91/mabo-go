package varint_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/dnaka91/mabo/internal/varint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestU16(t *testing.T) {
	for _, value := range []uint16{0, 1, math.MaxUint16} {
		buf := varint.EncodeU16(value)
		got, _, err := varint.DecodeU16(buf)
		require.NoError(t, err)
		assert.Equal(t, value, got)
	}
}

func TestU32(t *testing.T) {
	for _, value := range []uint32{0, 1, math.MaxUint32} {
		buf := varint.EncodeU32(value)
		got, _, err := varint.DecodeU32(buf)
		require.NoError(t, err)
		assert.Equal(t, value, got)
	}
}

func TestU64(t *testing.T) {
	for _, value := range []uint64{0, 1, math.MaxUint64} {
		buf := varint.EncodeU64(value)
		got, _, err := varint.DecodeU64(buf)
		require.NoError(t, err)
		assert.Equal(t, value, got)
	}
}

func TestU128(t *testing.T) {
	for _, value := range []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
		new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil), big.NewInt(1)),
	} {
		buf := varint.EncodeU128(value)
		got, _, err := varint.DecodeU128(buf)
		require.NoError(t, err)
		assert.Equal(t, value, got)
	}
}
