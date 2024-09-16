package dto

import "time"

type Newsletter struct {
	ClientID       int64
	ClientPublicID int64
	ID             int64
	PublicID       int64
	TotalCount     int64
	Name           string
	Description    *string
	CreatedAt      time.Time
}
