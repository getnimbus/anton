package parser

import (
	"encoding/base64"

	"github.com/getnimbus/anton/internal/app"
)

var _ app.ParserService = (*Service)(nil)

type Service struct {
	*app.ParserConfig

	bcConfigBase64 string
}

func NewService(cfg *app.ParserConfig) *Service {
	s := new(Service)
	s.ParserConfig = cfg
	s.bcConfigBase64 = base64.StdEncoding.EncodeToString(cfg.BlockchainConfig.ToBOC())
	return s
}
