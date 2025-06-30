package token

import "time"

type Maker interface {
	GenerateToken(username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
