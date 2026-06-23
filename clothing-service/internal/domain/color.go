package domain

import "time"

type Color struct {
	ID        int64
	Name      string
	HexCode   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
