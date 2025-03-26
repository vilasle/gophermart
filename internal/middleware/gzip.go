package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vilasle/gophermart/internal/logger"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzMW(h http.Handler) http.Handler {
	gzipFunc := func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetRequestLogger(req)

		// copy original request
		or := req

		// check, that the client has sent to the server compressed data in gzip format
		contentEncoding := req.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			fmt.Println("[INFO] Content-Encoding is gzip")
			if !(strings.Contains(req.Header.Get("Content-Type"), "application/json") || strings.Contains(req.Header.Get("Content-Type"), "text/html")) {

				log.Info("Unacceptable Content-Type => continue without gzip")

				h.ServeHTTP(res, or)
				return
			}

			acceptEncoding := req.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if !supportsGzip {

				h.ServeHTTP(res, or)
				return
			}
			// OK
			// wrap request body into io.Reader with decompression available
			cr, err := newCompressReader(req.Body)
			if err != nil {
				log.Error("can not decompress request", "error", err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			req.Body = cr
			defer cr.Close()
		}

		compressedResponse := newCompressWriter(res)

		defer compressedResponse.Close()

		// call handler with modified res and req
		h.ServeHTTP(compressedResponse, req)

	}
	return http.HandlerFunc(gzipFunc)
}
