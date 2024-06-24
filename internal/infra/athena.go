package infra

import (
	"database/sql"

	"github.com/getnimbus/anton/internal/conf"

	drv "github.com/uber/athenadriver/go"
)

func NewAthenaSession(outputBucket string) (*sql.DB, error) {
	conf, _ := drv.NewDefaultConfig(
		outputBucket,
		conf.Config.AwsRegion,
		conf.Config.AwsAccessKeyId,
		conf.Config.AwsSecretAccessKey,
	)

	dsn := conf.Stringify()
	return sql.Open(drv.DefaultDBName, dsn)
}
