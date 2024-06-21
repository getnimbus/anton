package indexer

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/getnimbus/anton/internal/app"
	"github.com/getnimbus/anton/internal/core"
	"github.com/getnimbus/anton/internal/core/repository"
	"github.com/getnimbus/anton/internal/core/repository/account"
	"github.com/getnimbus/anton/internal/core/repository/block"
	"github.com/getnimbus/anton/internal/core/repository/msg"
	"github.com/getnimbus/anton/internal/core/repository/tx"
)

var _ app.IndexerService = (*Service)(nil)

type Service struct {
	*app.IndexerConfig

	blockRepo   core.BlockRepository
	txRepo      core.TransactionRepository
	msgRepo     repository.Message
	accountRepo core.AccountRepository

	run bool
	mx  sync.RWMutex
	wg  sync.WaitGroup
}

func NewService(cfg *app.IndexerConfig) *Service {
	var s = new(Service)

	s.IndexerConfig = cfg

	// validate config
	if s.Workers < 1 {
		s.Workers = 1
	}
	if s.FromBlock < 2 {
		s.FromBlock = 2
	}

	//ch, pg := s.DB.CH, s.DB.PG
	pg := s.DB.PG
	kafkaProducer := s.PRODUCER
	s.txRepo = tx.NewRepository(pg, kafkaProducer)
	s.msgRepo = msg.NewRepository(pg, kafkaProducer)
	s.blockRepo = block.NewRepository(pg, kafkaProducer)
	s.accountRepo = account.NewRepository(pg)

	return s
}

func (s *Service) running() bool {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.run
}

func (s *Service) Start() error {
	ctx := context.Background()

	fromBlock := s.FromBlock

	lastMaster, err := s.blockRepo.GetLastMasterBlock(ctx)
	switch {
	case err == nil:
		fromBlock = lastMaster.SeqNo + 1
	case err != nil && !errors.Is(err, core.ErrNotFound):
		return errors.Wrap(err, "cannot get last masterchain block")
	}

	s.mx.Lock()
	s.run = true
	s.mx.Unlock()

	blocksChan := make(chan *core.Block, s.Workers*2)

	s.wg.Add(1)
	go s.fetchMasterLoop(fromBlock, blocksChan)

	s.wg.Add(1)
	go s.saveBlocksLoop(blocksChan)

	log.Info().
		Uint32("from_block", fromBlock).
		Int("workers", s.Workers).
		Msg("started")

	return nil
}

func (s *Service) Stop() {
	s.mx.Lock()
	s.run = false
	s.mx.Unlock()

	s.wg.Wait()
}
