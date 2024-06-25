package block

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"

	"github.com/getnimbus/anton/internal/conf"
	"github.com/getnimbus/anton/internal/core"
	"github.com/getnimbus/anton/internal/core/repository"
	"github.com/getnimbus/anton/internal/infra"
)

var _ repository.Block = (*Repository)(nil)

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
	_, err := pgDB.NewCreateIndex().
		Model(&core.Block{}).
		Using("BTREE").
		Column("workchain", "seq_no").
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "block workchain pg create index")
	}

	_, err = pgDB.NewCreateIndex().
		Model(&core.Block{}).
		Using("BTREE").
		Column("master_workchain", "master_shard", "master_seq_no").
		Where("workchain != -1").
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "block workchain pg create index")
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
	//	Model(&core.Block{}).
	//	Exec(ctx)
	//if err != nil {
	//	return errors.Wrap(err, "block ch create table")
	//}

	_, err := pgDB.NewCreateTable().
		Model(&core.Block{}).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "block pg create table")
	}

	return createIndexes(ctx, pgDB)
}

func (r *Repository) AddBlocks(ctx context.Context, tx bun.Tx, info []*core.Block) error {
	if len(info) == 0 {
		return nil
	}
	for _, b := range info {
		if conf.Config.IsBackfill() {
			b.SyncType = core.BlockSyncType_BACKFILL
		}
		_, err := tx.NewInsert().Model(b).Exec(ctx)
		if err != nil {
			return err
		}
		if conf.Config.IsRealtime() {
			if err = r.kafkaProducer.SendJson(ctx, conf.Config.TonBlocksTopic, b); err != nil {
				return err
			}
		}
	}

	// send kafka

	// clickhouse insert
	//_, err := r.ch.NewInsert().Model(&info).Exec(ctx)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (r *Repository) GetLastMasterBlock(ctx context.Context) (*core.Block, error) {
	ret := new(core.Block)

	q := r.pg.NewSelect().Model(ret).
		Where("workchain = ?", -1).
		Order("seq_no DESC").
		Limit(1)
	if conf.Config.IsBackfill() {
		q = q.Where("sync_type = ?", core.BlockSyncType_BACKFILL)
	} else {
		q = q.Where("sync_type = ?", core.BlockSyncType_REALTIME)
	}

	err := q.Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Repository) CountMasterBlocks(ctx context.Context) (int, error) {
	//ret, err := r.ch.NewSelect().Model((*core.Block)(nil)).
	//	Where("workchain = -1").
	//	Count(ctx)
	//if err != nil {
	//	return 0, err
	//}
	//return ret, nil
	return 0, nil
}
