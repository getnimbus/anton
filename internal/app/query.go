package app

import (
	"context"

	"github.com/xssnick/tonutils-go/ton"

	"github.com/getnimbus/anton/abi"
	"github.com/getnimbus/anton/internal/core"
	"github.com/getnimbus/anton/internal/core/aggregate"
	"github.com/getnimbus/anton/internal/core/aggregate/history"
	"github.com/getnimbus/anton/internal/core/filter"
	"github.com/getnimbus/anton/internal/core/repository"
	"github.com/getnimbus/anton/internal/infra"
)

type QueryConfig struct {
	DB *repository.DB

	PRODUCER infra.KafkaSyncProducer

	API ton.APIClientWrapped
}

type QueryService interface {
	GetStatistics(ctx context.Context) (*aggregate.Statistics, error)

	GetDefinitions(context.Context) (map[abi.TLBType]abi.TLBFieldsDesc, error)
	GetInterfaces(ctx context.Context) ([]*core.ContractInterface, error)
	GetOperations(ctx context.Context) ([]*core.ContractOperation, error)

	filter.BlockRepository

	GetLabelCategories(context.Context) ([]core.LabelCategory, error)

	filter.AccountRepository
	filter.TransactionRepository
	filter.MessageRepository

	aggregate.AccountRepository
	aggregate.MessageRepository

	history.AccountRepository
	history.TransactionRepository
	history.MessageRepository
}
