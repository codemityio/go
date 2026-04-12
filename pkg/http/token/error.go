package token

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg describes the base package error.
	ErrPkg = errors.New(`token`)
	// ErrSignerUnableGenerateUUID unable to generate UUID.
	ErrSignerUnableGenerateUUID = fmt.Errorf(`%w: %s`, ErrPkg, `unable to generate random UUID`)
	// ErrSignerUnableToSignToken token signature error.
	ErrSignerUnableToSignToken = fmt.Errorf(`%w: %s`, ErrPkg, `unable to sign token`)
	// ErrUUIDProviderNewRandom error.
	ErrUUIDProviderNewRandom = fmt.Errorf(`%w: %s`, ErrPkg, `unable to generate new random UUID`)
)
