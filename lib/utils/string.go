package utils

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

var lettersLen = len(letters)

func RandString(n int) string {
	nR := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[nR.Intn(lettersLen)]
	}
	return string(b)
}
