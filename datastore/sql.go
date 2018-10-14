package datastore

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid"
	"math/rand"
	"prepaidcard/models"
	"time"
)

var tables = [...]string{
	`CREATE TABLE IF NOT EXISTS cards (
	card_number varchar(256) NOT NULL PRIMARY KEY,
	full_balance bigint NOT NULL,
	blocked_balance bigint NOT NULL,
	created_at timestamp without time zone,
	updated_at timestamp without time zone
);`,

	`CREATE TABLE IF NOT EXISTS merchants (
	id varchar(256) NOT NULL PRIMARY KEY,
	name varchar(256) NOT NULL UNIQUE,
	type varchar(256) NOT NULL,
	address text NOT NULL,
	created_at timestamp without time zone,
	updated_at timestamp without time zone
);`,

	`CREATE TABLE IF NOT EXISTS transactions (
	id varchar(256) NOT NULL PRIMARY KEY,
	card_id varchar(256) NOT NULL,
	merchant_id varchar(256) NOT NULL,
	original_amount	bigint NOT NULL,
	authorized_amount bigint NOT NULL,
	captured_amount bigint NOT NULL,
	created_at timestamp without time zone,
	updated_at timestamp without time zone
);`,

	`CREATE OR REPLACE VIEW user_transaction_list AS
	SELECT cards.card_number card_id, transactions.id transaction_id, merchants.type merchant_type, merchants.name merchant_name, transactions.original_amount auth_amount, transactions.captured_amount amount, transactions.created_at auth_time FROM transactions
	JOIN merchants ON transactions.merchant_id = merchants.id
	JOIN cards ON transactions.card_id = cards.card_number
;`,
}

const (
	cardNumberLength = 16
	cardNumbers = "0123456789"
	cardIdSelector = `SELECT * FROM cards WHERE card_number=?`
	merchantIdSelector = `SELECT * FROM merchants WHERE id=?`
	transactionIdSelector = `SELECT * FROM transactions WHERE id=?`
	transactionListQuery = `SELECT * FROM user_transaction_list WHERE card_id=? ORDER BY user_transaction_list.auth_time DESC`
)

type SQLStore struct {
	db	*sqlx.DB
}

func newId(createdTime time.Time) ulid.ULID {
	now := ulid.Timestamp(createdTime)
	id, _ := ulid.New(now, nil) // Only err if createdTime > max time in unix ms
	return id
}

func InitDB(db *sqlx.DB) (*SQLStore, error) {
	ds := &SQLStore{db: db}
	tx, err := ds.db.Beginx()
	if err != nil {
		return nil, err
	}
	for _, v := range tables {
		_, err := tx.Exec(v)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return ds, nil
}

func (s *SQLStore) CreateCard() (*models.PrepaidCard, error) {
	var card models.PrepaidCard
	card.CreatedAt = time.Now()
	card.UpdatedAt = card.CreatedAt
	// TODO: Feels kind of hacky, but gets the job done for now
	newNumberBytes := make([]byte, cardNumberLength)
	rand.Seed(time.Now().Unix())
	for i := range newNumberBytes {
		newNumberBytes[i] = cardNumbers[rand.Intn(len(cardNumbers))]
	}
	card.CardNumber = string(newNumberBytes)
	query := s.db.Rebind(`INSERT INTO cards (
			card_number,
			full_balance,
			blocked_balance,
			created_at,
			updated_at
	)
	VALUES (
			:card_number,
			:full_balance,
			:blocked_balance,
			:created_at,
			:updated_at
	);`)
	_, err := s.db.NamedExec(query, &card)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (s *SQLStore) GetCard(cardId string) (*models.PrepaidCard, error) {
	var card models.PrepaidCard
	query := s.db.Rebind(cardIdSelector)
	row := s.db.QueryRowx(query, cardId)
	err := row.StructScan(&card)
	if err == sql.ErrNoRows {
		return nil, models.NotFound
	}
	if err != nil {
		return nil, err
	}
	return &card, err
}

func (s *SQLStore) LoadCard(cardId string, amount int64) (*models.PrepaidCard, error) {
	var card models.PrepaidCard
	query := s.db.Rebind(cardIdSelector)
	row := s.db.QueryRowx(query, cardId)
	err := row.StructScan(&card)
	if err == sql.ErrNoRows {
		return nil, models.NotFound
	}
	if err != nil {
		return nil, err
	}
	card.FullBalance = card.FullBalance + amount
	query = s.db.Rebind(`UPDATE cards SET full_balance=:full_balance WHERE card_number=:card_number`)
	_, err = s.db.NamedExec(query, card)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (s *SQLStore) TransactionList(cardId string) (*models.SpendingList, error) {
	var listModel models.SpendingList
	list, err := s.transactionList(cardId)
	listModel.SpendingList = list
	if err == sql.ErrNoRows {
		return &listModel, nil
	}
	if err != nil {
		return nil, err
	}
	return &listModel, nil
}

func (s *SQLStore) transactionList(cardId string) ([]*models.Spending, error) {
	var list []*models.Spending
	query := s.db.Rebind(transactionListQuery)
	rows, err := s.db.Queryx(query, cardId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var spending models.Spending
		err := rows.StructScan(&spending)
		if err != nil {
			if err == sql.ErrNoRows {
				return list, nil
			}
			return list, err
		}
		list = append(list, &spending)
	}
	if err := rows.Err(); err != nil {
		return list, err
	}
	return list, nil
}

func (s *SQLStore) CreateMerchant(newMerchant *models.Merchant) (*models.Merchant, error) {
	merchant := new(models.Merchant)
	*merchant = *newMerchant
	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = merchant.CreatedAt
	query := s.db.Rebind(`INSERT INTO merchants (
			id,
			name,
			type,
			address,
			created_at,
			updated_at
	)
	VALUES (
			:id,
			:name,
			:type,
			:address,
			:created_at,
			:updated_at
	);`)
	_, err := s.db.NamedExec(query, merchant)
	if err != nil {
		return nil, err
	}
	return merchant, nil
}

func (s *SQLStore) GetMerchant (merchantId string) (*models.Merchant, error) {
	var merchant models.Merchant
	query := s.db.Rebind(merchantIdSelector)
	row := s.db.QueryRowx(query, merchantId)
	err := row.StructScan(&merchant)
	if err == sql.ErrNoRows {
		return nil, models.NotFound
	}
	if err != nil {
		return nil, err
	}
	return &merchant, err
}

func (s *SQLStore) GetTransaction(transactionId string) (*models.Transaction, error) {
	var transaction models.Transaction
	query := s.db.Rebind(transactionIdSelector)
	row := s.db.QueryRowx(query, transactionId)
	err := row.StructScan(&transaction)
	if err == sql.ErrNoRows {
		return nil, models.NotFound
	}
	if err != nil {
		return nil, err
	}
	return &transaction, err
}

/*
	Performs a card Auth
	- Check that amount < full_balance - blocked_balance
	- Create transaction with amount for Card & Merchant
	- Add amount to the blocked_balance
 */
func (s *SQLStore) Auth(card *models.PrepaidCard, merchant *models.Merchant, amount int64) (*models.Transaction, error) {
	var transaction models.Transaction
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = transaction.CreatedAt
	transaction.ID = newId(transaction.CreatedAt).String()
	transaction.OriginalAmount = amount
	transaction.AuthorizedAmount = amount
	tx, err:= s.db.Beginx()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if amount > (card.FullBalance - card.BlockedBalance) { // Should be checked in the API, but let's make it defensive
		tx.Rollback()
		return nil, models.InvalidCardBalance
	}
	card.BlockedBalance = card.BlockedBalance + amount
	transaction.CardID = card.CardNumber
	transaction.Card = card
	transaction.MerchantID = merchant.ID
	transaction.Merchant = merchant
	query := tx.Rebind(`INSERT INTO transactions (
			id,
			card_id,
			merchant_id,
			original_amount,
			authorized_amount,
			captured_amount,
			created_at,
			updated_at
	)
	VALUES (
			:id,
			:card_id,
			:merchant_id,
			:original_amount,
			:authorized_amount,
			:captured_amount,
			:created_at,
			:updated_at
	);`)
	_, err = tx.NamedExec(query, transaction)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	query = tx.Rebind(`UPDATE cards SET blocked_balance=:blocked_balance WHERE card_number=:card_number`)
	_, err = tx.NamedExec(query, card)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &transaction, nil
}

/*
	Performs a transaction capture
	- Check amount >= authorized_amount
	- Remove amount from authed and append to captured
	- Remove amount from card Full + Blocked balances
 */
func (s *SQLStore) Capture(transaction *models.Transaction, amount int64) error {
	tx, err:= s.db.Beginx()
	if err != nil {
		tx.Rollback()
		return err
	}
	if amount > transaction.AuthorizedAmount {
		tx.Rollback()
		return models.InvalidTransactionAuth
	}
	// Need to update the card amount within the transaction
	var card models.PrepaidCard
	query := tx.Rebind(cardIdSelector)
	row := tx.QueryRowx(query, transaction.CardID)
	err = row.StructScan(&card)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return models.NotFound
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	transaction.AuthorizedAmount = transaction.AuthorizedAmount - amount
	transaction.CapturedAmount = transaction.CapturedAmount + amount
	query = tx.Rebind(`UPDATE transactions SET authorized_amount=:authorized_amount, captured_amount=:captured_amount WHERE id=:id`)
	_, err = tx.NamedExec(query, transaction)
	if err != nil {
		tx.Rollback()
		return err
	}
	card.FullBalance = card.FullBalance - amount
	card.BlockedBalance = card.BlockedBalance - amount
	query = tx.Rebind(`UPDATE cards SET blocked_balance=:blocked_balance, full_balance=:full_balance WHERE card_number=:card_number`)
	_, err = tx.NamedExec(query, card)
	if err != nil {
		tx.Rollback()
		return err
	}
	transaction.Card = &card
	tx.Commit()
	return nil
}

/*
	Performs a reverse on an auth
	- Check amount <= authorized_amount
	- Remove amount from authed
	- Remove amount from Blocked balance
 */
func (s *SQLStore) Reverse(transaction *models.Transaction, amount int64) error {
	tx, err:= s.db.Beginx()
	if err != nil {
		tx.Rollback()
		return err
	}
	if amount > transaction.AuthorizedAmount {
		tx.Rollback()
		return models.InvalidTransactionAuth
	}
	transaction.AuthorizedAmount = transaction.AuthorizedAmount - amount
	// Need to update the card amount within the transaction
	var card models.PrepaidCard
	query := tx.Rebind(cardIdSelector)
	row := tx.QueryRowx(query, transaction.CardID)
	err = row.StructScan(&card)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return models.NotFound
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	card.BlockedBalance = card.BlockedBalance - amount
	query = tx.Rebind(`UPDATE transactions SET authorized_amount=:authorized_amount WHERE id=:id`)
	_, err = tx.NamedExec(query, transaction)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = tx.Rebind(`UPDATE cards SET blocked_balance=:blocked_balance WHERE card_number=:card_number`)
	_, err = tx.NamedExec(query, card)
	if err != nil {
		tx.Rollback()
		return err
	}
	transaction.Card = &card
	tx.Commit()
	return nil
}

/*
	Performs a refund on captured funds
	- Check amount <= captured_amount
	- Add to card full_balance
	- remove captured amount
 */
func (s *SQLStore) Refund(transaction *models.Transaction, amount int64) error {
	tx, err:= s.db.Beginx()
	if err != nil {
		tx.Rollback()
		return err
	}
	if amount > transaction.CapturedAmount {
		tx.Rollback()
		return models.InvalidTransactionCaptured
	}
	transaction.CapturedAmount = transaction.CapturedAmount - amount
	// Need to update the card amount within the transaction
	var card models.PrepaidCard
	query := tx.Rebind(cardIdSelector)
	row := tx.QueryRowx(query, transaction.CardID)
	err = row.StructScan(&card)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return models.NotFound
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	card.FullBalance = card.FullBalance + amount
	query = tx.Rebind(`UPDATE cards SET full_balance=:full_balance WHERE card_number=:card_number`)
	_, err = tx.NamedExec(query, card)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = tx.Rebind(`UPDATE transactions SET captured_amount=:captured_amount WHERE id=:id`)
	_, err = tx.NamedExec(query, transaction)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
