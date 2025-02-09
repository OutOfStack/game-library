// Package td contains functions for data testing
package td

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint: gosec
}

// Int64 returns random int64 value
func Int64() int64 {
	sign := random.Intn(2)
	if sign == 0 {
		sign = -1
	}
	return random.Int63() * int64(sign)
}

// Int32 returns random int32 value
func Int32() int32 {
	sign := random.Intn(2)
	if sign == 0 {
		sign = -1
	}
	return random.Int31() * int32(sign)
}

// Int31 returns random non-negative int32 value
func Int31() int32 {
	return random.Int31()
}

// Intn returns random value in range [0, n)
func Intn(n int) int {
	return random.Intn(n)
}

// Uint64 returns random uint64 value
func Uint64() uint64 {
	return random.Uint64()
}

// Uint32 returns random uint32 value
func Uint32() uint32 {
	return random.Uint32()
}

// Uint8 returns random uint8 value
func Uint8() uint8 {
	return uint8(random.Intn(math.MaxUint8) + 1)
}

// Date returns random date
func Date() time.Time {
	return time.Date(1970+random.Intn(100), time.Month(1+random.Intn(12)), 1+random.Intn(28), random.Intn(24),
		random.Intn(60), random.Intn(60), 0, time.UTC)
}

// String returns random string value
func String() string {
	return Stringn(10)
}

// Email returns random email address
func Email() string {
	return fmt.Sprintf("%s@%s.%s", String(), String(), Stringn(3))
}

// Stringn returns random string value with specified length
func Stringn(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for range n {
		c := alphabet[random.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Float64 returns random float64 value
func Float64() float64 {
	return random.Float64()
}

// Bool returns random bool value
func Bool() bool {
	return random.Intn(2) == 0
}
