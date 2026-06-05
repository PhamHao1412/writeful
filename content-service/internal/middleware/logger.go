package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		bodySize := c.Writer.Size()

		// Choose color based on status code
		statusColor := getStatusColor(statusCode)
		methodColor := getMethodColor(method)

		// Build log message
		logMsg := fmt.Sprintf("%s[%s]%s %s%s%s %s",
			colorCyan, time.Now().Format("2006-01-02 15:04:05"), colorReset,
			methodColor, method, colorReset,
			path,
		)

		if query != "" {
			logMsg += colorYellow + "?" + query + colorReset
		}

		logMsg += fmt.Sprintf(" | %s%d%s | %s | %v | %d bytes",
			statusColor, statusCode, colorReset,
			clientIP,
			duration,
			bodySize,
		)

		if errorMessage != "" {
			logMsg += colorRed + " | " + errorMessage + colorReset
		}

		fmt.Println(logMsg)
	}
}

func getStatusColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return colorGreen
	case statusCode >= 300 && statusCode < 400:
		return colorBlue
	case statusCode >= 400 && statusCode < 500:
		return colorYellow
	default:
		return colorRed
	}
}

func getMethodColor(method string) string {
	switch method {
	case "GET":
		return colorBlue
	case "POST":
		return colorGreen
	case "PUT":
		return colorYellow
	case "DELETE":
		return colorRed
	case "PATCH":
		return colorPurple
	default:
		return colorWhite
	}
}
