package middleware

import (
	"project-name/internal/entity"
	"project-name/internal/service"
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware(logService service.APILogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for the logs endpoint itself to avoid infinite loop
		if strings.HasPrefix(c.Request.URL.Path, "/api/logs") {
			c.Next()
			return
		}

		startTime := time.Now()

		// Read request body
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			// Restore the body for further processing
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Create custom response writer to capture response
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// Process request
		c.Next()

		// Calculate response time
		responseTime := time.Since(startTime).Milliseconds()

		// Determine access source from User-Agent
		userAgent := c.Request.UserAgent()
		accessFrom := determineAccessFrom(userAgent)

		// Get user ID, company ID, and branch ID if authenticated
		var userID *uuid.UUID
		var companyID *uuid.UUID
		var branchID *uuid.UUID
		
		if id, exists := c.Get("user_id"); exists {
			if uid, ok := id.(uuid.UUID); ok {
				userID = &uid
			}
		}
		
		if cid, exists := c.Get("company_id"); exists {
			if cidStr, ok := cid.(string); ok {
				if parsed, err := uuid.Parse(cidStr); err == nil {
					companyID = &parsed
				}
			}
		}
		
		if bid, exists := c.Get("branch_id"); exists {
			if bidStr, ok := bid.(string); ok {
				if parsed, err := uuid.Parse(bidStr); err == nil {
					branchID = &parsed
				}
			}
		}

		// Get response body (limit size to avoid huge logs)
		responseBody := blw.body.String()
		if len(responseBody) > 5000 {
			responseBody = responseBody[:5000] + "... (truncated)"
		}

		// Limit request body size
		if len(requestBody) > 5000 {
			requestBody = requestBody[:5000] + "... (truncated)"
		}

		// Get error message if any
		var errorMessage string
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		// Create log entry
		apiLog := &entity.APILog{
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			StatusCode:   c.Writer.Status(),
			ResponseTime: responseTime,
			IPAddress:    c.ClientIP(),
			UserAgent:    userAgent,
			AccessFrom:   accessFrom,
			UserID:       userID,
			CompanyID:    companyID,
			BranchID:     branchID,
			RequestBody:  requestBody,
			ResponseBody: responseBody,
			ErrorMessage: errorMessage,
		}

		// Save log asynchronously to avoid blocking the response
		go func() {
			_ = logService.CreateLog(apiLog)
		}()
	}
}

func determineAccessFrom(userAgent string) string {
	userAgent = strings.ToLower(userAgent)

	// Check for Postman
	if strings.Contains(userAgent, "postman") {
		return "postman"
	}

	// Check for mobile devices
	if strings.Contains(userAgent, "mobile") || 
	   strings.Contains(userAgent, "android") || 
	   strings.Contains(userAgent, "iphone") || 
	   strings.Contains(userAgent, "ipad") {
		return "mobile"
	}

	// Check for common API clients
	if strings.Contains(userAgent, "curl") {
		return "curl"
	}
	if strings.Contains(userAgent, "insomnia") {
		return "insomnia"
	}
	if strings.Contains(userAgent, "httpie") {
		return "httpie"
	}

	// Check for browsers
	if strings.Contains(userAgent, "mozilla") || 
	   strings.Contains(userAgent, "chrome") || 
	   strings.Contains(userAgent, "safari") || 
	   strings.Contains(userAgent, "firefox") || 
	   strings.Contains(userAgent, "edge") {
		return "website"
	}

	// Default
	return "unknown"
}
