package models

import (
	"time"
)

type PrepaidCard struct {
	CardNumber 		string		`json:"card_number" db:"card_number"`
	FullBalance 	int64		`json:"full_balance" db:"full_balance"`
	BlockedBalance 	int64		`json:"blocked_balance" db:"blocked_balance"`
	CreatedAt		time.Time	`json:"created_at,omitempty" db:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at,omitempty" db:"updated_at"`
}
