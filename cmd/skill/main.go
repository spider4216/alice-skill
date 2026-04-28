package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/spider4216/alice-skill/internal/logger"
	"github.com/spider4216/alice-skill/internal/store/pg"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	ApiVer = "1.0"
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	conn, err := sql.Open("pgx", flagDatabaseURI)
	if err != nil {
		return err
	}

	store := pg.NewStore(conn)
	store.Bootstrap(context.Background())

	appInstance := newApp(store)

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))

	fmt.Println("Run on", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, http.HandlerFunc(logger.RequestLogger(gzipMiddleware(appInstance.webhook))))
}

func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")

		if supportGzip {
			cw := NewCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			cr, err := NewCompressReader(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	}
}
