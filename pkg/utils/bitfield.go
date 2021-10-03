package utils

type bitField struct{}

func (*bitField) AddBits(sum int64, add int64) int64 {
	sum |= add
	return sum
}

func (*bitField) RemoveBits(sum int64, remove int64) int64 {
	sum &= ^remove
	return sum
}

func (*bitField) HasBits(sum int64, bit int64) bool {
	return (sum & bit) == bit
}

var BitField = bitField{}
