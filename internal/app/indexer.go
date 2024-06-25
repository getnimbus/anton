package app

import (
	"github.com/xssnick/tonutils-go/ton"

	"github.com/getnimbus/anton/internal/core/repository"
	"github.com/getnimbus/anton/internal/infra"
)

type IndexerConfig struct {
	DB *repository.DB

	PRODUCER infra.KafkaSyncProducer

	S3 repository.S3Service

	API ton.APIClientWrapped

	Fetcher FetcherService
	Parser  ParserService

	FromBlock uint32
	Workers   int
}

type IndexerService interface {
	Start() error
	Stop()
}
