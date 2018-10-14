package models

import "time"

type Transaction struct {
	ID 					string			`json:"id" db:"id"`
	CardID 				string			`json:"-" db:"card_id"`
	Card 				*PrepaidCard	`json:"card" db:"-"`
	MerchantID 			string			`json:"-" db:"merchant_id"`
	Merchant 			*Merchant		`json:"merchant" db:"-"`
	OriginalAmount		int64			`json:"original_amount" db:"original_amount"`
	AuthorizedAmount 	int64			`json:"authorized_amount" db:"authorized_amount"`
	CapturedAmount 		int64			`json:"captured_amount" db:"captured_amount"`
	CreatedAt			time.Time		`json:"created_at,omitempty" db:"created_at"`
	UpdatedAt			time.Time		`json:"updated_at,omitempty" db:"updated_at"`
}
