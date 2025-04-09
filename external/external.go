// Package external provides interfaces for external service integrations
package external

import "github.com/zokiio/mukabi/external/raiderio"

// ExternalAPI defines the interface for accessing external services
type ExternalAPI interface {
	RaiderIO() *raiderio.Client
}

// Services implements ExternalAPI interface and holds external service clients
type Services struct {
	raiderIO *raiderio.Client
}

// NewServices creates a new Services instance with configured external clients
func NewServices(raiderIOKey string) *Services {
	return &Services{
		raiderIO: raiderio.New(raiderIOKey),
	}
}

// RaiderIO returns the RaiderIO client instance
func (s *Services) RaiderIO() *raiderio.Client {
	return s.raiderIO
}
