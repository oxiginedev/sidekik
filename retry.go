package sidekik

import (
	"context"
	"fmt"
	"time"
)

type RetryOptions struct {
	// Tries is the number of times the function should be retried
	Tries uint
	// Backoff is the exponential backoff strategy used for retrying
	// Both Tries and SleepSeconds will be ignored if a value for Backoff is set
	Backoff []uint
	// SleepSeconds is the number of seconds to wait before retrying
	// This will be ignored if a value for Backoff is set
	SleepSeconds uint
	// RetryableFunc is a function that returns true if the error is retryable
	// If RetryableFunc is not set, all errors are considered retryable
	RetryableFunc func(uint, error) bool
}

// Retry retries a given function based on the given options
func Retry(ctx context.Context, fn func() error, opts RetryOptions) error {
	tries := opts.Tries
	backoff := opts.Backoff

	if len(backoff) > 0 {
		tries = uint(len(backoff)) + 1
	} else if tries == 0 {
		tries = 1
	}

	sleepSeconds := opts.SleepSeconds
	if sleepSeconds == 0 {
		sleepSeconds = 1
	}

	var err error
	for attempts := uint(0); attempts < tries; attempts++ {
		if err = ctx.Err(); err != nil {
			return fmt.Errorf("retry cancelled after %d attempt(s): %w", attempts, err)
		}

		if attempts > 0 {
			delay := sleepSeconds
			if len(backoff) > 0 {
				delay = backoff[attempts-1]
			}

			// Use select to respect context cancellation during sleep
			timer := time.NewTimer(time.Duration(delay) * time.Second)
			select {
			case <-ctx.Done():
				timer.Stop()
				return fmt.Errorf("retry cancelled after %d attempt(s): %w", attempts, ctx.Err())
			case <-timer.C:
			}
		}

		if err = fn(); err != nil {
			if opts.RetryableFunc != nil && !opts.RetryableFunc(attempts, err) {
				return err
			}
			continue
		}

		return nil
	}

	return err
}
