package token

import (
	"github.com/auronvila/simple-bank/util"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJwtMaker(t *testing.T) {
	role := util.DepositorRole
	maker, err := NewJwtMaker(util.RandomString(32), role)
	require.NoError(t, err)

	username := util.RandomOwner()

	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.GenerateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJwtToken(t *testing.T) {
	role := util.DepositorRole
	maker, err := NewJwtMaker(util.RandomString(32), role)
	require.NoError(t, err)

	token, payload, err := maker.GenerateToken(util.RandomOwner(), role, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestInvalidJwtToken(t *testing.T) {
	role := util.DepositorRole
	payload, err := NewPayload(util.RandomOwner(), role, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJwtMaker(util.RandomString(32), role)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}
