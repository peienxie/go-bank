package api

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// randomInt generates a random number in [min, max]
func randomInt(min, max int64) int64 {
	return rand.Int63n(max - min + 1)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randomString generates a random string of given length n
func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// randomUsername returns a random username only contains letters
func randomUsername() string {
	return randomString(10)
}

// randomMoney generates a random amount of money
func randomMoney() int64 {
	return randomInt(0, 1000)
}

// randomCurrency generates a random currency type
func randomCurrency() string {
	currencies := []string{"USD", "TWD"}
	return currencies[rand.Intn(len(currencies))]
}
