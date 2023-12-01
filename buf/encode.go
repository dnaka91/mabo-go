package buf

import (
	"encoding/binary"
	"math"
	"math/big"

	"github.com/dnaka91/stef/internal/varint"
)

func EncodeBool(w []byte, value bool) []byte {
	if value {
		return append(w, 1)
	} else {
		return append(w, 0)
	}
}

func EncodeU8(w []byte, value uint8) []byte {
	return append(w, value)
}

func EncodeU16(w []byte, value uint16) []byte {
	return append(w, varint.EncodeU16(value)...)
}

func EncodeU32(w []byte, value uint32) []byte {
	return append(w, varint.EncodeU32(value)...)
}

func EncodeU64(w []byte, value uint64) []byte {
	return append(w, varint.EncodeU64(value)...)
}

func EncodeU128(w []byte, value *big.Int) []byte {
	return append(w, varint.EncodeU128(value)...)
}

func EncodeI8(w []byte, value int8) []byte {
	return append(w, uint8(value))
}

func EncodeI16(w []byte, value int16) []byte {
	return append(w, varint.EncodeI16(value)...)
}

func EncodeI32(w []byte, value int32) []byte {
	return append(w, varint.EncodeI32(value)...)
}

func EncodeI64(w []byte, value int64) []byte {
	return append(w, varint.EncodeI64(value)...)
}

func EncodeI128(w []byte, value *big.Int) []byte {
	return append(w, varint.EncodeI128(value)...)
}

func EncodeF32(w []byte, value float32) []byte {
	return binary.BigEndian.AppendUint32(w, math.Float32bits(value))
}

func EncodeF64(w []byte, value float64) []byte {
	return binary.BigEndian.AppendUint64(w, math.Float64bits(value))
}

func EncodeString(w []byte, value string) []byte {
	return EncodeBytes(w, []byte(value))
}

func EncodeBytes(w []byte, value []byte) []byte {
	w = EncodeU64(w, uint64(len(value)))
	return append(w, value...)
}

func EncodeVec[T any](w []byte, vec []T, encode func([]byte, T) []byte) []byte {
	w = EncodeU64(w, uint64(len(vec)))

	for _, value := range vec {
		w = encode(w, value)
	}

	return w
}

func EncodeHashMap[K comparable, V any](
	w []byte,
	m map[K]V,
	encodeKey func([]byte, K) []byte,
	encodeValue func([]byte, V) []byte,
) []byte {
	w = EncodeU64(w, uint64(len(m)))

	for key, value := range m {
		w = encodeKey(w, key)
		w = encodeValue(w, value)
	}

	return w
}

func EncodeHashSet[T comparable](
	w []byte,
	m map[T]struct{},
	encode func([]byte, T) []byte,
) []byte {
	w = EncodeU64(w, uint64(len(m)))

	for value := range m {
		w = encode(w, value)
	}

	return w
}

func EncodeOption[T any](
	w []byte,
	value *T,
	encode func([]byte, T) []byte,
) []byte {
	if value != nil {
		w = EncodeU8(w, 1)
		return encode(w, *value)
	} else {
		return EncodeU8(w, 0)
	}
}

func EncodeID(w []byte, id uint32) []byte {
	return EncodeU32(w, id)
}

func EncodeField(w []byte, id uint32, encode func([]byte) []byte) []byte {
	w = EncodeID(w, id)
	return encode(w)
}

func EncodeFieldOption[T any](w []byte, id uint32, option *T, encode func([]byte, T) []byte) []byte {
	if option != nil {
		w = EncodeID(w, id)
		return encode(w, *option)
	}

	return w
}

type Encode interface {
	Encode(w []byte) []byte
}
