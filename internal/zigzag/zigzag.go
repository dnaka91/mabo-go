package zigzag

import "math/big"

func EncodeI16(value int16) uint16 {
	return uint16((value << 1) ^ (value >> 15))
}

func DecodeI16(value uint16) int16 {
	return int16(value>>1) ^ (-(int16(value & 0b1)))
}

func EncodeI32(value int32) uint32 {
	return uint32((value << 1) ^ (value >> 31))
}

func DecodeI32(value uint32) int32 {
	return int32(value>>1) ^ (-(int32(value & 0b1)))
}

func EncodeI64(value int64) uint64 {
	return uint64((value << 1) ^ (value >> 63))
}

func DecodeI64(value uint64) int64 {
	return int64(value>>1) ^ (-(int64(value & 0b1)))
}

func EncodeI128(value *big.Int) *big.Int {
	return new(big.Int).Xor(
		new(big.Int).Lsh(value, 1),
		new(big.Int).Rsh(value, 127),
	)
}

func DecodeI128(value *big.Int) *big.Int {
	return new(big.Int).Xor(
		new(big.Int).Rsh(value, 1),
		new(big.Int).Neg(new(big.Int).And(value, big.NewInt(0b1))),
	)
}
