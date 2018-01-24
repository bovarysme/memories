package crypto

func hashCode(str string) int {
	var hash int32
	for i := 0; i < len(str); i++ {
		hash = 31*hash + int32(str[i])
	}

	return int(hash)
}

func generateKey(iv int) []byte {
	const length = 16
	var key [length]byte
	var pad [length]int8

	pad[0] = int8(iv)
	pad[1] = pad[0] - 71
	pad[2] = pad[1] - 71
	for i := 3; i < length; i++ {
		pad[i] = int8(int(pad[i-3]) ^ int(pad[i-2]) ^ 0xb9 ^ i)
	}

	factor := iv
	if iv > -2 && iv < 2 {
		factor = -313187 + 13819823*iv
	}

	term := -7
	for i := 1; i < length+1; i++ {
		index := i & (length - 1)
		value := int(pad[index])*factor + term

		term = int(int8(value >> 32))
		value = int(int32(value + term))
		if value < term {
			term++
			value++
		}

		value = -value - 2
		key[index] = byte(value)
	}

	return key[:]
}
