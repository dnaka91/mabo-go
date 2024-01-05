package buf

import (
	"math/big"

	"github.com/dnaka91/mabo/internal/varint"
)

func SizeBool(value bool) int {
	return 1
}

func SizeU8(value uint8) int {
	return 1
}

func SizeU16(value uint16) int {
	return varint.SizeU16(value)
}

func SizeU32(value uint32) int {
	return varint.SizeU32(value)
}

func SizeU64(value uint64) int {
	return varint.SizeU64(value)
}

func SizeU128(value *big.Int) int {
	return varint.SizeU128(value)
}

func SizeI8(value int8) int {
	return 1
}

func SizeI16(value int16) int {
	return varint.SizeI16(value)
}

func SizeI32(value int32) int {
	return varint.SizeI32(value)
}

func SizeI64(value int64) int {
	return varint.SizeI64(value)
}

func SizeI128(value *big.Int) int {
	return varint.SizeI128(value)
}

func SizeF32(value float32) int {
	return 4
}

func SizeF64(value float64) int {
	return 8
}

func SizeString(value string) int {
	return SizeU64(uint64(len(value))) + len(value)
}

func SizeBytes(value []byte) int {
	return SizeU64(uint64(len(value))) + len(value)
}

func SizeVec[T any](vec []T, size func(T) int) int {
	s := SizeU64(uint64(len(vec)))

	for _, value := range vec {
		s += size(value)
	}

	return s
}

func SizeHashMap[K comparable, V any](
	m map[K]V,
	sizeKey func(K) int,
	sizeValue func(V) int,
) int {
	s := SizeU64(uint64(len(m)))

	for key, value := range m {
		s += sizeKey(key)
		s += sizeValue(value)
	}

	return s
}

func SizeHashSet[T comparable](
	m map[T]struct{},
	size func(T) int,
) int {
	s := SizeU64(uint64(len(m)))

	for value := range m {
		s += size(value)
	}

	return s
}

func SizeOption[T any](
	value *T,
	size func(T) int,
) int {
	if value != nil {
		s := SizeU8(1)
		return s + size(*value)
	} else {
		return SizeU8(0)
	}
}

func SizeID(id uint32) int {
	return SizeU32(id)
}

func SizeField(id uint32, size func() int) int {
	s := SizeID(id)
	return s + size()
}

func SizeFieldOption[T any](id uint32, option *T, size func(T) int) int {
	if option != nil {
		s := SizeID(id)
		return s + size(*option)
	}

	return 0
}

type Size interface {
	Size(w []byte) []byte
}
