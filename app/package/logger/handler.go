package logger

import (
	"context"
	"io"
	"log/slog"
	"sync"
)

type Handler struct {
	w    io.Writer
	opts *slog.HandlerOptions

	mu    *sync.Mutex
	attrs []slog.Attr
}

func NewHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	return &Handler{
		w:     w,
		opts:  opts,
		mu:    &sync.Mutex{},
		attrs: make([]slog.Attr, 0),
	}
}

func (h Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h Handler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	attrs := h.attrs
	h.mu.Unlock()

	level := "[" + r.Level.String() + "]"
	timestamp := r.Time.Format("[02/01 15:04:05.000]")
	message := r.Message

	var fields string

	for _, a := range attrs {
		fields += " " + a.Key + "=" + a.Value.String()
	}

	r.Attrs(func(a slog.Attr) bool {
		fields += " " + a.Key + "=" + a.Value.String()
		return true
	})

	_, err := h.w.Write([]byte(level + "\t" + timestamp + " " + message + fields + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()

	return Handler{
		w:     h.w,
		opts:  h.opts,
		mu:    h.mu,
		attrs: append(h.attrs, attrs...),
	}
}

func (h Handler) WithGroup(name string) slog.Handler {
	// TODO implement

	return h
}
