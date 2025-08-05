package middleware

import (
	"bytes"
	"gateway-service/internal/model"
	"gateway-service/internal/repository"
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LogrusMiddleware(repo repository.LogRepository) echo.MiddlewareFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()

			// === BAGIAN BARU UNTUK MEMBACA BODY ===
			var requestBodyBytes []byte
			if req.Body != nil {
				requestBodyBytes, _ = io.ReadAll(req.Body)
			}
			// Kembalikan body ke request agar bisa dibaca lagi oleh handler selanjutnya
			req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
			// === AKHIR BAGIAN BARU ===

			err := next(c)

			stop := time.Now()
			latency := stop.Sub(start)
			res := c.Response()

			// Buat entri log terstruktur
			entry := logger.WithFields(logrus.Fields{
				"method":     req.Method,
				"uri":        req.RequestURI,
				"status":     res.Status,
				"remote_ip":  c.RealIP(),
				"latency":    latency.String(),
				"user_agent": req.UserAgent(),
			})

			// Buat model untuk disimpan ke DB
			logEntry := &model.ServiceLog{
				Timestamp:   stop,
				ServiceName: "gateway-service",
				Method:      req.Method,
				URI:         req.RequestURI,
				StatusCode:  res.Status,
				LatencyMs:   latency.Milliseconds(),
				RemoteIP:    c.RealIP(),
				UserAgent:   req.UserAgent(),
				RequestBody: string(requestBodyBytes),
			}

			if err != nil {
				entry = entry.WithError(err)
				logEntry.ErrorMessage = err.Error()
			}
			
			// Cetak log ke terminal (stdout)
			entry.Info("request processed")
			
			// Simpan log ke database
			repo.CreateLog(logEntry)

			return err
		}
	}
}