package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password1 := "GoodPassword"
	password2 := "BadPassword"

	hashedPass, err := HashPassword(password1)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass)

	err = CheckPassword(password1, hashedPass)
	require.NoError(t, err)

	err = CheckPassword(password2, hashedPass)
	require.Error(t, err)
}
