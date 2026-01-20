package sidekik

import "time"

type RetryOptions struct {
	// Tries is the number of times the function should be retried
	Tries uint
	// Backoff is the exponential backoff strategy used for retrying
	// Both Tries and SleepSeconds will be ignored if a value for Backoff is set
	Backoff []uint
	// SleepSeconds is the number of seconds to wait before retrying
	// This will be ignored if a value for Backoff is set
	SleepSeconds uint
}

// Retry retries a given function based on the given options
func Retry(opts RetryOptions, fn func() error) error {
	tries := opts.Tries
	backoff := opts.Backoff

	if tries == 0 && len(backoff) == 0 {
		tries = 1
	}

	if len(backoff) > 0 {
		tries = uint(len(backoff)) + 1
	}

	sleepSeconds := opts.SleepSeconds
	if sleepSeconds == 0 {
		sleepSeconds = 1
	}

	var err error
	for attempts := uint(0); attempts < tries; attempts++ {
		if attempts > 0 {
			delay := sleepSeconds
			if len(backoff) > 0 {
				delay = backoff[attempts-1]
			}

			time.Sleep(time.Duration(delay) * time.Second)
		}

		if err = fn(); err != nil {
			continue
		}

		return nil
	}

	return err
}
