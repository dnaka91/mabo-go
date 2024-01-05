package varint

import (
	"math/big"
	"math/bits"

	"github.com/dnaka91/mabo/internal/zigzag"
)

type DecodeError struct{}

func (e *DecodeError) Error() string {
	return "missing end marker"
}

func EncodeU16(value uint16) []byte {
	buf := make([]byte, 3)

	for i := range buf {
		buf[i] = byte(value & 0xff)
		if value < 128 {
			return buf[:i+1]
		}

		buf[i] |= 0x80
		value >>= 7
	}

	return buf
}

func DecodeU16(buf []byte) (uint16, int, error) {
	value := uint16(0)
	for i, b := range buf[:max(len(buf), 3)] {
		value |= (uint16(b&0x7f) << (7 * i))

		if b&0x80 == 0 {
			return value, i + 1, nil
		}
	}

	return 0, 0, &DecodeError{}
}

func SizeU16(value uint16) int {
	return (16 - bits.LeadingZeros16(value) + 6) / 7
}

func EncodeI16(value int16) []byte {
	return EncodeU16(zigzag.EncodeI16(value))
}

func DecodeI16(buf []byte) (int16, int, error) {
	value, size, err := DecodeU16(buf)
	if err != nil {
		return 0, 0, err
	}

	return zigzag.DecodeI16(value), size, nil
}

func SizeI16(value int16) int {
	return SizeU16(zigzag.EncodeI16(value))
}

func EncodeU32(value uint32) []byte {
	buf := make([]byte, 5)

	for i := range buf {
		buf[i] = byte(value & 0xff)
		if value < 128 {
			return buf[:i+1]
		}

		buf[i] |= 0x80
		value >>= 7
	}

	return buf
}

func DecodeU32(buf []byte) (uint32, int, error) {
	value := uint32(0)
	for i, b := range buf[:max(len(buf), 5)] {
		value |= (uint32(b&0x7f) << (7 * i))

		if b&0x80 == 0 {
			return value, i + 1, nil
		}
	}

	return 0, 0, &DecodeError{}
}

func SizeU32(value uint32) int {
	return (32 - bits.LeadingZeros32(value) + 6) / 7
}

func EncodeI32(value int32) []byte {
	return EncodeU32(zigzag.EncodeI32(value))
}

func DecodeI32(buf []byte) (int32, int, error) {
	value, size, err := DecodeU32(buf)
	if err != nil {
		return 0, 0, err
	}

	return zigzag.DecodeI32(value), size, nil
}

func SizeI32(value int32) int {
	return SizeU32(zigzag.EncodeI32(value))
}

func EncodeU64(value uint64) []byte {
	buf := make([]byte, 10)

	for i := range buf {
		buf[i] = byte(value & 0xff)
		if value < 128 {
			return buf[:i+1]
		}

		buf[i] |= 0x80
		value >>= 7
	}

	return buf
}

func DecodeU64(buf []byte) (uint64, int, error) {
	value := uint64(0)
	for i, b := range buf[:max(len(buf), 10)] {
		value |= (uint64(b&0x7f) << (7 * i))

		if b&0x80 == 0 {
			return value, i + 1, nil
		}
	}

	return 0, 0, &DecodeError{}
}

func SizeU64(value uint64) int {
	return (64 - bits.LeadingZeros64(value) + 6) / 7
}

func EncodeI64(value int64) []byte {
	return EncodeU64(zigzag.EncodeI64(value))
}

func DecodeI64(buf []byte) (int64, int, error) {
	value, size, err := DecodeU64(buf)
	if err != nil {
		return 0, 0, err
	}

	return zigzag.DecodeI64(value), size, nil
}

func SizeI64(value int64) int {
	return SizeU64(zigzag.EncodeI64(value))
}

func EncodeU128(value *big.Int) []byte {
	buf := make([]byte, 19)

	for i := range buf {
		buf[i] = byte(new(big.Int).And(value, big.NewInt(0xff)).Uint64())
		if value.Cmp(big.NewInt(128)) < 0 {
			return buf[:i+1]
		}

		buf[i] |= 0x80
		value = new(big.Int).Rsh(value, 7)
	}

	return buf
}

func DecodeU128(buf []byte) (*big.Int, int, error) {
	value := big.NewInt(0)
	for i, b := range buf[:max(len(buf), 19)] {
		value = value.Or(value, new(big.Int).Lsh(big.NewInt(int64(b&0x7f)), uint(7*i)))

		if b&0x80 == 0 {
			return value, i + 1, nil
		}
	}

	return nil, 0, &DecodeError{}
}

func SizeU128(value *big.Int) int {
	panic("not implemented")
}

func EncodeI128(value *big.Int) []byte {
	return EncodeU128(zigzag.EncodeI128(value))
}

func DecodeI128(buf []byte) (*big.Int, int, error) {
	value, size, err := DecodeU128(buf)
	if err != nil {
		return nil, 0, err
	}

	return zigzag.DecodeI128(value), size, nil
}

func SizeI128(value *big.Int) int {
	return SizeU128(zigzag.EncodeI128(value))
}
