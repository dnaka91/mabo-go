package buf

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"unicode/utf8"

	"github.com/dnaka91/mabo/internal/varint"
)

type InsufficientDataError struct{}

func (e *InsufficientDataError) Error() string {
	return "insufficient data"
}

type NonUtf8Error struct{}

func (e *NonUtf8Error) Error() string {
	return "string is not valid UTF-8"
}

type MissingFieldError struct {
	ID    uint32
	Field string
}

func (e *MissingFieldError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("missing field named `%v` with ID %v", e.Field, e.ID)
	}
	return fmt.Sprintf("missing field with ID %v", e.ID)
}

type UnknownVariantError struct {
	ID uint32
}

func (e *UnknownVariantError) Error() string {
	return fmt.Sprintf("unknown enum variant with ID %v", e.ID)
}

const EndMarker uint32 = 0

func DecodeBool(r []byte) ([]byte, bool, error) {
	if len(r) < 1 {
		return nil, false, &InsufficientDataError{}
	}
	return r[1:], r[0] != 0, nil
}

func DecodeU8(r []byte) ([]byte, byte, error) {
	if len(r) < 1 {
		return nil, 0, &InsufficientDataError{}
	}
	return r[1:], r[0], nil
}

func DecodeU16(r []byte) ([]byte, uint16, error) {
	value, consumed, err := varint.DecodeU16(r)
	if err != nil {
		return nil, 0, err
	}
	return r[consumed:], value, nil
}

func DecodeU32(r []byte) ([]byte, uint32, error) {
	value, consumed, err := varint.DecodeU32(r)
	if err != nil {
		return nil, 0, err
	}
	return r[consumed:], value, nil
}

func DecodeU64(r []byte) ([]byte, uint64, error) {
	value, consumed, err := varint.DecodeU64(r)
	if err != nil {
		return nil, 0, err
	}
	return r[consumed:], value, nil
}

func DecodeU128(r []byte) ([]byte, *big.Int, error) {
	value, consumed, err := varint.DecodeU128(r)
	if err != nil {
		return nil, nil, err
	}
	return r[consumed:], value, nil
}

func DecodeI8(r []byte) ([]byte, int8, error) {
	if len(r) < 1 {
		return nil, 0, &InsufficientDataError{}
	}
	return r[1:], int8(r[0]), nil
}

func DecodeI16(r []byte) ([]byte, int16, error) {
	value, consumed, err := varint.DecodeI16(r)
	if err != nil {
		return nil, 0, err
	}
	return r[consumed:], value, nil
}

func DecodeI32(r []byte) ([]byte, int32, error) {
	value, consumed, err := varint.DecodeI32(r)
	if err != nil {
		return nil, 0, err
	}
	return r[consumed:], value, nil
}

func DecodeI64(r []byte) ([]byte, int64, error) {
	value, consumed, err := varint.DecodeI64(r)
	if err != nil {
		return nil, 0, err
	}
	return r[consumed:], value, nil
}

func DecodeI128(r []byte) ([]byte, *big.Int, error) {
	value, consumed, err := varint.DecodeI128(r)
	if err != nil {
		return nil, nil, err
	}
	return r[consumed:], value, nil
}

func DecodeF32(r []byte) ([]byte, float32, error) {
	if len(r) < 4 {
		return nil, 0, &InsufficientDataError{}
	}

	return r[4:], math.Float32frombits(binary.BigEndian.Uint32(r)), nil
}

func DecodeF64(r []byte) ([]byte, float64, error) {
	if len(r) < 8 {
		return nil, 0, &InsufficientDataError{}
	}

	return r[8:], math.Float64frombits(binary.BigEndian.Uint64(r)), nil
}

func DecodeString(r []byte) ([]byte, string, error) {
	r, buf, err := DecodeBytes(r)
	if err != nil {
		return nil, "", err
	}

	if utf8.Valid(buf) {
		return r, string(buf), nil
	}

	return nil, "", &NonUtf8Error{}
}

func DecodeBytes(r []byte) ([]byte, []byte, error) {
	r, size, err := DecodeU64(r)
	if err != nil {
		return nil, nil, err
	}

	if len(r) < int(size) {
		return nil, nil, &InsufficientDataError{}
	}

	buf := make([]byte, size)
	copy(buf, r[:size])

	return r[size:], buf, nil
}

func DecodeVec[T any](r []byte, decode func([]byte) ([]byte, T, error)) ([]byte, []T, error) {
	r, size, err := DecodeU64(r)
	if err != nil {
		return nil, nil, err
	}

	vec := make([]T, size)
	for i := 0; i < int(size); i++ {
		r2, value, err := decode(r)
		if err != nil {
			return nil, nil, err
		}

		vec[i] = value
		r = r2
	}

	return r, vec, nil
}

func DecodeHashMap[K comparable, V any](
	r []byte,
	decodeKey func([]byte) ([]byte, K, error),
	decodeValue func([]byte) ([]byte, V, error),
) ([]byte, map[K]V, error) {
	r, size, err := DecodeU64(r)
	if err != nil {
		return nil, nil, err
	}

	m := make(map[K]V, size)
	for i := 0; i < int(size); i++ {
		r2, key, err := decodeKey(r)
		if err != nil {
			return nil, nil, err
		}

		r2, value, err := decodeValue(r2)
		if err != nil {
			return nil, nil, err
		}

		m[key] = value
		r = r2
	}

	return r, m, nil
}

func DecodeHashSet[T comparable](r []byte, decode func([]byte) ([]byte, T, error)) ([]byte, map[T]struct{}, error) {
	r, size, err := DecodeU64(r)
	if err != nil {
		return nil, nil, err
	}

	set := make(map[T]struct{}, size)
	for i := 0; i < int(size); i++ {
		r2, value, err := decode(r)
		if err != nil {
			return nil, nil, err
		}

		set[value] = struct{}{}
		r = r2
	}

	return r, set, nil
}

func DecodeOption[T any](r []byte, decode func([]byte) ([]byte, T, error)) ([]byte, *T, error) {
	r, some, err := DecodeU8(r)
	if err != nil {
		return nil, nil, err
	}

	if some == 1 {
		r, value, err := decode(r)
		if err != nil {
			return nil, nil, err
		}

		return r, &value, nil
	}

	return r, nil, nil
}

func DecodeID(r []byte) ([]byte, uint32, error) {
	return DecodeU32(r)
}

type Decode interface {
	Decode(r []byte) ([]byte, error)
}
