package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomNumber(num int) int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(num)))
	if err != nil {
		fmt.Println("error:", err)
		return 0
	}
	return n.Int64()
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randomNumber(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
