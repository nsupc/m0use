package ns

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	client             http.Client
	user               string
	requests           []time.Time
	ratelimitLimit     int
	ratelimitRemaining int
	ratelimitResetIn   time.Duration
	maxRequests        int
}

func (c *Client) clearBucket() {
	now := time.Now()

	filtered := []time.Time{}

	for _, instant := range c.requests {
		if now.Sub(instant) <= 30*time.Second {
			filtered = append(filtered, instant)
		}
	}

	c.requests = filtered
}

func (c *Client) acquireFatal() error {
	c.clearBucket()

	if len(c.requests) >= c.maxRequests {
		return errors.New("too many requests")
	}

	if c.ratelimitRemaining <= 1 {
		return errors.New("too many requests")
	}

	now := time.Now()

	c.requests = append(c.requests, now)

	return nil
}

func (c *Client) acquire() error {
	c.clearBucket()

	if len(c.requests) == 0 {
		c.requests = append(c.requests, time.Now())
		return nil
	}

	sleepDuration := float64(c.maxRequests) / 30

	slog.Debug("sleeping", slog.Float64("duration", sleepDuration))
	time.Sleep(time.Duration(sleepDuration))

	c.requests = append(c.requests, time.Now())

	return nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	err := c.acquire()
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.user)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	limit := resp.Header.Get("ratelimit-limit")
	if limit != "" {
		limit, err := strconv.Atoi(limit)
		if err != nil {
			slog.Warn("failed to convert ratelimit-limit to int", slog.Any("error", err))
		} else {
			c.ratelimitLimit = limit
		}

	}

	remaining := resp.Header.Get("ratelimit-remaining")
	if remaining != "" {
		remaining, err := strconv.Atoi(remaining)
		if err != nil {
			slog.Warn("failed to convert ratelimit-remaining to int", slog.Any("error", err))
		} else {
			c.ratelimitRemaining = remaining
		}
	}

	reset := resp.Header.Get("ratelimit-reset")
	if reset != "" {
		reset, err := strconv.Atoi(reset)
		if err != nil {
			slog.Warn("failed to convert ratelimit-reset to int", slog.Any("error", err))
		} else {
			c.ratelimitResetIn = time.Duration(reset)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("request failed with status %s", resp.Status))
	}

	return resp, nil
}

func NewClient(user string, maxRequests int) *Client {
	client := Client{
		client:             http.Client{},
		user:               user,
		requests:           []time.Time{},
		ratelimitLimit:     50,
		ratelimitRemaining: 50,
		ratelimitResetIn:   30 * time.Second,
		maxRequests:        maxRequests,
	}

	return &client
}
