package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type RequestInfo struct {
	URL      string
	Method   string
	Duration time.Duration
}

type ResponseInfo struct {
	Status   int
	BodySize int
}

type Logger struct {
	logger *slog.Logger
}

func NewMiddleware(logger *slog.Logger) *Logger {
	return &Logger{logger: logger}
}

// WithLogger добавление логов в запросы
func (m *Logger) WithLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestInfo := &RequestInfo{
			URL:    r.URL.String(),
			Method: r.Method,
		}
		responseInfo := &ResponseInfo{}
		writerWithLog := NewWriterWithLog(w, responseInfo)
		start := time.Now()

		h.ServeHTTP(writerWithLog, r)

		requestInfo.Duration = time.Since(start)
		m.logger.Debug("Info",
			slog.Group("Request",
				"URL", requestInfo.URL,
				"Method", requestInfo.Method,
				"Duration", requestInfo.Duration.String(),
			),
			slog.Group("Response",
				"Status", responseInfo.Status,
				"BodySize", responseInfo.BodySize),
		)
	})
}

// WriterWithLog структура для подмены writer для логирования
type WriterWithLog struct {
	w       http.ResponseWriter
	logInfo *ResponseInfo
}

func NewWriterWithLog(w http.ResponseWriter, responseInfo *ResponseInfo) *WriterWithLog {
	return &WriterWithLog{w: w, logInfo: responseInfo}
}

func (w *WriterWithLog) Write(body []byte) (int, error) {
	if w.logInfo.Status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	size, err := w.w.Write(body)
	w.logInfo.BodySize += size
	return size, err
}

func (w *WriterWithLog) WriteHeader(statusCode int) {
	if w.logInfo.Status != 0 {
		return
	}
	w.w.WriteHeader(statusCode)
	w.logInfo.Status = statusCode
}

func (w *WriterWithLog) Header() http.Header {
	return w.w.Header()
}
