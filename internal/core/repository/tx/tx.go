package tx

import (
	"context"
	"github.com/getnimbus/anton/internal/conf"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"

	"github.com/getnimbus/anton/internal/core"
	"github.com/getnimbus/anton/internal/core/repository"
	"github.com/getnimbus/anton/internal/infra"
)

var _ repository.Transaction = (*Repository)(nil)

type Repository struct {
	//ch *ch.DB
	pg            *bun.DB
	kafkaProducer infra.KafkaSyncProducer
}

func NewRepository(
	//ck *ch.DB,
	pg *bun.DB,
	kafkaProducer infra.KafkaSyncProducer,
) *Repository {
	return &Repository{
		//ch: ck,
		pg:            pg,
		kafkaProducer: kafkaProducer,
	}
}

func createIndexes(ctx context.Context, pgDB *bun.DB) error {
	var err error

	// transactions

	_, err = pgDB.NewCreateIndex().
		Model(&core.Transaction{}).
		Using("HASH").
		Column("address").
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "transaction address pg create index")
	}

	_, err = pgDB.NewCreateIndex().
		Model(&core.Transaction{}).
		Unique().
		Column("address", "created_lt").
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "transaction account lt pg create index")
	}

	_, err = pgDB.NewCreateIndex().
		Model(&core.Transaction{}).
		Column("workchain", "shard", "block_seq_no").
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "tx block id pg create unique index")
	}

	_, err = pgDB.NewCreateIndex().
		Model(&core.Transaction{}).
		Using("BTREE").
		Column("created_lt").
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "tx created_lt pg create index")
	}

	_, err = pgDB.NewCreateIndex().
		Model(&core.Transaction{}).
		Using("HASH").
		Column("in_msg_hash").
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "tx in_msg hash pg create index")
	}

	return nil
}

func CreateTables(
	ctx context.Context,
	//chDB *ch.DB,
	pgDB *bun.DB,
) error {
	//_, err := chDB.NewCreateTable().
	//	IfNotExists().
	//	Engine("ReplacingMergeTree").
	//	Model(&core.Transaction{}).
	//	Exec(ctx)
	//if err != nil {
	//	return errors.Wrap(err, "transaction ch create table")
	//}

	_, err := pgDB.NewCreateTable().
		Model(&core.Transaction{}).
		IfNotExists().
		// ForeignKey(`(in_msg_hash) REFERENCES messages(hash)`). // TODO: fix tests with fk
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "transaction pg create table")
	}

	if err := createIndexes(ctx, pgDB); err != nil {
		return err
	}

	return nil
}

func (r *Repository) AddTransactions(ctx context.Context, tx bun.Tx, transactions []*core.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	if conf.Config.IsRealtime() {
		for _, t := range transactions {
			if err := r.kafkaProducer.SendJson(ctx, conf.Config.TonTxsTopic, t); err != nil {
				return err
			}
		}
	}

	//_, err := tx.NewInsert().Model(&transactions).Exec(ctx)
	//if err != nil {
	//	return err
	//}

	// clickhouse insert
	//_, err = r.ch.NewInsert().Model(&transactions).Exec(ctx)
	//if err != nil {
	//	return err
	//}

	return nil
}
