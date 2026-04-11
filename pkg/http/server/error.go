package server

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg a common package error.
	ErrPkg = errors.New("server")
	// ErrServerUnableListenAndServe server error.
	ErrServerUnableListenAndServe = fmt.Errorf("%w: server unable to listen and serve", ErrPkg)
	// ErrServerCertificateDoesNotExist error.
	ErrServerCertificateDoesNotExist = fmt.Errorf("%w: certificate file does not exist", ErrPkg)
	// ErrServerKeyDoesNotExist error.
	ErrServerKeyDoesNotExist = fmt.Errorf("%w: key file does not exist", ErrPkg)
	// ErrServerUnableToMarshal error.
	ErrServerUnableToMarshal = fmt.Errorf("%w: unable to marshal", ErrPkg)
	// ErrServerShutdownFailure error.
	ErrServerShutdownFailure = fmt.Errorf("%w: unable to shutdown", ErrPkg)
)
