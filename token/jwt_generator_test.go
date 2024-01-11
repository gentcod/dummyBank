package token

import (
	"testing"
	"time"

	"github.com/gentcod/DummyBank/util"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTGenerator(t *testing.T) {
	username := util.RandomOwner()
	duration := time.Minute

	maker, err := NewJWTGenerator(util.RandomStr(32))
	require.NoError(t, err)

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	duration := time.Minute

	maker, err := NewJWTGenerator(util.RandomStr(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err,ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTToken(t *testing.T) {
	//When secretKey length is invalid
	secret1 := util.RandomStr(20)

	maker1, err := NewJWTGenerator(secret1)
	require.Error(t, err)
	require.Nil(t, maker1)

	//When none signature token type is used
	payloadAlgNone, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payloadAlgNone)
	tokenAlg, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	makerAlg, err := NewJWTGenerator(util.RandomStr(32))
	require.NoError(t, err)

	payloadAlgNone, err = makerAlg.VerifyToken(tokenAlg)
	require.Error(t, err)
	require.EqualError(t, err,ErrInvalidToken.Error())
	require.Nil(t, payloadAlgNone)
}