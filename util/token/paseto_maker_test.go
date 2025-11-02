package token

import (
	"testing"
	"time"

	"github.com/jasonwebb3152/webb-ai/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	username := util.RandomString(10)
	symmetricKey := util.RandomString(32)
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Minute * 30)
	token, payload, err := maker.CreateToken(username, time.Minute*30)

	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)
	require.NoError(t, err)

	require.NotZero(t, payload.Id)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.Id)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	username := util.RandomString(10)
	symmetricKey := util.RandomString(32)
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	token, payload, err := maker.CreateToken(username, time.Minute*-1)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
