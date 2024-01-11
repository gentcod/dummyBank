package token

import "time"

//Generator is an interface for managing tokens
type Generator interface {
	//CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	//VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}