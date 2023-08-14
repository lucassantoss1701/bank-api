package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"lucassantoss1701/bank/internal/entity"

	"github.com/sirupsen/logrus"
)

var (
	once   sync.Once
	logger *logrus.Logger
)

func getLoggerInstance() *logrus.Logger {
	once.Do(func() {
		logger = logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
			ForceQuote:    true,
		})
		logger.SetOutput(os.Stdout)
	})
	return logger
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "swagger") {
			next.ServeHTTP(w, r)
			return
		}

		uuid := entity.NewUUID()

		bodyCopy, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(bodyCopy))

		rw := &responseWriter{w, &responseData{status: http.StatusInternalServerError}}

		fields := logrus.Fields{
			"url":          r.URL.Path,
			"method":       r.Method,
			"request_body": string(bodyCopy),
		}

		logger := getLoggerInstance()

		logger.WithFields(fields).Info(fmt.Sprintf("Requisição recebida: %s", uuid))

		defer func() {
			fields["response_body"] = rw.ResponseData.buf.String()
			fields["response_code"] = rw.ResponseData.status

			logger.WithFields(fields).Info(fmt.Sprintf("Requisição concluída: %s", uuid))
		}()

		next.ServeHTTP(rw, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	ResponseData *responseData
}

type responseData struct {
	status int
	buf    bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
	rw.ResponseData.status = code
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.ResponseData.buf.Write(b)
	return size, err
}
