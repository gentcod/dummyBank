package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

//PasetoGenerator is a PASETO Token maker
type PasetoGenerator struct {
	paseto *paseto.V2
	symmetricKey []byte
}

//NewPasetoGenerator creates a new PasetoGenerator
func NewPasetoGenerator(symmetricKey string) (Generator, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoGenerator{
		paseto: paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

	//CreateToken creates a new token for a specific username and duration
	func(maker *PasetoGenerator) CreateToken(username string, duration time.Duration) (string, error) {
		payload, err := NewPayload(username, duration)
		if err != nil {
			return "", err
		}

		return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	}

	//VerifyToken checks if the token is valid or not
	func(maker *PasetoGenerator) VerifyToken(token string) (*Payload, error) {
		payload := &Payload{}

		err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
		if err != nil {
			return nil, ErrInvalidToken
		}

		err = payload.Valid()
		if err != nil {
			return nil, err
		}

		return payload, nil
	}