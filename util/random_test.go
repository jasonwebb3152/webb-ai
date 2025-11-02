package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandom(t *testing.T) {
	number := RandomInt(0, 10)
	require.LessOrEqual(t, number, int64(10))
	require.GreaterOrEqual(t, number, int64(0))

	randomString := RandomString(25)
	require.Equal(t, 25, len(randomString))

	randomPassword := RandomPassword(12)
	require.Equal(t, 15, len(randomPassword))
	require.Equal(t, "A1!", randomPassword[12:])

	randomEmail := RandomEmail()
	require.Equal(t, 20, len(randomEmail))
}
