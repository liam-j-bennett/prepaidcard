package models


type CardStore interface {
	CreateCard() (*PrepaidCard, error)
	GetCard(cardId string) (*PrepaidCard, error)
	LoadCard(cardId string, amount int64) (*PrepaidCard, error)
	TransactionList(cardId string) (*SpendingList, error)
	CreateMerchant(newMerchant *Merchant) (*Merchant, error)
	GetMerchant(merchantId string) (*Merchant, error)
	GetTransaction(transactionId string) (*Transaction, error)
	Auth(card *PrepaidCard, merchant *Merchant, amount int64) (*Transaction, error)
	Capture(transaction *Transaction, amount int64) error
	Reverse(transaction *Transaction, amount int64) error
	Refund(transaction *Transaction, amount int64) error
}
