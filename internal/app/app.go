package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/aintsashqa/simple-todo/internal"
	"github.com/aintsashqa/simple-todo/internal/config"
	"github.com/aintsashqa/simple-todo/internal/delivery/http"
	"github.com/aintsashqa/simple-todo/internal/repository"
	"github.com/aintsashqa/simple-todo/pkg/database/postgresql"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func Run() {
	ctx := context.Background()
	conf := config.GetConfig()
	var l log.Logger
	{
		l = log.NewJSONLogger(os.Stdout)
		l = log.NewSyncLogger(l)
		l = log.With(l,
			"service", "simple-todo",
			"time", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(l).Log("msg", "service started")
	defer level.Info(l).Log("msg", "service stopped")

	var r internal.Repository
	{
		level.Info(l).Log("msg", "init database connection")
		c, err := postgresql.NewClient(ctx, postgresql.Config{
			Host:     conf.Database.Host,
			Port:     conf.Database.Port,
			Username: conf.Database.Username,
			Password: conf.Database.Password,
			Database: conf.Database.Name,
		}, nil)

		if err != nil {
			level.Error(l).Log("err", err)
			os.Exit(-1)
		}

		r = repository.NewRepository(c, l)
	}

	errs := make(chan error)

	{
		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt, os.Kill)
			errs <- fmt.Errorf("%s", <-sig)
		}()

		level.Info(l).Log("msg", "init http server")
		h := http.NewHandler(r, l)
		srv := http.NewServer(conf.Server, h)

		go func() {
			level.Info(l).Log("msg", "http server running")
			errs <- srv.ListenAndServe()
		}()
	}

	level.Error(l).Log("exit", <-errs)
}
