package gin_accesslog

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"time"

	"go.uber.org/zap"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Ginzap(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 获取请求体
		requestBody, _ := c.GetRawData()

		// 请求包体写回。
		if len(requestBody) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		crw := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = crw

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}
		respBody := string(crw.body.Bytes())

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		}

		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(timeFormat)),
			zap.String("requestBody", string(requestBody)),
			zap.String("responseBody", respBody),
			zap.Duration("latency", latency),
		)
	}
}
