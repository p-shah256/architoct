package logger

import (
	"time"

	"github.com/labstack/echo/v4"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        start := time.Now()
        req := c.Request()
        res := c.Response()

        err := next(c)

        if err == nil && res.Status < 400 {
            L.Info().
                Int("status", res.Status).
                Str("method", req.Method).
                Str("path", req.URL.Path).
                Dur("latency", time.Since(start)).
                Str("ip", c.RealIP()).
                Send()
        } else {
            errLogger := L.Error()
            if err != nil {
                errLogger = errLogger.Err(err)
            }
            errLogger.
                Int("status", res.Status).
                Str("method", req.Method).
                Str("path", req.URL.Path).
                Str("user_agent", req.UserAgent()).
                Str("remote_ip", c.RealIP()).
                Dur("latency", time.Since(start)).
                Interface("headers", req.Header).
                Send()
        }
        return err
    }
}
