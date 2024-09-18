package dto

import "time"

type Newsletter struct {
	ClientPublicID int64
	PublicID       int64
	TotalCount     int64
	Name           string
	Description    *string
	CreatedAt      time.Time
}
