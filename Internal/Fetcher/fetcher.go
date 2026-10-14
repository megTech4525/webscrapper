package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Fetcher struct {
	client    *http.Client
	userAgent string
	mu        sync.Mutex
	limiters  map[string]*rate.Limiter
	perSecond rate.Limit
	burst     int
}

func New(useAgent string, requestsPerSecond float64) *Fetcher {
	return &Fetcher{
		client:    &http.Client{Timeout: 15 * time.Second},
		userAgent: useAgent,
		limiters:  make(map[string]*rate.Limiter),
		perSecond: rate.Limit(requestsPerSecond),
		burst:     10,
	}
}

//creates a limiter for a host, creating it for the first time

func (f *Fetcher) limitFor(host string) *rate.Limiter {
	f.mu.Lock()

	lim, ok := f.limiters[host]
	if !ok {
		lim = rate.NewLimiter(f.perSecond, f.burst)
		f.limiters[host] = lim
	}
	return lim
}

func (f *Fetcher) Fetch(ctx context.Context, rawURL string) ([]byte, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("bad url %q: %w", rawURL, err)

	}
	if err := f.limitFor(u.Host).Wait(ctx); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", f.userAgent)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s, returned status %d", rawURL, resp.StatusCode)

	}
	return io.ReadAll(resp.Body)
}
