package td_test

import (
	"strings"
	"testing"
	"unicode"

	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/assert"
)

var n = 100

func TestInt64(t *testing.T) {
	for range n {
		val := td.Int64()
		assert.IsType(t, int64(0), val)
	}
}

func TestInt32(t *testing.T) {
	for range n {
		val := td.Int32()
		assert.IsType(t, int32(0), val)
	}
}

func TestInt31(t *testing.T) {
	for range n {
		val := td.Int31()
		assert.IsType(t, int32(0), val)
		assert.GreaterOrEqual(t, val, int32(0))
	}
}

func TestIntn(t *testing.T) {
	for range n {
		val := td.Intn(100)
		assert.GreaterOrEqual(t, val, 0)
		assert.Less(t, val, 100)
	}
}

func TestUint64(t *testing.T) {
	for range n {
		val := td.Uint64()
		assert.IsType(t, uint64(0), val)
	}
}

func TestUint32(t *testing.T) {
	for range n {
		val := td.Uint32()
		assert.IsType(t, uint32(0), val)
	}
}

func TestUint8(t *testing.T) {
	for range n {
		val := td.Uint8()
		assert.IsType(t, uint8(0), val)
		assert.GreaterOrEqual(t, val, uint8(1))
		assert.LessOrEqual(t, val, uint8(255))
	}
}

func TestDate(t *testing.T) {
	for range n {
		val := td.Date()
		assert.GreaterOrEqual(t, val.Year(), 1970)
		assert.LessOrEqual(t, val.Year(), 2069)
		assert.GreaterOrEqual(t, int(val.Month()), 1)
		assert.LessOrEqual(t, int(val.Month()), 12)
		assert.GreaterOrEqual(t, val.Day(), 1)
		assert.LessOrEqual(t, val.Day(), 28)
	}
}

func TestString(t *testing.T) {
	for range n {
		val := td.String()
		assert.Len(t, val, 10)
		for _, c := range val {
			assert.True(t, unicode.IsLetter(c))
		}
	}
}

func TestEmail(t *testing.T) {
	for range n {
		val := td.Email()
		assert.Contains(t, val, "@")
		parts := strings.Split(val, "@")
		assert.Len(t, parts, 2)
		assert.NotEmpty(t, parts[0])
		assert.Contains(t, parts[1], ".")
		domainParts := strings.Split(parts[1], ".")
		assert.Len(t, domainParts, 2)
		assert.NotEmpty(t, domainParts[0])
		assert.NotEmpty(t, domainParts[1])
	}
}

func TestStringn(t *testing.T) {
	for range n {
		number := 5 + td.Intn(15)
		val := td.Stringn(number)
		assert.Len(t, val, number)
		for _, c := range val {
			assert.True(t, unicode.IsLetter(c))
		}
	}
}

func TestFloat64(t *testing.T) {
	for range n {
		val := td.Float64()
		assert.IsType(t, float64(0), val)
		assert.GreaterOrEqual(t, val, 0.0)
		assert.Less(t, val, 1.0)
	}
}

func TestBool(t *testing.T) {
	m := make(map[bool]int)
	for range n {
		val := td.Bool()
		m[val]++
		assert.IsType(t, false, val)
		assert.Contains(t, []bool{true, false}, val)
	}
	assert.NotZero(t, m[true])
	assert.NotZero(t, m[false])
}

func TestBytesn(t *testing.T) {
	for range n {
		size := 5 + td.Intn(15)
		val := td.Bytesn(size)
		assert.Len(t, val, size)
	}
}

func TestBytes(t *testing.T) {
	for range n {
		val := td.Bytes()
		assert.GreaterOrEqual(t, len(val), 5)
		assert.LessOrEqual(t, len(val), 20)
	}
}
