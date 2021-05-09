package errors

import "errors"

var (
	UnknownReceiverType = errors.New("unknown receiver type")
	NoSuchRecipient     = errors.New("there's no recipient with such name")
)
