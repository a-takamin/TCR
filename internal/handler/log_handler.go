package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func LogMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		// TODO: もっといい方法がある？
		if c.Request.URL.Path != "/health" {

			slog.Info("request metadata",
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				slog.Group("header",
					slog.String("Host", c.Request.Header.Get("Host")),
					slog.String("Content-Type", c.Request.Header.Get("Content-Type")),
					slog.String("Content-Length", c.Request.Header.Get("Content-Length")),
					slog.String("Range", c.Request.Header.Get("Range")),
				),
			)

		}

		c.Next()

		if c.Request.URL.Path != "/health" {

			slog.Info("response metadata",
				slog.Group("header",
					slog.String("Docker-Content-Digest", c.Writer.Header().Get("Docker-Content-Digest")),
					slog.String("Range", c.Writer.Header().Get("Range")),
					slog.String("Content-Range", c.Writer.Header().Get("Content-Range")),
				),
			)
		}

	}
}
