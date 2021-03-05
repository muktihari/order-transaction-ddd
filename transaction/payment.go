package transaction

import (
	"encoding/base64"
	"errors"
)

var (
	// ErrPaymentTypeNotAllowed tells that payment type to purchase the order is not allowed.
	ErrPaymentTypeNotAllowed = errors.New("error payment method not allowed")
	// ErrPaymentProofIsNotBase64EncodedString tells that payment proof is not a base64 encoded string.
	ErrPaymentProofIsNotBase64EncodedString = errors.New("payment proof is not base64 encoded string")
)

// PaymentSpecification contains information about a payment: its type,
// name holder, and identifier ID that can be verified.
type PaymentSpecification struct {
	Type         PaymentType  `bson:"type" json:"type"`
	NameHolder   string       `bson:"name_holder" json:"name_holder"`
	IdentifierID string       `bson:"identifier_id" json:"identifier_id"`
	Proof        PaymentProof `bson:"proof" json:"proof"`
}

// PaymentType type of payment
type PaymentType int

const (
	// PaymentTypeBankTransfer tells that payment method is a bank transfer.
	PaymentTypeBankTransfer PaymentType = iota + 1
)

// PaymentProof is a base64 string image
type PaymentProof string

// Validate verified that payment specification is valid.
func (p PaymentSpecification) Validate() error {
	if p.Type != PaymentTypeBankTransfer {
		return ErrPaymentTypeNotAllowed
	}
	_, err := base64.RawStdEncoding.DecodeString(string(p.Proof))
	if err != nil {
		return ErrPaymentProofIsNotBase64EncodedString
	}

	return nil
}
