package customer

import (
	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidCustomer = errors.New("invalid customer")
)

type Customer struct {
	Name        string `validate:"required,min=3"`
	Lastname    string `validate:"required,min=3"`
	PhoneNumber string `validate:"startswith=0"`
	Email       string `validate:"email"`
}

func (c Customer) Validate() error {
	err := validator.New().Struct(c)
	if err != nil {
		return errors.Errorf("%s: %w", err, ErrInvalidCustomer)
	}
	return nil
}
