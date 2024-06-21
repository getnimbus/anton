package db

import (
	"context"
	"github.com/getnimbus/anton/internal/conf"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"

	"github.com/getnimbus/anton/internal/core/repository"
	"github.com/getnimbus/anton/migrations/pgmigrations"
)

func newMigrators() (
	pg *migrate.Migrator,
	//ch *chmigrate.Migrator,
	err error,
) {
	//chURL := env.GetString("DB_CH_URL", "")
	//pgURL := env.GetString("DB_PG_URL", "")
	//chURL := conf.Config.DbChUrl
	pgURL := conf.Config.DbPgUrl

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := repository.ConnectDB(ctx,
		//chURL,
		pgURL,
	)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to the databases")
	}

	mpg := migrate.NewMigrator(conn.PG, pgmigrations.Migrations)
	//mch := chmigrate.NewMigrator(conn.CH, chmigrations.Migrations)
	//
	//return mpg, mch, nil
	return mpg, nil
}

func pgUnlock(ctx context.Context, m *migrate.Migrator) {
	if err := m.Unlock(ctx); err != nil {
		log.Error().Err(err).Msg("cannot unlock pg")
	}
}

//func chUnlock(ctx context.Context, m *chmigrate.Migrator) {
//	if err := m.Unlock(ctx); err != nil {
//		log.Error().Err(err).Msg("cannot unlock ch")
//	}
//}

var Command = &cli.Command{
	Name:  "migrate",
	Usage: "Migrates database",

	Subcommands: []*cli.Command{
		{
			Name:  "init",
			Usage: "Creates migration tables",
			Action: func(c *cli.Context) error {
				//mpg, mch, err := newMigrators()
				mpg, err := newMigrators()
				if err != nil {
					return err
				}
				if err := mpg.Init(c.Context); err != nil {
					return err
				}
				//if err := mch.Init(c.Context); err != nil {
				//	return err
				//}
				return nil
			},
		},
		{
			Name:      "create",
			Usage:     "Creates up and down SQL migrations",
			ArgsUsage: "migration_name",
			Action: func(c *cli.Context) error {
				name := strings.Join(c.Args().Slice(), "_")

				mpg := migrate.NewMigrator(nil, pgmigrations.Migrations)
				//mch := chmigrate.NewMigrator(nil, chmigrations.Migrations)

				pgFiles, err := mpg.CreateSQLMigrations(c.Context, name)
				if err != nil {
					return err
				}
				for _, mf := range pgFiles {
					log.Info().Str("name", mf.Name).Str("path", mf.Path).Msg("created pg migration")
				}

				//chFiles, err := mch.CreateSQLMigrations(c.Context, name)
				//if err != nil {
				//	return err
				//}
				//for _, mf := range chFiles {
				//	log.Info().Str("name", mf.Name).Str("path", mf.Path).Msg("created ch migration")
				//}

				return nil
			},
		},
		{
			Name:  "up",
			Usage: "Migrates database",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "init",
					Aliases: []string{"i"},
					Value:   false,
					Usage:   "create migration tables",
				},
			},
			Action: func(c *cli.Context) error {
				//mpg, mch, err := newMigrators()
				mpg, err := newMigrators()
				if err != nil {
					return err
				}

				if c.Bool("init") {
					if err := mpg.Init(c.Context); err != nil {
						return err
					}
					//if err := mch.Init(c.Context); err != nil {
					//	return err
					//}
				}

				// postgresql

				if err := mpg.Lock(c.Context); err != nil {
					return err
				}
				defer pgUnlock(c.Context, mpg)

				pgGroup, err := mpg.Migrate(c.Context)
				if err != nil {
					return err
				}
				if pgGroup.IsZero() {
					log.Info().Msg("there are no new migrations to run (pg database is up to date)")
				} else {
					log.Info().Str("group", pgGroup.String()).Msg("pg migrated")
				}

				// clickhouse
				//if err := mch.Lock(c.Context); err != nil {
				//	return err
				//}
				//defer chUnlock(c.Context, mch)
				//
				//chGroup, err := mch.Migrate(c.Context)
				//if err != nil {
				//	return err
				//}
				//if chGroup.IsZero() {
				//	log.Info().Msg("there are no new migrations to run (ch database is up to date)")
				//	return nil
				//} else {
				//	log.Info().Str("group", chGroup.String()).Msg("ch migrated")
				//}

				return nil
			},
		},
		{
			Name:  "down",
			Usage: "Rollbacks the last migration group",
			Action: func(c *cli.Context) error {
				//mpg, mch, err := newMigrators()
				mpg, err := newMigrators()
				if err != nil {
					return err
				}

				// postgresql

				if err := mpg.Lock(c.Context); err != nil {
					return err
				}
				defer pgUnlock(c.Context, mpg)

				pgGroup, err := mpg.Rollback(c.Context)
				if err != nil {
					return err
				}
				if pgGroup.IsZero() {
					log.Info().Msg("there are no pg groups to roll back")
				} else {
					log.Info().Str("group", pgGroup.String()).Msg("pg rolled back")
				}

				// clickhouse
				//if err := mch.Lock(c.Context); err != nil {
				//	return err
				//}
				//defer chUnlock(c.Context, mch)
				//
				//chGroup, err := mch.Rollback(c.Context)
				//if err != nil {
				//	return err
				//}
				//if chGroup.IsZero() {
				//	log.Info().Msg("there are no ch groups to roll back")
				//} else {
				//	log.Info().Str("group", chGroup.String()).Msg("ch rolled back")
				//}

				return nil
			},
		},
		{
			Name:  "status",
			Usage: "Prints migrations status",
			Action: func(c *cli.Context) error {
				//mpg, mch, err := newMigrators()
				mpg, err := newMigrators()
				if err != nil {
					return err
				}

				spg, err := mpg.MigrationsWithStatus(c.Context)
				if err != nil {
					return err
				}
				log.Info().Str("slice", spg.String()).Msg("pg all")
				log.Info().Str("slice", spg.Unapplied().String()).Msg("pg unapplied")
				log.Info().Str("group", spg.LastGroup().String()).Msg("pg last migration")

				//sch, err := mch.MigrationsWithStatus(c.Context)
				//if err != nil {
				//	return err
				//}
				//log.Info().Str("slice", sch.String()).Msg("ch all")
				//log.Info().Str("slice", sch.Unapplied().String()).Msg("ch unapplied")
				//log.Info().Str("group", sch.LastGroup().String()).Msg("ch last migration")

				return nil
			},
		},
		{
			Name:  "lock",
			Usage: "Locks migrations",
			Action: func(c *cli.Context) error {
				//mpg, mch, err := newMigrators()
				mpg, err := newMigrators()
				if err != nil {
					return err
				}
				if err := mpg.Lock(c.Context); err != nil {
					return err
				}
				//if err := mch.Lock(c.Context); err != nil {
				//	return err
				//}
				return nil
			},
		},
		{
			Name:  "unlock",
			Usage: "Unlocks migrations",
			Action: func(c *cli.Context) error {
				//mpg, mch, err := newMigrators()
				mpg, err := newMigrators()
				if err != nil {
					return err
				}
				if err := mpg.Unlock(c.Context); err != nil {
					return err
				}
				//if err := mch.Unlock(c.Context); err != nil {
				//	return err
				//}
				return nil
			},
		},
		{
			Name:  "markApplied",
			Usage: "Marks migrations as applied without actually running them",
			Action: func(c *cli.Context) error {
				//mpg, mch, err := newMigrators()
				mpg, err := newMigrators()
				if err != nil {
					return err
				}
				pgGroup, err := mpg.Migrate(c.Context, migrate.WithNopMigration())
				if err != nil {
					return err
				}
				if pgGroup.IsZero() {
					log.Info().Msg("there are no new pg migrations to mark as applied")
					return nil
				}
				log.Info().Str("group", pgGroup.String()).Msg("pg marked as applied")

				//chGroup, err := mch.Migrate(c.Context, chmigrate.WithNopMigration())
				//if err != nil {
				//	return err
				//}
				//if chGroup.IsZero() {
				//	log.Info().Msg("there are no new ch migrations to mark as applied\n")
				//	return nil
				//}
				log.Info().Str("group", pgGroup.String()).Msg("ch marked as applied")

				return nil
			},
		},
		//{
		//	Name:  "transferCodeData",
		//	Usage: "Moves account states code and data from ClickHouse to RocksDB",
		//	Flags: []cli.Flag{
		//		&cli.IntFlag{
		//			Name:  "limit",
		//			Value: 10000,
		//			Usage: "batch size for update",
		//		},
		//		&cli.Uint64Flag{
		//			Name:  "start-from",
		//			Value: 0,
		//			Usage: "last tx lt to start from",
		//		},
		//	},
		//	Action: func(c *cli.Context) error {
		//		chURL := env.GetString("DB_CH_URL", "")
		//		pgURL := env.GetString("DB_PG_URL", "")
		//
		//		conn, err := repository.ConnectDB(c.Context, chURL, pgURL)
		//		if err != nil {
		//			return errors.Wrap(err, "cannot connect to the databases")
		//		}
		//		defer conn.Close()
		//
		//		lastTxLT := c.Uint64("start-from")
		//		limit := c.Int("limit")
		//		for {
		//			var nextTxLT uint64
		//
		//			err := conn.CH.NewRaw(`
		//					SELECT last_tx_lt
		//					FROM account_states
		//					WHERE last_tx_lt >= ?
		//					ORDER BY last_tx_lt ASC
		//					OFFSET ? ROW FETCH FIRST 1 ROWS ONLY`, lastTxLT, limit).
		//				Scan(c.Context, &nextTxLT)
		//			if errors.Is(err, sql.ErrNoRows) {
		//				nextTxLT = 0
		//				log.Info().Uint64("last_tx_lt", lastTxLT).Msg("finishing")
		//			} else if err != nil {
		//				return errors.Wrap(err, "get next tx lt")
		//			}
		//
		//			f := fmt.Sprintf("last_tx_lt >= %d", lastTxLT)
		//			if nextTxLT != 0 {
		//				f += fmt.Sprintf(" AND last_tx_lt < %d", nextTxLT)
		//			}
		//
		//			err = conn.CH.NewRaw(`
		//					INSERT INTO account_states_code
		//					SELECT code_hash, any(code)
		//					FROM (
		//						SELECT code_hash, code
		//						FROM account_states
		//						WHERE ` + f + ` AND length(code) > 0
		//					)
		//					GROUP BY code_hash`).
		//				Scan(c.Context)
		//			if err != nil {
		//				return errors.Wrapf(err, "transfer code from %d to %d", lastTxLT, nextTxLT)
		//			}
		//
		//			err = conn.CH.NewRaw(`
		//					INSERT INTO account_states_data
		//					SELECT data_hash, any(data)
		//					FROM (
		//						SELECT data_hash, data
		//						FROM account_states
		//						WHERE ` + f + ` AND length(data) > 0
		//					)
		//					GROUP BY data_hash`).
		//				Scan(c.Context)
		//			if err != nil {
		//				return errors.Wrapf(err, "transfer data from %d to %d", lastTxLT, nextTxLT)
		//			}
		//
		//			if nextTxLT == 0 {
		//				log.Info().Msg("finished")
		//				return nil
		//			} else {
		//				log.Info().Uint64("last_tx_lt", lastTxLT).Uint64("next_tx_lt", nextTxLT).Msg("transferred new batch")
		//			}
		//			lastTxLT = nextTxLT
		//		}
		//	},
		//},
		//{
		//	Name:  "clearCodeData",
		//	Usage: "Removes account states code and data from PostgreSQL",
		//	Flags: []cli.Flag{
		//		&cli.IntFlag{
		//			Name:  "limit",
		//			Value: 10000,
		//			Usage: "batch size for update",
		//		},
		//		&cli.Uint64Flag{
		//			Name:  "start-from",
		//			Value: 0,
		//			Usage: "last tx lt to start from",
		//		},
		//	},
		//	Action: func(c *cli.Context) error {
		//		chURL := env.GetString("DB_CH_URL", "")
		//		pgURL := env.GetString("DB_PG_URL", "")
		//
		//		conn, err := repository.ConnectDB(c.Context, chURL, pgURL)
		//		if err != nil {
		//			return errors.Wrap(err, "cannot connect to the databases")
		//		}
		//		defer conn.Close()
		//
		//		lastTxLT := c.Uint64("start-from")
		//		limit := c.Int("limit")
		//		for {
		//			var nextTxLT uint64
		//
		//			err := conn.CH.NewRaw(`
		//					SELECT last_tx_lt
		//					FROM account_states
		//					WHERE last_tx_lt >= ?
		//					ORDER BY last_tx_lt ASC
		//					OFFSET ? ROW FETCH FIRST 1 ROWS ONLY`, lastTxLT, limit).
		//				Scan(c.Context, &nextTxLT)
		//			if errors.Is(err, sql.ErrNoRows) {
		//				nextTxLT = 0
		//				log.Info().Uint64("last_tx_lt", lastTxLT).Msg("finishing")
		//			} else if err != nil {
		//				return errors.Wrap(err, "get next tx lt")
		//			}
		//
		//			q := conn.PG.NewUpdate().Model((*core.AccountState)(nil)).
		//				Set("code = NULL").
		//				Where("last_tx_lt >= ?", lastTxLT)
		//			if nextTxLT != 0 {
		//				q = q.Where("last_tx_lt < ?", nextTxLT)
		//			}
		//			if _, err := q.Exec(c.Context); err != nil {
		//				return errors.Wrapf(err, "clear code from %d to %d", lastTxLT, nextTxLT)
		//			}
		//
		//			q = conn.PG.NewUpdate().Model((*core.AccountState)(nil)).
		//				Set("data = NULL").
		//				Where("last_tx_lt >= ?", lastTxLT)
		//			if nextTxLT != 0 {
		//				q = q.Where("last_tx_lt < ?", nextTxLT)
		//			}
		//			if _, err := q.Exec(c.Context); err != nil {
		//				return errors.Wrapf(err, "clear data from %d to %d", lastTxLT, nextTxLT)
		//			}
		//
		//			log.Info().Uint64("last_tx_lt", lastTxLT).Uint64("next_tx_lt", nextTxLT).Msg("cleared new batch")
		//
		//			if nextTxLT == 0 {
		//				log.Info().Msg("finished")
		//				return nil
		//			}
		//			lastTxLT = nextTxLT
		//		}
		//	},
		//},
	},
}
