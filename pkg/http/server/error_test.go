package server

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultServer_IsPredicates(t *testing.T) {
	errBoom := errors.New("error")

	tests := map[string]struct {
		sentinel error
		check    func(s *DefaultServer, err error) bool
	}{
		"IsServerUnableToListenAndServe": {
			sentinel: ErrServerUnableToListenAndServe,
			check:    (*DefaultServer).IsServerUnableToListenAndServe,
		},
		"IsServerCertificateDoesNotExist": {
			sentinel: ErrServerCertificateDoesNotExist,
			check:    (*DefaultServer).IsServerCertificateDoesNotExist,
		},
		"IsServerKeyDoesNotExist": {
			sentinel: ErrServerKeyDoesNotExist,
			check:    (*DefaultServer).IsServerKeyDoesNotExist,
		},
		"IsServerUnableToMarshal": {
			sentinel: ErrServerUnableToMarshal,
			check:    (*DefaultServer).IsServerUnableToMarshal,
		},
		"IsServerShutdownFailure": {
			sentinel: ErrServerShutdownFailure,
			check:    (*DefaultServer).IsServerShutdownFailure,
		},
	}

	srvr := &DefaultServer{}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.True(t, test.check(srvr, test.sentinel), "bare sentinel")
			require.True(
				t,
				test.check(srvr, fmt.Errorf("%w: %w", test.sentinel, errBoom)),
				"wrapped sentinel",
			)
			require.False(t, test.check(srvr, errBoom), "unrelated error")
			require.False(t, test.check(srvr, nil), "nil error")
		})
	}
}
