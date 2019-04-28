package service

import (
	"github.com/zpnk/go-bitly"
	"go.uber.org/zap"
)

// Bitly wraps bitly client dependencies.
type Bitly struct {
	client *bitly.Client
	logger *zap.Logger
}

// NewBitly wraps dependencies for URL shortening.
func NewBitly(client *bitly.Client, logger *zap.Logger) *Bitly {
	return &Bitly{client: client, logger: logger}
}

// ShortURL shorten long URL
// / github.com/zpnk/go-bitly go bitly client
func (b *Bitly) ShortURL(longURL string) (string, error) {
	short, err := b.client.Links.Shorten(longURL)
	if err != nil {
		b.logger.Error("shorten-link-error", zap.Error(err))
		return "", err
	}
	return short.URL, nil
}
