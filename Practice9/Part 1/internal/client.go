package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ExecutePayment(ctx context.Context, client *http.Client, url string) error {
	maxRetries := 5

	for attempt := 1; attempt <= maxRetries; attempt++ {

		req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)

		
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("Attempt %d: Success!\n", attempt)
			return nil
		}

		
		if !IsRetryable(resp, err) {
			return fmt.Errorf("non-retryable error: %v", err)
		}

		if attempt == maxRetries {
			return fmt.Errorf("max retries reached")
		}

		delay := CalculateBackoff(attempt)
		fmt.Printf("Attempt %d failed: waiting %v...\n", attempt, delay)

		
		select {
		case <-time.After(delay):
			
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}