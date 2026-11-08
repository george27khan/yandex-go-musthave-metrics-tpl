package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipWriter структура подмены ResponseWriter со сжатием тела ответа
type GzipWriter struct {
	w              http.ResponseWriter
	zw             *gzip.Writer
	compressed     bool // Флаг активации сжатия
	firstWriteDone bool // Флаг первой записи тела ответа для установки заголовка
}

func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Write переопределение метода для запуска сжатия записываемых данных
func (gw *GzipWriter) Write(p []byte) (int, error) {
	for _, value := range gw.w.Header()["Content-Type"] {
		if value == "application/json" || value == "text/html" {
			if !gw.firstWriteDone {
				gw.firstWriteDone = true
				gw.compressed = true
				gw.w.Header().Set("Content-Encoding", "gzip")
			}
			return gw.zw.Write(p)
		}
	}
	return gw.w.Write(p)
}

func (gw *GzipWriter) Header() http.Header {
	return gw.w.Header()
}

func (gw *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 && gw.compressed {
		gw.w.Header().Set("Content-Encoding", "gzip") // если сжатие активно, то добавляем заголовок
	}
	gw.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (gw *GzipWriter) Close() error {
	return gw.zw.Close()
}

// GzipReader структура переопределения Request для распаковки тела от сжатия
type GzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewGzipReader(r io.ReadCloser) (*GzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &GzipReader{
		r:  r,
		zr: zr,
	}, nil
}

func (cr *GzipReader) Read(p []byte) (n int, err error) {
	return cr.zr.Read(p)
}

func (cr *GzipReader) Close() error {
	return cr.zr.Close()
}

// WithCompression middleware для добавления сжатия и распаковки тела
func WithCompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//если указан алгоритм сжатия ответа, то делаем сжатие тела
		responseEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(responseEncoding, "gzip") {
			cw := NewGzipWriter(w)
			// меняем тело запроса на новое
			w = cw
			defer cw.Close()
		}
		//если указан алгоритм сжатия запроса, то делаем распаковку тела
		requestEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(requestEncoding, "gzip") {
			cr, err := NewGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}
		next.ServeHTTP(w, r)

	})
}
