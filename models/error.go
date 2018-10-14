package models

import "errors"

var (
	NotFound = ApiError{
		code: 404,
		error: errors.New("object not found"),
	}
	InvalidAmount = ApiError{
		code: 400,
		error: errors.New("invalid amount"),
	}
	InvalidCardBalance = ApiError{
		code: 409,
		error: errors.New("invalid balance on card"),
	}
	InvalidTransactionAuth = ApiError{
		code: 409,
		error: errors.New("invalid authorized amount on transaction"),
	}
	InvalidTransactionCaptured = ApiError{
		code: 409,
		error: errors.New("invalid captured amount on transaction"),
	}
)

type Error interface {
	Code() int
	error
}

type ApiError struct {
	code int
	error
}

func (e ApiError) Code() int { return e.code }
