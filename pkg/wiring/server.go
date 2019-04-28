package wiring

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/zpnk/go-bitly"
	"go.uber.org/zap"

	"github.com/nikhil-github/sms-api/pkg/service"
)

// Start wires the services and start the app.
func Start(cfg *Config, logger *zap.Logger) error {

	ctx := context.Background()
	bitly := service.NewBitly(bitly.New(cfg.BITLY.Token), logger)
	svc := service.New(cfg.TRANSMIT.Apikey, cfg.TRANSMIT.Secret, &http.Client{Timeout: time.Second * 5}, logger, bitly)
	router := NewRouter(&Params{Logger: logger, Formatter: svc, Sender: svc})

	errs := make(chan error)
	serveHTTP(cfg.HTTP.Port, logger, router, errs)

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return nil
	}
}

func serveHTTP(port int, logger *zap.Logger, h http.Handler, errs chan error) {
	addr := fmt.Sprintf(":%d", port)
	s := &http.Server{Addr: addr, Handler: disableCors(h)}

	go func() {
		logger.Info("Listening for requests .....", zap.String("http.address", addr))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- errors.Wrapf(err, "error serving HTTP on address %s", addr)
		}
	}()
}

func disableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}
