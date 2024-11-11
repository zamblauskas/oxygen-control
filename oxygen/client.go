package oxygen

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Boost() error {
	return c.withRetry(func() error {
		log.Info().Msg("triggering oxygen boost")
		url := c.baseURL + "/cmd?sc73=1"
		resp, err := c.client.Get(url)
		if err != nil {
			return &retryableError{
				message: fmt.Sprintf("http request failed: %v", err),
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 {
			body, _ := io.ReadAll(resp.Body)
			return &retryableError{
				message: fmt.Sprintf("server error: %s, body: %s", resp.Status, string(body)),
			}
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			log.Warn().
				Str("url", url).
				Str("status", resp.Status).
				Str("body", string(body)).
				Msg("non-2xx HTTP from url")
			return fmt.Errorf("unexpected status: %s", resp.Status)
		}

		log.Info().Msg("Oxygen boost triggered successfully")
		return nil
	})
}

type retryableError struct {
	message string
}

func (e *retryableError) Error() string {
	return e.message
}

func (c *Client) withRetry(operation func() error) error {
	maxRetries := 3
	backoff := 1 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		// Check if error is retryable
		if _, isRetryable := err.(*retryableError); !isRetryable {
			return err
		}

		log.Error().
			Int("attempt", attempt+1).
			Int("max_retries", maxRetries+1).
			Err(err).
			Msg("attempt failed")
		if attempt == maxRetries {
			return err
		}

		time.Sleep(backoff)
		backoff *= 2
	}
	return fmt.Errorf("max retries exceeded")
}
