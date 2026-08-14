package server

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg a common package error.
	ErrPkg = errors.New("server")
	// ErrServerUnableToListenAndServe server error.
	ErrServerUnableToListenAndServe = fmt.Errorf("%w: server unable to listen and serve", ErrPkg)
	// ErrServerCertificateDoesNotExist error.
	ErrServerCertificateDoesNotExist = fmt.Errorf("%w: certificate file does not exist", ErrPkg)
	// ErrServerKeyDoesNotExist error.
	ErrServerKeyDoesNotExist = fmt.Errorf("%w: key file does not exist", ErrPkg)
	// ErrServerUnableToMarshal error.
	ErrServerUnableToMarshal = fmt.Errorf("%w: unable to marshal", ErrPkg)
	// ErrServerShutdownFailure error.
	ErrServerShutdownFailure = fmt.Errorf("%w: unable to shutdown", ErrPkg)
)

// IsServerUnableToListenAndServe reports whether err is or wraps ErrServerUnableToListenAndServe.
func (s *DefaultServer) IsServerUnableToListenAndServe(err error) bool {
	return errors.Is(err, ErrServerUnableToListenAndServe)
}

// IsServerCertificateDoesNotExist reports whether err is or wraps ErrServerCertificateDoesNotExist.
func (s *DefaultServer) IsServerCertificateDoesNotExist(err error) bool {
	return errors.Is(err, ErrServerCertificateDoesNotExist)
}

// IsServerKeyDoesNotExist reports whether err is or wraps ErrServerKeyDoesNotExist.
func (s *DefaultServer) IsServerKeyDoesNotExist(err error) bool {
	return errors.Is(err, ErrServerKeyDoesNotExist)
}

// IsServerUnableToMarshal reports whether err is or wraps ErrServerUnableToMarshal.
func (s *DefaultServer) IsServerUnableToMarshal(err error) bool {
	return errors.Is(err, ErrServerUnableToMarshal)
}

// IsServerShutdownFailure reports whether err is or wraps ErrServerShutdownFailure.
func (s *DefaultServer) IsServerShutdownFailure(err error) bool {
	return errors.Is(err, ErrServerShutdownFailure)
}
