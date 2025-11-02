package token

import (
	"testing"
	"time"

	"github.com/jasonwebb3152/webb-ai/util"
	"github.com/stretchr/testify/require"
)

func TestPayload(t *testing.T) {
	username := util.RandomString(8)
	now := time.Now()
	nowPlus15 := now.Add(time.Minute * 15)
	duration := time.Minute * 15

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload.Id)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, now, payload.IssuedAt, time.Second*5)
	require.WithinDuration(t, nowPlus15, payload.ExpiresAt, time.Second*5)

	err = payload.Valid()
	require.NoError(t, err)

	payload.ExpiresAt = time.Now().Add(time.Minute * -10)
	err = payload.Valid()
	require.Error(t, err)
}
