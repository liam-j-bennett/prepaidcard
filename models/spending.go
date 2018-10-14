package models

import "time"

type Spending struct {
	CardNumber		string		`json:"card_number" db:"card_id"`
	TransactionId	string		`json:"transaction_id" db:"transaction_id"`
	MerchantType	string		`json:"merchant_type" db:"merchant_type"`
	MerchantName	string		`json:"merchant_name" db:"merchant_name"`
	OriginalAmount	int64		`json:"authorized_amount" db:"auth_amount"`
	CapturedAmount	int64		`json:"amount" db:"amount"`
	Time 			time.Time	`json:"time" db:"auth_time"`
}

type SpendingList struct {
	SpendingList	[]*Spending `json:"spending"`
}