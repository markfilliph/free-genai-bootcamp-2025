package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"lang-portal/internal/logger"
)

// responseWriter is a custom response writer that captures the response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body while writing it
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestLogging middleware logs incoming HTTP requests and their responses
func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore the request body for later use
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create custom response writer to capture response
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:          &bytes.Buffer{},
		}
		c.Writer = w

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Prepare log fields
		fields := logger.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration.String(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		// Add request body for non-GET requests (limited to 1000 characters)
		if c.Request.Method != "GET" && len(requestBody) > 0 {
			reqBody := string(requestBody)
			if len(reqBody) > 1000 {
				reqBody = reqBody[:1000] + "..."
			}
			fields["request_body"] = reqBody
		}

		// Add response body for errors (limited to 1000 characters)
		if c.Writer.Status() >= 400 {
			respBody := w.body.String()
			if len(respBody) > 1000 {
				respBody = respBody[:1000] + "..."
			}
			fields["response_body"] = respBody
		}

		// Log based on status code
		if c.Writer.Status() >= 500 {
			logger.Error("Server error", fields)
		} else if c.Writer.Status() >= 400 {
			logger.Warn("Client error", fields)
		} else {
			logger.Info("Request completed", fields)
		}
	}
}
