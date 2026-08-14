package token

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg describes the base package error.
	ErrPkg = errors.New(`token`)
	// ErrSignerUnableToGenerateUUID unable to generate UUID.
	ErrSignerUnableToGenerateUUID = fmt.Errorf(`%w: %s`, ErrPkg, `unable to generate random UUID`)
	// ErrSignerUnableToSignToken token signature error.
	ErrSignerUnableToSignToken = fmt.Errorf(`%w: %s`, ErrPkg, `unable to sign token`)
	// ErrUUIDProviderUnableToGenerateUUID error.
	ErrUUIDProviderUnableToGenerateUUID = fmt.Errorf(
		`%w: %s`,
		ErrPkg,
		`unable to generate new random UUID`,
	)
)
