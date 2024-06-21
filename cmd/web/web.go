package web

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
	"github.com/getnimbus/anton/internal/api/http"
	"github.com/getnimbus/anton/internal/app"
	"github.com/getnimbus/anton/internal/app/query"
	"github.com/getnimbus/anton/internal/core/repository"
	"github.com/getnimbus/anton/internal/core/repository/contract"
)

var Command = &cli.Command{
	Name:  "web",
	Usage: "HTTP JSON API",

	Action: func(ctx *cli.Context) error {
		//chURL := env.GetString("DB_CH_URL", "")
		//pgURL := env.GetString("DB_PG_URL", "")
		pgURL := conf.Config.DbPgUrl

		conn, err := repository.ConnectDB(
			ctx.Context,
			//chURL,
			pgURL,
		)
		if err != nil {
			return errors.Wrap(err, "cannot connect to a database")
		}

		def, err := contract.NewRepository(conn.PG).GetDefinitions(ctx.Context)
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

		qs, err := query.NewService(ctx.Context, &app.QueryConfig{
			DB:  conn,
			API: api,
		})
		if err != nil {
			return err
		}

		srv := http.NewServer(
			//env.GetString("LISTEN", "0.0.0.0:80"),
			"0.0.0.0:8080",
		)
		srv.RegisterRoutes(http.NewController(qs))

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-c
			conn.Close()
			os.Exit(0)
		}()

		if err = srv.Run(); err != nil {
			return err
		}

		return nil
	},
}
