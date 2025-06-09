package main

//go:generate make api.json
//go:generate go tool oapi-codegen -config cfg.yaml api.json

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"github.com/vitaminmoo/poe-schema-claims/api"
	"github.com/vitaminmoo/poe-schema-claims/ctxutil"
	"github.com/vitaminmoo/poe-schema-claims/log"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

//go:embed api.json
var specFile []byte

//go:embed static
var static embed.FS

func main() {
	z := log.Init("DEBUG")
	ctx := log.Save(context.Background(), z)
	defer z.Sync()
	Run(ctx)
}

func Run(ctx context.Context) error {
	z := log.Load(ctx)
	g, ctx := errgroup.WithContext(ctx)
	db, err := sql.Open("sqlite3", "file:demo.db")
	if err != nil {
		z.Fatal("opening database", zap.Error(err))
	}
	defer db.Close()

	if err := dbSetup(ctx, db); err != nil {
		z.Fatal("setting up database", zap.Error(err))
	}

	if err := startHttpServer(ctx, db, g); err != nil {
		return fmt.Errorf("starting http server: %w", err)
	}

	g.Go(func() error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		z.Info("shutting down due to signal")
		return context.Canceled
	})

	err = g.Wait()
	if err != nil {
		if err == context.Canceled {
			return nil
		}
		return err
	}
	return nil
}

func startHttpServer(ctx context.Context, db *sql.DB, g *errgroup.Group) error {
	z := log.Load(ctx)
	spec, err := openapi3.NewLoader().LoadFromData(specFile)
	if err != nil {
		z.Fatal("loading openapi spec", zap.Error(err))
	}

	server := api.NewServer(db)

	rootMux := http.NewServeMux()
	rootMux.Handle("GET /healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"healthy":true}`))
	}))
	rootMux.Handle("GET /static/", http.FileServerFS(static))
	rootMux.Handle("GET /openapi.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specFile))
	}))

	rootMux.Handle("/v1/", api.HandlerWithOptions(
		api.NewStrictHandlerWithOptions(
			server,
			[]api.StrictMiddlewareFunc{
				strictMiddlewareExample,
			},
			api.StrictHTTPServerOptions{
				RequestErrorHandlerFunc:  strictRequestErrorHandler,
				ResponseErrorHandlerFunc: strictResponseErrorHandler,
			},
		),
		api.StdHTTPServerOptions{
			BaseURL: "/v1",
			Middlewares: []api.MiddlewareFunc{
				nethttpmiddleware.OapiRequestValidatorWithOptions(spec,
					&nethttpmiddleware.Options{
						Options:               openapi3filter.Options{},
						SilenceServersWarning: true,
						ErrorHandler:          middlewareErrorHandler,
					},
				),
			},
			ErrorHandlerFunc: errorHandler,
		},
	))

	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: func() []string {
			return nil
		}(),
		AllowOriginFunc: func(origin string) bool {
			//return slices.Contains([]string{"https://poe-schema.obsoleet.org"}, origin)
			return false
		},
		AllowCredentials:    true,
		AllowPrivateNetwork: false,
		OptionsPassthrough:  false,
	})

	// ctxutil MUST be first, logging SHOULD be second
	n := negroni.New(ctxutil.NewMiddleware(), log.New(), c)
	n.UseHandler(rootMux)

	s := &http.Server{
		Handler:           n,
		ReadHeaderTimeout: 1 * time.Second,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      5 * time.Second,
		Addr:              "0.0.0.0:8080",
	}

	g.Go(func() error {
		z.Info("starting the http server")
		return s.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done()
		z.Info("shutting down the http server")
		return s.Shutdown(ctx)
	})

	return nil
}

func strictMiddlewareExample(next strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (response any, err error) {
		return next(ctx, w, r, request)
	}
}

func dbSetup(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		DROP TABLE IF EXISTS enum;
		CREATE TABLE enum(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source VARCHAR(255) NOT NULL,
			name VARCHAR(128) NOT NULL,
			client_labels TEXT NOT NULL,
			server_labels TEXT NOT NULL,
			vals TEXT NOT NULL,
			zero_indexed BOOLEAN NOT NULL
		);
	`)
	return err
}
