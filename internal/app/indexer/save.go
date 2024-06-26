package indexer

import (
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/alitto/pond"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/getnimbus/anton/addr"
	"github.com/getnimbus/anton/internal/app"
	"github.com/getnimbus/anton/internal/conf"
	"github.com/getnimbus/anton/internal/core"
)

func (s *Service) insertData(
	ctx context.Context,
	acc []*core.AccountState,
	msg []*core.Message,
	tx []*core.Transaction,
	b []*core.Block,
) error {
	dbTx, err := s.DB.PG.Begin()
	if err != nil {
		return errors.Wrap(err, "cannot begin db tx")
	}
	defer func() {
		_ = dbTx.Rollback()
	}()

	for _, message := range msg {
		err := s.Parser.ParseMessagePayload(ctx, message)
		if errors.Is(err, app.ErrImpossibleParsing) {
			continue
		}
		if err != nil {
			log.Error().Err(err).
				Hex("msg_hash", message.Hash).
				Hex("src_tx_hash", message.SrcTxHash).
				Str("src_addr", message.SrcAddress.String()).
				Hex("dst_tx_hash", message.DstTxHash).
				Str("dst_addr", message.DstAddress.String()).
				Uint32("op_id", message.OperationID).
				Msg("parse message payload")
		}
	}

	eg, childCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := func() error {
			defer app.TimeTrack(time.Now(), "AddAccountStates(%d)", len(acc))
			return s.accountRepo.AddAccountStates(childCtx, dbTx, acc)
		}(); err != nil {
			return errors.Wrap(err, "add account states")
		}

		if err := func() error {
			defer app.TimeTrack(time.Now(), "AddMessages(%d)", len(msg))
			sort.Slice(msg, func(i, j int) bool { return msg[i].CreatedLT < msg[j].CreatedLT })
			return s.msgRepo.AddMessages(childCtx, dbTx, msg)
		}(); err != nil {
			return errors.Wrap(err, "add messages")
		}

		if err := func() error {
			defer app.TimeTrack(time.Now(), "AddTransactions(%d)", len(tx))
			return s.txRepo.AddTransactions(childCtx, dbTx, tx)
		}(); err != nil {
			return errors.Wrap(err, "add transactions")
		}

		if err := func() error {
			defer app.TimeTrack(time.Now(), "AddBlocks(%d)", len(b))
			return s.blockRepo.AddBlocks(childCtx, dbTx, b)
		}(); err != nil {
			return errors.Wrap(err, "add blocks")
		}

		return nil
	})

	if conf.Config.IsBackfill() {
		eg.Go(func() error {
			return s.storeS3(childCtx, b, tx, msg)
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "cannot insert data")
	}

	if err := dbTx.Commit(); err != nil {
		return errors.Wrap(err, "cannot commit db tx")
	}
	return nil
}

func (s *Service) uniqAccounts(transactions []*core.Transaction) []*core.AccountState {
	var ret []*core.AccountState

	uniqAcc := make(map[addr.Address]map[uint64]*core.AccountState)

	for _, tx := range transactions {
		if tx.Account == nil {
			continue
		}
		if uniqAcc[tx.Account.Address] == nil {
			uniqAcc[tx.Account.Address] = map[uint64]*core.AccountState{}
		}
		uniqAcc[tx.Account.Address][tx.Account.LastTxLT] = tx.Account
	}

	for _, accounts := range uniqAcc {
		for _, a := range accounts {
			if len(a.LastTxHash) > 0 {
				a.LastTxHashHex = hex.EncodeToString(a.LastTxHash)
			}
			ret = append(ret, a)
		}
	}

	return ret
}

func (s *Service) addMessage(msg *core.Message, uniqMsg map[string]*core.Message) {
	id := string(msg.Hash)

	if _, ok := uniqMsg[id]; !ok {
		uniqMsg[id] = msg
		return
	}

	switch {
	case msg.SrcTxLT != 0:
		uniqMsg[id].SrcTxLT, uniqMsg[id].SrcTxHash =
			msg.SrcTxLT, msg.SrcTxHash
		uniqMsg[id].SrcWorkchain, uniqMsg[id].SrcShard, uniqMsg[id].SrcBlockSeqNo =
			msg.SrcWorkchain, msg.SrcShard, msg.SrcBlockSeqNo
		uniqMsg[id].SrcState = msg.SrcState

	case msg.DstTxLT != 0:
		uniqMsg[id].DstTxLT, uniqMsg[id].DstTxHash =
			msg.DstTxLT, msg.DstTxHash
		uniqMsg[id].DstWorkchain, uniqMsg[id].DstShard, uniqMsg[id].DstBlockSeqNo =
			msg.DstWorkchain, msg.DstShard, msg.DstBlockSeqNo
		uniqMsg[id].DstState = msg.DstState
	}
}

func (s *Service) getMessageSource(ctx context.Context, msg *core.Message) (skip bool) {
	source, err := s.msgRepo.GetMessage(context.Background(), msg.Hash)
	if err == nil {
		msg.SrcTxLT, msg.SrcShard, msg.SrcBlockSeqNo, msg.SrcState =
			source.SrcTxLT, source.SrcShard, source.SrcBlockSeqNo, source.SrcState
		return false
	}
	if err != nil && !errors.Is(err, core.ErrNotFound) {
		//panic(errors.Wrapf(err, "get message with hash %s", msg.Hash))
		log.Info().Msg(fmt.Sprintf("get message with hash %s", msg.Hash))
		return true
	}

	// some masterchain messages does not have source
	if msg.SrcAddress.Workchain() == -1 || msg.DstAddress.Workchain() == -1 {
		return false
	}

	//blocks, err := s.blockRepo.CountMasterBlocks(ctx)
	//if err != nil {
	//	panic(errors.Wrap(err, "count masterchain blocks"))
	//}
	//if blocks < 1000 {
	//	log.Debug().
	//		Hex("dst_tx_hash", msg.DstTxHash).
	//		Int32("dst_workchain", msg.DstWorkchain).Int64("dst_shard", msg.DstShard).Uint32("dst_block_seq_no", msg.DstBlockSeqNo).
	//		Str("src_address", msg.SrcAddress.String()).Str("dst_address", msg.DstAddress.String()).
	//		Msg("cannot find source message")
	//	return true
	//}
	//
	//panic(fmt.Errorf("unknown source of message with dst tx hash %x on block (%d, %d, %d) from %s to %s",
	//	msg.DstTxHash, msg.DstWorkchain, msg.DstShard, msg.DstBlockSeqNo, msg.SrcAddress.String(), msg.DstAddress.String()))

	log.Info().Msg(fmt.Sprintf("unknown source of message with dst tx hash %x on block (%d, %d, %d) from %s to %s",
		msg.DstTxHash, msg.DstWorkchain, msg.DstShard, msg.DstBlockSeqNo, msg.SrcAddress.String(), msg.DstAddress.String()))
	return true
}

func (s *Service) uniqMessages(ctx context.Context, transactions []*core.Transaction) []*core.Message {
	var ret []*core.Message

	uniqMsg := make(map[string]*core.Message)

	for j := range transactions {
		tx := transactions[j]

		if tx.InMsg != nil {
			s.addMessage(tx.InMsg, uniqMsg)
		}
		for _, out := range tx.OutMsg {
			s.addMessage(out, uniqMsg)
		}
	}

	for _, msg := range uniqMsg {
		if msg.Type == core.Internal && (msg.SrcTxLT == 0 && msg.DstTxLT != 0) {
			if s.getMessageSource(ctx, msg) {
				continue
			}
		}

		msg = msg.WithDateKey()
		if len(msg.Hash) > 0 {
			msg.HashHex = hex.EncodeToString(msg.Hash)
		}
		if len(msg.SrcTxHash) > 0 {
			msg.SrcTxHashHex = hex.EncodeToString(msg.SrcTxHash)
		}
		if len(msg.DstTxHash) > 0 {
			msg.DstTxHashHex = hex.EncodeToString(msg.DstTxHash)
		}

		ret = append(ret, msg)
	}

	return ret
}

var lastLog = time.Now()

func (s *Service) saveBlock(ctx context.Context, master *core.Block) {
	newBlocks := append([]*core.Block{master}, master.Shards...)

	var newTransactions []*core.Transaction
	for i := range newBlocks {
		newBlocks[i] = newBlocks[i].WithDateKey()
		if len(newBlocks[i].FileHash) > 0 {
			newBlocks[i].FileHashHex = hex.EncodeToString(newBlocks[i].FileHash)
		}
		if len(newBlocks[i].RootHash) > 0 {
			newBlocks[i].RootHashHex = hex.EncodeToString(newBlocks[i].RootHash)
		}

		for _, tx := range newBlocks[i].Transactions {
			tx = tx.WithDateKey()
			if len(tx.Hash) > 0 {
				tx.HashHex = hex.EncodeToString(tx.Hash)
			}
			if len(tx.PrevTxHash) > 0 {
				tx.PrevTxHashHex = hex.EncodeToString(tx.PrevTxHash)
			}
			if len(tx.InMsgHash) > 0 {
				tx.InMsgHashHex = hex.EncodeToString(tx.InMsgHash)
			}
		}
		newTransactions = append(newTransactions, newBlocks[i].Transactions...)
	}

	if err := s.insertData(ctx, s.uniqAccounts(newTransactions), s.uniqMessages(ctx, newTransactions), newTransactions, newBlocks); err != nil {
		//panic(err)
		log.Error().Err(err).Msg("failed to insert data")
		return
	}

	lvl := log.Debug()
	if time.Since(lastLog) > 10*time.Minute {
		lvl = log.Info()
		lastLog = time.Now()
	}
	lvl.Uint32("last_inserted_seq", master.SeqNo).Msg("inserted new block")
}

func (s *Service) saveBlocksLoop(results <-chan *core.Block) {
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	for s.running() {
		var b *core.Block

		select {
		case b = <-results:
		case <-t.C:
			continue
		}

		log.Debug().
			Uint32("master_seq_no", b.SeqNo).
			Int("master_tx", len(b.Transactions)).
			Int("shards", len(b.Shards)).
			Msg("new master")

		s.saveBlock(context.Background(), b)
	}
}

func (s *Service) storeS3(
	ctx context.Context,
	blocks []*core.Block,
	txs []*core.Transaction,
	msgs []*core.Message,
) error {
	log.Info().Msg("start goroutine store s3...")

	pool := pond.New(10, 0, pond.MinWorkers(3))
	defer pool.StopAndWait()

	// store blocks to S3
	pool.Submit(func() {
		if len(blocks) == 0 {
			return
		}

		errCh := make(chan error, 1)
		defer close(errCh)

		pw := s.S3.FileStreamWriter(ctx, conf.Config.AwsBucket, fmt.Sprintf("blocks/ton-blocks/datekey=%v/%v.json.gz", blocks[0].DateKey, blocks[0].PartitionKey()), errCh)
		zw, err := gzip.NewWriterLevel(pw, gzip.BestSpeed)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("failed to create gzip writer: %v", err))
			return
		}

		// add data to gzip
		for _, block := range blocks {
			data, err := json.Marshal(block)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("failed to marshal block: %v", err))
				return
			}
			_, err = zw.Write(data)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("failed to write block to S3: %v", err))
				return
			}
			zw.Write([]byte("\n"))
		}

		// flush data to S3
		zw.Close()
		pw.Close()

		returnErr := <-errCh
		if returnErr != nil {
			log.Error().Msg(fmt.Sprintf("failed to store block to S3: %v", err))
			return
		}
		log.Info().Msg(fmt.Sprintf("[%v] submit blocks %v to S3 success", blocks[0].DateKey, blocks[0].PartitionKey()))
	})

	// store transactions to S3
	pool.Submit(func() {
		if len(txs) == 0 {
			return
		}

		errCh := make(chan error, 1)
		defer close(errCh)

		pw := s.S3.FileStreamWriter(ctx, conf.Config.AwsBucket, fmt.Sprintf("txs/ton-txs/datekey=%v/%v.json.gz", txs[0].DateKey, txs[0].BlockSeqNo), errCh)
		zw, err := gzip.NewWriterLevel(pw, gzip.BestSpeed)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("failed to create gzip writer: %v", err))
			return
		}

		// add data to gzip
		for _, tx := range txs {
			data, err := json.Marshal(tx)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("failed to marshal tx: %v", err))
				return
			}
			_, err = zw.Write(data)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("failed to write tx to S3: %v", err))
				return
			}
			zw.Write([]byte("\n"))
		}

		// flush data to S3
		zw.Close()
		pw.Close()

		returnErr := <-errCh
		if returnErr != nil {
			log.Error().Msg(fmt.Sprintf("failed to store tx to S3: %v", err))
			return
		}
		log.Info().Msg(fmt.Sprintf("[%v] submit txs in checkpoint %v to S3 success", txs[0].DateKey, txs[0].BlockSeqNo))
	})

	// store messages to S3
	pool.Submit(func() {
		if len(msgs) == 0 {
			return
		}

		errCh := make(chan error, 1)
		defer close(errCh)

		pw := s.S3.FileStreamWriter(ctx, conf.Config.AwsBucket, fmt.Sprintf("messages/ton-messages/datekey=%v/%v.json.gz", msgs[0].DateKey, msgs[0].PartitionKey()), errCh)
		zw, err := gzip.NewWriterLevel(pw, gzip.BestSpeed)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("failed to create gzip writer: %v", err))
			return
		}

		// add data to gzip
		for _, msg := range msgs {
			data, err := json.Marshal(msg)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("failed to marshal msg: %v", err))
				return
			}
			_, err = zw.Write(data)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("failed to write msg to S3: %v", err))
				return
			}
			zw.Write([]byte("\n"))
		}

		// flush data to S3
		zw.Close()
		pw.Close()

		returnErr := <-errCh
		if returnErr != nil {
			log.Error().Msg(fmt.Sprintf("failed to store msg to S3: %v", err))
			return
		}
		log.Info().Msg(fmt.Sprintf("[%v] submit msgs in checkpoint %v to S3 success", msgs[0].DateKey, msgs[0].PartitionKey()))
	})

	return nil
}
