package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) int64 {
	/** Generates a random integer from min to max. */
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	/** Generates a random string of length n. */
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomPassword(n int) string {
	s1 := RandomString(n)
	s2 := "A1!"

	var builder strings.Builder
	builder.WriteString(s1)
	builder.WriteString(s2)
	return builder.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(10))
}
