package models

import "time"

type Merchant struct {
	ID 			string		`json:"id" db:"id"`
	Name 		string		`json:"name" db:"name"`
	Type 		string		`json:"type" db:"type"`
	Address 	string		`json:"address" db:"address"`
	CreatedAt	time.Time	`json:"created_at,omitempty" db:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at,omitempty" db:"updated_at"`
}
