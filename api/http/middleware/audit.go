package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"lbe/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// bodyLogWriter wraps gin.ResponseWriter and captures the response body.
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(data []byte) (int, error) {
	// write into our buffer
	w.body.Write(data)
	// write out to the real ResponseWriter
	return w.ResponseWriter.Write(data)
}

// if your handlers use WriteString:
func (w *bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func AuditLogger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// capture request body
		var reqBody string
		if c.Request.Body != nil {
			buf, _ := io.ReadAll(c.Request.Body)
			reqBody = string(buf)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
		}

		// wrap the ResponseWriter
		blw := &bodyLogWriter{body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
		c.Writer = blw

		// run the handler
		c.Next()

		// after handler: grab response
		respBody := blw.body.String()

		actor := c.GetString("app_id")
		// 2) if empty, fall back to the header
		if actor == "" {
			actor = c.GetHeader("AppID")
		}

		// build audit entry
		entry := model.AuditLog{
			ActorID:      actor,
			Method:       c.Request.Method,
			Path:         c.FullPath(),
			StatusCode:   c.Writer.Status(),
			ClientIP:     c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			RequestBody:  reqBody,
			ResponseBody: respBody,
			LatencyMs:    time.Since(start).Milliseconds(),
		}

		// persist asynchronously
		go func() {
			if err := db.Create(&entry).Error; err != nil {
				log.Printf("audit log persistence error: %v", err)
			}
		}()
	}
}
