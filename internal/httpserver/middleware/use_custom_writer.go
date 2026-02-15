package middleware

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

type bodyWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

// UseCustomWriter wraps the response writer to capture the response body for logging and error handling.
func UseCustomWriter() func(c *gin.Context) {
	return func(c *gin.Context) {
		buf := &bytes.Buffer{}
		writer := &bodyWriter{
			ResponseWriter: c.Writer,
			buf:            buf,
		}

		c.Writer = writer

		c.Next()
	}
}
