package resilience

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewRetrier(t *testing.T) {
	tests := []struct {
		name          string
		options       []Option
		wantBaseDelay time.Duration
		wantMaxDelay  time.Duration
		wantRetries   int
	}{
		{
			name:          "defaults are applied when no options are given",
			options:       nil,
			wantBaseDelay: BaseDelayDefault,
			wantMaxDelay:  MaxDelayDefault,
			wantRetries:   RetriesDefault,
		},
		{
			name: "options override defaults",
			options: []Option{
				WithBaseDelay(5 * time.Millisecond),
				WithMaxDelay(50 * time.Millisecond),
				WithRetries(7),
			},
			wantBaseDelay: 5 * time.Millisecond,
			wantMaxDelay:  50 * time.Millisecond,
			wantRetries:   7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrier := NewRetrier(tt.options...)

			require.Equal(t, tt.wantBaseDelay, retrier.baseDelay)
			require.Equal(t, tt.wantMaxDelay, retrier.maxDelay)
			require.Equal(t, tt.wantRetries, retrier.retries)
		})
	}
}

func TestRetrier_delay(t *testing.T) {
	tests := []struct {
		name      string
		baseDelay time.Duration
		maxDelay  time.Duration
		attempt   int
		want      time.Duration
	}{
		{
			name:      "first attempt has no delay",
			baseDelay: 10 * time.Millisecond,
			maxDelay:  100 * time.Millisecond,
			attempt:   0,
			want:      0,
		},
		{
			name:      "negative attempt has no delay",
			baseDelay: 10 * time.Millisecond,
			maxDelay:  100 * time.Millisecond,
			attempt:   -1,
			want:      0,
		},
		{
			name:      "second attempt doubles base delay",
			baseDelay: 10 * time.Millisecond,
			maxDelay:  100 * time.Millisecond,
			attempt:   1,
			want:      20 * time.Millisecond,
		},
		{
			name:      "third attempt doubles again",
			baseDelay: 10 * time.Millisecond,
			maxDelay:  100 * time.Millisecond,
			attempt:   2,
			want:      40 * time.Millisecond,
		},
		{
			name:      "fourth attempt doubles again",
			baseDelay: 10 * time.Millisecond,
			maxDelay:  100 * time.Millisecond,
			attempt:   3,
			want:      80 * time.Millisecond,
		},
		{
			name:      "fifth attempt is capped at max delay",
			baseDelay: 10 * time.Millisecond,
			maxDelay:  100 * time.Millisecond,
			attempt:   4,
			want:      100 * time.Millisecond,
		},
		{
			name:      "shift overflow is capped at max delay",
			baseDelay: time.Second,
			maxDelay:  time.Minute,
			attempt:   100,
			want:      time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrier := NewRetrier(WithBaseDelay(tt.baseDelay), WithMaxDelay(tt.maxDelay))

			require.Equal(t, tt.want, retrier.delay(tt.attempt))
		})
	}
}

func TestRetrier_Execute(t *testing.T) {
	errBoom := errors.New("error")

	tests := []struct {
		name           string
		retrier        *Retrier
		ctx            func() (context.Context, context.CancelFunc)
		handler        func(calls *int) HandlerFunc
		wantErrIs      []error
		wantCalls      int
		wantMaxElapsed time.Duration
	}{
		{
			name: "succeeds on first attempt without waiting",
			retrier: NewRetrier(
				WithBaseDelay(50*time.Millisecond),
				WithMaxDelay(time.Second),
				WithRetries(3),
			),
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			handler: func(calls *int) HandlerFunc {
				return func(_ context.Context) error {
					*calls++

					return nil
				}
			},
			wantErrIs:      []error{nil},
			wantCalls:      1,
			wantMaxElapsed: 25 * time.Millisecond,
		},
		{
			name: "succeeds after a couple of failed attempts",
			retrier: NewRetrier(
				WithBaseDelay(time.Millisecond),
				WithMaxDelay(10*time.Millisecond),
				WithRetries(5),
			),
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			handler: func(calls *int) HandlerFunc {
				return func(_ context.Context) error {
					*calls++
					if *calls < 3 {
						return errBoom
					}

					return nil
				}
			},
			wantErrIs: []error{nil},
			wantCalls: 3,
		},
		{
			name: "returns ErrRetriesExhausted when the handler always fails",
			retrier: NewRetrier(
				WithBaseDelay(time.Millisecond),
				WithMaxDelay(5*time.Millisecond),
				WithRetries(3),
			),
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			handler: func(calls *int) HandlerFunc {
				return func(_ context.Context) error {
					*calls++

					return errBoom
				}
			},
			wantErrIs: []error{ErrRetriesExhausted, errBoom},
			wantCalls: 3,
		},
		{
			name:    "returns ErrRetriesExhausted immediately when retries is zero",
			retrier: NewRetrier(WithRetries(0)),
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			handler: func(calls *int) HandlerFunc {
				return func(_ context.Context) error {
					*calls++

					return nil
				}
			},
			wantErrIs: []error{ErrRetriesExhausted},
			wantCalls: 0,
		},
		{
			name: "returns ErrRetryWaitCanceled when the context is done during backoff",
			retrier: NewRetrier(
				WithBaseDelay(200*time.Millisecond),
				WithMaxDelay(200*time.Millisecond),
				WithRetries(3),
			),
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 20*time.Millisecond)
			},
			handler: func(calls *int) HandlerFunc {
				return func(_ context.Context) error {
					*calls++

					return errBoom
				}
			},
			wantErrIs: []error{ErrRetryWaitCanceled, context.DeadlineExceeded},
			wantCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.ctx()
			defer cancel()

			calls := 0

			start := time.Now()
			err := tt.retrier.Execute(ctx, tt.handler(&calls))
			elapsed := time.Since(start)

			for _, wantErr := range tt.wantErrIs {
				require.ErrorIs(t, err, wantErr)
			}

			require.Equal(t, tt.wantCalls, calls)

			if tt.wantMaxElapsed > 0 {
				require.Less(t, elapsed, tt.wantMaxElapsed, "attempt should not have waited")
			}
		})
	}
}

func TestRetrier_IsRetriesExhausted(t *testing.T) {
	errBoom := errors.New("error")

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "true for the bare sentinel",
			err:  ErrRetriesExhausted,
			want: true,
		},
		{
			name: "true when wrapped with the last handler error",
			err:  fmt.Errorf("%w: %w", ErrRetriesExhausted, errBoom),
			want: true,
		},
		{
			name: "false for an unrelated error",
			err:  errBoom,
			want: false,
		},
		{
			name: "false for a different sentinel",
			err:  ErrRetryWaitCanceled,
			want: false,
		},
		{
			name: "false for nil",
			err:  nil,
			want: false,
		},
	}

	retrier := NewRetrier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, retrier.IsRetriesExhausted(tt.err))
		})
	}
}

func TestRetrier_IsRetryWaitCanceled(t *testing.T) {
	errBoom := errors.New("error")

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "true for the bare sentinel",
			err:  ErrRetryWaitCanceled,
			want: true,
		},
		{
			name: "true when wrapped with the context error",
			err:  fmt.Errorf("%w: %w", ErrRetryWaitCanceled, context.DeadlineExceeded),
			want: true,
		},
		{
			name: "false for an unrelated error",
			err:  errBoom,
			want: false,
		},
		{
			name: "false for a different sentinel",
			err:  ErrRetriesExhausted,
			want: false,
		},
		{
			name: "false for nil",
			err:  nil,
			want: false,
		},
	}

	retrier := NewRetrier()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, retrier.IsRetryWaitCanceled(tt.err))
		})
	}
}
