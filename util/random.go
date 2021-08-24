package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// run every time package is run
func init() {
	// ensure the random data is always random
	rand.Seed(time.Now().UnixNano())
}

//generates a random int between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max -min + 1) // return int between 0 and max - min
}

// Generate Random String of length n
func RandomString(n int) string {
	// create string builder object
	var sb strings.Builder
	k := len(alphabet)


	// for loop to generate random string of n random characters 
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)] //assign random char 
		sb.WriteByte(c) // write to string builder 
	}
	return sb.String()
}

// generate a random owner name of len 6
func RandomOwner() string {
	return RandomString(6)
}

// generate random int between 0 - 1000
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// generate random currency 
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
