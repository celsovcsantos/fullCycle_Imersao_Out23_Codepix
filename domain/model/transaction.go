package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending   = "pending"
	TransactionCompleted = "completed"
	TransactionError     = "error"
	TransactionConfirmed = "confirmed"
)

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transaction []*Transaction
}

type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"-"`
	Amount            float64  `json:"amount" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	Status            string   `json:"status" valid:"notnull"`
	Description       string   `json:"description" valid:"notnull"`
	CancelDescription string   `json:"cancel_description" valid:"-"`
}

func (t *Transaction) isValid() error {
	_, error := govalidator.ValidateStruct(t)

	if t.Amount <= 0 {
		return errors.New("the amount must be greater than 0")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionError {
		return errors.New("invalid status for the transaction")
	}

	if t.PixKeyTo.AccountID == t.AccountFrom.ID {
		return errors.New("the source and destination account cannot be the same")
	}

	if error != nil {
		return error
	}
	return nil
}

func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount:      amount,
		PixKeyTo:    pixKeyTo,
		Status:      TransactionPending,
		Description: description,
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	error := transaction.isValid()
	if error != nil {
		return nil, error
	}

	return &transaction, nil
}

func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()
	error := t.isValid()
	return error
}

func (t *Transaction) Confirm() error {
	t.Status = TransactionConfirmed
	t.UpdatedAt = time.Now()
	error := t.isValid()
	return error
}

func (t *Transaction) Cancel(description string) error {
	t.Status = TransactionError
	t.UpdatedAt = time.Now()
	t.Description = description
	error := t.isValid()
	return error
}
