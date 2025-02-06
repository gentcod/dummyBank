package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// responseInterceptor is a custom ResponseWriter to capture response body
type responseInterceptor struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseInterceptor) Write(data []byte) (int, error) {
	return r.body.Write(data)
}

func (r *responseInterceptor) WriteString(s string) (int, error) {
	return r.body.WriteString(s)
}

// ResponseInterceptor is a middleware to intercept and modify response in production.
// The response is set to be custom and generic
func interceptor(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		originalWriter := c.Writer

		interceptor := &responseInterceptor{
			ResponseWriter: originalWriter,
			body:           bytes.NewBuffer([]byte{}),
		}
		c.Writer = interceptor

		c.Next()

		statusCode := c.Writer.Status()
		if server.config.Environment == "development" && statusCode >= 500 {
			for k := range originalWriter.Header() {
				originalWriter.Header().Del(k)
			}

			originalWriter.Header().Set("Content-Type", "application/json")
			originalWriter.WriteHeader(statusCode)

			response := handleInternalResponse(ApiResponse[any]{
				StatusCode: statusCode,
				Data:       nil,
			})

			json.NewEncoder(originalWriter).Encode(response)
			c.Abort()
			return
		}

		originalWriter.Write(interceptor.body.Bytes())
	}
}
