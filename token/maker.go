package token

import "time"

type Maker interface {
	GenerateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
