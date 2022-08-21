package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type AttemptsOptions struct {
	Count   int
	Delay   time.Duration
	Timeout time.Duration
}

func NewClient(ctx context.Context, conf Config, opt *AttemptsOptions) (client *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)

	if opt == nil {
		opt = &AttemptsOptions{
			Count:   5,
			Delay:   time.Second * 5,
			Timeout: time.Second * 1,
		}
	}

	err = attempt(func() error {
		timeoutCtx, cancel := context.WithTimeout(ctx, opt.Timeout)
		defer cancel()

		pgxConf, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return err
		}

		client, err = pgxpool.ConnectConfig(timeoutCtx, pgxConf)
		if err != nil {
			return err
		}

		return nil
	}, *opt)

	return
}

func attempt(action func() error, opt AttemptsOptions) (err error) {
	for opt.Count > 0 {
		if err = action(); err != nil {
			time.Sleep(opt.Delay)
			opt.Count--

			continue
		}

		return
	}

	return
}
