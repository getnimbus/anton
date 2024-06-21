package main

import (
	"fmt"
	"github.com/getnimbus/anton/internal/conf"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/getnimbus/anton/cmd/archive"
	"github.com/getnimbus/anton/cmd/contract"
	"github.com/getnimbus/anton/cmd/db"
	"github.com/getnimbus/anton/cmd/indexer"
	"github.com/getnimbus/anton/cmd/label"
	"github.com/getnimbus/anton/cmd/rescan"
	"github.com/getnimbus/anton/cmd/web"
)

func init() {
	// load env
	if err := conf.LoadConfig("."); err != nil {
		panic(fmt.Errorf("cannot load config: %v", err))
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	level := zerolog.InfoLevel
	if conf.Config.IsDebug() {
		level = zerolog.DebugLevel
	}

	// add file and line number to log
	log.Logger = log.With().Caller().Logger().Level(level)
}

func main() {
	app := &cli.App{
		Name:  "anton",
		Usage: "an indexing project",
		Commands: []*cli.Command{
			db.Command,
			indexer.Command,
			web.Command,
			archive.Command,
			contract.Command,
			label.Command,
			rescan.Command,
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
