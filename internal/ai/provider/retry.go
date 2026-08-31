package provider

import (
	"context"
	"fmt"
	"time"
)

// Executes an AI operation with retry support.
func retry(
	ctx context.Context,
	maxRetries int,
	operation func(context.Context) error,
) error {

	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		lastErr = operation(ctx)
		if lastErr == nil {
			return nil
		}

		// No delay after the final attempt.
		if attempt == maxRetries {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt+1) * 200 * time.Millisecond):
		}
	}

	return fmt.Errorf("provider retry exhausted: %w", lastErr)
}