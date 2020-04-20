package helpers

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	alphanumerics = `ABCDEFGHJKMNOPQRSTUVWXYZ0123456789`
)

// ParseInt64 helper to avoid code repetition
func ParseInt64(stringToParse string) (int64, error) {
	intID, err := strconv.ParseInt(stringToParse, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse string to int")
	}
	return intID, nil
}

// GenerateAlphanumericToken helper to generate a length n alphanumeric string
func GenerateAlphanumericToken(length int) string {
	rand.Seed(time.Now().UnixNano())
	var grand = rand.New(rand.NewSource(time.Now().UnixNano()))
	var b bytes.Buffer
	for i := 0; i < length; i++ {
		b.WriteRune(rune(alphanumerics[grand.Intn(len(alphanumerics))]))
	}
	return strings.ToLower(b.String())
}
