// test data
package td

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Int64 returns random int64 value
func Int64() int64 {
	sign := rand.Intn(2)
	if sign == 0 {
		sign = -1
	}
	return rand.Int63() * int64(sign)
}

// Int32 returns random int32 value
func Int32() int32 {
	sign := rand.Intn(2)
	if sign == 0 {
		sign = -1
	}
	return rand.Int31() * int32(sign)
}

// Uint64 returns random uint64 value
func Uint64() uint64 {
	return rand.Uint64()
}

// Uint32 returns random uint32 value
func Uint32() uint32 {
	return rand.Uint32()
}

// Uint8 returns random uint8 value
func Uint8() uint8 {
	return uint8(rand.Intn(math.MaxUint8) + 1)
}

// DateString returns random date
func Date() time.Time {
	return time.Date(1970+rand.Intn(100), time.Month(1+rand.Intn(12)), 1+rand.Intn(28), rand.Intn(24), rand.Intn(60), rand.Intn(60), 0, time.Local)
}

// String returns random string value
func String() string {
	return Stringn(10)
}

// Stringn returns random string value with specified length
func Stringn(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
