package logger

import (
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	DebugLevel = -4
	InfoLevel  = 0
	WarnLevel  = 4
	ErrorLevel = 8
)

var logger *slog.Logger

func init() {
	Init(os.Stdout, InfoLevel)
}

func Init(wrt io.Writer, level int) {
	if wrt == nil {
		wrt = os.Stdout
	}
	logger = slog.New(slog.NewJSONHandler(wrt, &slog.HandlerOptions{Level: slog.Level(level)}))
}

func With(args ...any) *slog.Logger {
	return logger.With(args...)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func GetRequestLogger(r *http.Request) *slog.Logger {
	reqID := r.Context().Value(middleware.RequestIDKey)
	if reqID == nil {
		reqID = ""
	}
	return With("id", reqID.(string))
}
