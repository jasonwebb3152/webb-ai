package db

import (
	"context"
	"testing"
	"time"

	"github.com/jasonwebb3152/webb-ai/util"
	"github.com/stretchr/testify/require"
)

func randomUser() CreateUserParams {
	return CreateUserParams{
		Username:       util.RandomString(8),
		FullName:       util.RandomString(8),
		Email:          util.RandomEmail(),
		HashedPassword: util.RandomString(20),
	}
}

func CreateRandomUser(t *testing.T) User {
	userParams := randomUser()
	hashedPassword, err := util.HashPassword(userParams.HashedPassword)
	require.NoError(t, err)
	userParams.HashedPassword = hashedPassword

	ctx := context.Background()
	user, err := testStore.CreateUser(ctx, userParams)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, userParams.Username, user.Username)
	require.Equal(t, userParams.FullName, user.FullName)
	require.Equal(t, userParams.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero()) // Defaults to zero time
	return user
}

func TestGetUser(t *testing.T) {
	user := CreateRandomUser(t)

	ctx := context.Background()
	res, err := testStore.GetUser(ctx, user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, user.Username, res.Username)
	require.Equal(t, user.Email, res.Email)
	require.Equal(t, user.FullName, res.FullName)
	require.Equal(t, user.HashedPassword, res.HashedPassword)

	require.WithinDuration(t, user.CreatedAt, res.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, res.PasswordChangedAt, time.Second)
}
