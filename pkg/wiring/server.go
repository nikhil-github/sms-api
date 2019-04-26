package wiring

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/nikhil-github/sms-api/pkg/service"
)

// Start wires the services and start the app.
func Start(cfg *Config, logger *zap.Logger) error {

	ctx := context.Background()
	svc := service.New(cfg.TRANSMIT.Apikey, cfg.TRANSMIT.Secret, http.DefaultClient, logger, cfg.BITLY.Token)
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
	s := &http.Server{Addr: addr, Handler: h}

	go func() {
		logger.Info("Listening for requests .....", zap.String("http.address", addr))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- errors.Wrapf(err, "error serving HTTP on address %s", addr)
		}
	}()
}
