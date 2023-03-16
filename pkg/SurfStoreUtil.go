package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"strconv"
)

type Pair struct {
	Key   int
	Value string
}

func GetChordPosition(N int, hash string) (int, error) {
	digits := int(math.Log2(float64(N))/4) + 1
	hashTail := hash[len(hash)-digits:]
	i, err := strconv.ParseInt(hashTail, 16, 64)
	if err != nil {
		return -1, err
	}
	return int(i) % N, nil
}

func GetHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
