package external

import "github.com/zokiio/mukabi/external/raiderio"

type ExternalAPI interface {
	RaiderIO() *raiderio.RaiderIO
}

type Services struct {
	raiderIO *raiderio.RaiderIO
}

func NewServices(RaiderIOKey string) *Services {
	return &Services{
		raiderIO: raiderio.New(RaiderIOKey),
	}
}

func (s *Services) RaiderIO() *raiderio.RaiderIO {
	return s.raiderIO
}
