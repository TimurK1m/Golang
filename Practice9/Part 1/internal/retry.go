package internal

import (
	"math/rand"
	"net"
	"net/http"
	"time"
)


func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			return true
		}
		return true
	}

	if resp == nil {
		return false
	}

	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	case 401, 404:
		return false
	default:
		return false
	}
}


func CalculateBackoff(attempt int) time.Duration {
	base := 500 * time.Millisecond
	maxDelay := base * time.Duration(1<<attempt)

	return time.Duration(rand.Int63n(int64(maxDelay)))
}