package rescan

import (
	"fmt"
	"github.com/getnimbus/anton/internal/conf"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"

	"github.com/getnimbus/anton/abi"
	"github.com/getnimbus/anton/internal/app"
	"github.com/getnimbus/anton/internal/app/parser"
	"github.com/getnimbus/anton/internal/app/rescan"
	"github.com/getnimbus/anton/internal/core/repository"
	"github.com/getnimbus/anton/internal/core/repository/account"
	"github.com/getnimbus/anton/internal/core/repository/block"
	"github.com/getnimbus/anton/internal/core/repository/contract"
	"github.com/getnimbus/anton/internal/core/repository/msg"
	rescanRepository "github.com/getnimbus/anton/internal/core/repository/rescan"
	"github.com/getnimbus/anton/internal/infra"
)

var Command = &cli.Command{
	Name: "rescan",

	Usage: "Updates account states and messages data",

	Action: func(ctx *cli.Context) error {
		//chURL := env.GetString("DB_CH_URL", "")
		//pgURL := env.GetString("DB_PG_URL", "")
		//chURL := conf.Config.DbChUrl
		pgURL := conf.Config.DbPgUrl

		conn, err := repository.ConnectDB(
			ctx.Context,
			//chURL,
			pgURL,
		)
		if err != nil {
			return errors.Wrap(err, "cannot connect to a database")
		}

		kafkaProducer, cleanup, err := infra.NewKafkaSyncProducer(ctx.Context)
		if err != nil {
			return err
		}

		contractRepo := contract.NewRepository(conn.PG)

		interfaces, err := contractRepo.GetInterfaces(ctx.Context)
		if err != nil {
			return errors.Wrap(err, "get interfaces")
		}
		if len(interfaces) == 0 {
			return errors.New("no contract interfaces")
		}

		def, err := contractRepo.GetDefinitions(ctx.Context)
		if err != nil {
			return errors.Wrap(err, "get definitions")
		}
		err = abi.RegisterDefinitions(def)
		if err != nil {
			return errors.Wrap(err, "get definitions")
		}

		client := liteclient.NewConnectionPool()
		api := ton.NewAPIClient(client, ton.ProofCheckPolicyUnsafe).WithRetry()
		for _, addr := range strings.Split(conf.Config.LiteServers, ",") {
			split := strings.Split(addr, "|")
			if len(split) != 2 {
				return fmt.Errorf("wrong server address format '%s'", addr)
			}
			host, key := split[0], split[1]
			if err := client.AddConnection(ctx.Context, host, key); err != nil {
				return errors.Wrapf(err, "cannot add connection with %s host and %s key", host, key)
			}
		}
		bcConfig, err := app.GetBlockchainConfig(ctx.Context, api)
		if err != nil {
			return errors.Wrap(err, "cannot get blockchain config")
		}

		p := parser.NewService(&app.ParserConfig{
			BlockchainConfig: bcConfig,
			ContractRepo:     contractRepo,
		})
		i := rescan.NewService(&app.RescanConfig{
			ContractRepo: contractRepo,
			RescanRepo:   rescanRepository.NewRepository(conn.PG),
			AccountRepo: account.NewRepository(
				//conn.CH,
				conn.PG,
			),
			BlockRepo: block.NewRepository(
				//conn.CH,
				conn.PG,
				kafkaProducer,
			),
			MessageRepo: msg.NewRepository(
				//conn.CH,
				conn.PG,
				kafkaProducer,
			),
			Parser:      p,
			Workers:     conf.Config.RescanWorkers,
			SelectLimit: conf.Config.RescanSelectLimit,
		})
		if err = i.Start(); err != nil {
			return err
		}

		c := make(chan os.Signal, 1)
		done := make(chan struct{}, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-c
			i.Stop()
			conn.Close()
			cleanup()
			done <- struct{}{}
		}()

		<-done

		return nil
	},
}
