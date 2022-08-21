package http

import (
	"fmt"
	"net/http"

	"github.com/aintsashqa/simple-todo/internal/config"
)

func NewServer(conf config.Server, h http.Handler) http.Server {
	return http.Server{
		Addr:           fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler:        h,
		ReadTimeout:    conf.ReadTimeout,
		WriteTimeout:   conf.WriteTimeout,
		MaxHeaderBytes: int(conf.MaxHeaderBytes) << 20,
	}
}
