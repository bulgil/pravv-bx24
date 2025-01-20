package client

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/bulgil/pravv-bx24/app/package/logger"
	"github.com/bulgil/pravv-bx24/app/package/queue"
)

const (
	decrementInterval = time.Second
)

var (
	errDoRequestTimeout = errors.New("do request timeout")
)

type Client struct {
	*http.Client
	logger logger.Logger

	bucket *bucket

	timeout time.Duration
}

type ClientOpts struct {
	Timeout          time.Duration
	CounterMax       int
	CounterDecrement int
}

func New(log logger.Logger, opts ClientOpts) *Client {
	return &Client{
		Client:  &http.Client{},
		logger:  log,
		bucket:  newBucket(opts.CounterMax, opts.CounterDecrement),
		timeout: opts.Timeout,
	}
}

func (c *Client) Run(ctx context.Context) {
	go c.decrementRequestCounter(ctx, decrementInterval)
	go c.serveQueue(ctx)

	c.logger.Info("request client is running")

	<-ctx.Done()
}

// DoRequest pushes request in client queue. Then client will process it in goroutine. If timeout is exceeded, error errDoRequestTimeout is returned.
func (c *Client) DoRequest(r *http.Request) (*http.Response, error) {
	respChan := make(chan *http.Response)
	errChan := make(chan error)
	defer func() {
		close(respChan)
		close(errChan)
	}()

	c.bucket.Queue.Push(&Request{
		Request: r,
		Resp:    respChan,
		Error:   errChan,
	})

	select {
	case res := <-respChan:
		return res, nil

	case err := <-errChan:
		return nil, err

	case <-time.After(c.timeout):
		return nil, errDoRequestTimeout
	}
}

func (c *Client) serveQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			b := c.bucket
			if b.counter <= b.counterMax {
				req := b.Pop()
				if req == nil {
					continue
				}

				go c.request(*req)
			}
		}
	}
}

func (c *Client) request(r Request) {
	log := c.logger
	c.bucket.increment()

	res, err := c.Do(r.Request)
	if err != nil {
		r.Error <- err
		return
	}

	log.Debug("request processed", slog.Attr{
		Key:   "status",
		Value: slog.StringValue(res.Status),
	}, slog.Attr{
		Key:   "request_counter",
		Value: slog.IntValue(c.bucket.Counter()),
	})

	r.Resp <- res
}

func (c *Client) decrementRequestCounter(ctx context.Context, interval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(interval)
			c.bucket.decrement()
		}
	}
}

type bucket struct {
	*queue.Queue[*Request]

	mu               *sync.Mutex
	counter          int
	counterMax       int
	counterDecrement int
}

func newBucket(counterMax, counterDecrement int) *bucket {
	return &bucket{
		Queue:            queue.New[*Request](),
		mu:               &sync.Mutex{},
		counter:          0,
		counterMax:       counterMax,
		counterDecrement: counterDecrement,
	}
}

func (b *bucket) increment() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.counter++
}

func (b *bucket) decrement() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.counter >= b.counterDecrement {
		b.counter -= b.counterDecrement
	} else {
		b.counter = 0
	}
}

func (b *bucket) Counter() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.counter
}
