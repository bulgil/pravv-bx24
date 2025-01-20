package server

import (
	"context"
	"net/http"

	"github.com/bulgil/pravv-bx24/app/internal/domain/models"
	"github.com/bulgil/pravv-bx24/app/internal/service/events"
	"github.com/bulgil/pravv-bx24/app/package/logger"
	"github.com/bulgil/pravv-bx24/app/package/queue"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*http.Server
	logger logger.Logger

	eventQueue *queue.Queue[models.Event]
}

type ServerOption struct {
	Host string
	Port string
}

func New(log logger.Logger, opts ServerOption) *Server {
	if opts.Host == "" || opts.Port == "" {
		panic("host and port must not be empty")
	}

	requestQueue := queue.New[models.Event]()
	engine := gin.Default()

	var server = &Server{
		Server:     &http.Server{Addr: opts.Host + ":" + opts.Port, Handler: engine},
		logger:     log,
		eventQueue: requestQueue,
	}

	engine.Use(putEventInQueueMiddleware(server))

	return server
}

func (s *Server) Run(ctx context.Context) {
	s.logger.Info("webhook server is running")

	go s.ListenAndServe()
	go s.serveQueue(ctx)

	<-ctx.Done()
}

func (s *Server) RegisterRoute(httpMethod, relativePath string, handlers ...gin.HandlerFunc) {
	h := s.Handler.(*gin.Engine)
	h.Handle(httpMethod, relativePath, handlers...)
}

func (s *Server) Event(e models.Event) {
	s.eventQueue.Push(e)
}

// TODO implement worker pool
func (s *Server) serveQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			event := s.eventQueue.Pop()
			if event.IsNil() {
				continue
			}

			events.EventRoute(s.logger, event)
		}
	}
}

func putEventInQueueMiddleware(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "OK")

		var event models.Event
		event.EventType = models.EventType(c.PostForm("event"))
		event.Data = models.Data(c.PostFormMap("data"))

		s.Event(event)
	}

}
