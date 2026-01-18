package pkg

import (
	"math/rand"
	"time"
)

var (
	letterBytes       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes       = "0123456789"
	letterNumberBytes = letterBytes + numberBytes
)

// GenerateRandomString generates a random string of a specified length using the provided charset template.
func GenerateRandomString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = letterNumberBytes[r.Intn(len(letterNumberBytes))]
	}
	return string(b)
}
