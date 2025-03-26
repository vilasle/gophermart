package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		start := time.Now()
		// mold collections to fill the structure
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		// fill the custom logging ResponseWriter
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		// serve an original request with custom ResponseWriter
		next.ServeHTTP(&lw, r)
		// get request duration
		duration := time.Since(start)

		id := r.Context().Value(middleware.RequestIDKey).(string)

		logger.Info("received request",
			"id", id,
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)

	})
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// structure to save response  info
type (
	responseData struct {
		status int
		size   int
	}
	// add the realization of http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// redefine the methods to get needed response data
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// get response using original http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	// new functionality
	r.responseData.size += size // get the size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// get statusCOde using original http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	// new functionality
	r.responseData.status = statusCode // get codeStatus
}
