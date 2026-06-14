package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	dictionary "github.com/localitas/localitas-app-dictionary"
	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	app := &cli.Command{
		Name:    "dictionary-server",
		Usage:   "dictionary app server",
		Version: version,
		Commands: []*cli.Command{
			serveCommand(),
		},
		DefaultCommand: "serve",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return serveAction(ctx, cmd)
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func serveCommand() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Start the server",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "listen", Value: ":0", Usage: "listen address"},
			&cli.StringFlag{Name: "base-path", Value: "/", Usage: "URL prefix for <base href>"},
			&cli.StringFlag{Name: "sources", Value: "dictionary,urban", Usage: "comma-separated lookup sources"},
		},
		Action: serveAction,
	}
}

func serveAction(ctx context.Context, cmd *cli.Command) error {
	basePath := cmd.String("base-path")

	sources := cmd.String("sources")

	a := dictionary.New(basePath, sources)
	mux := http.NewServeMux()
	a.RegisterRoutes(mux)
	mux.HandleFunc("GET /health.json", dictionary.HandleHealth)

	ln, err := net.Listen("tcp", cmd.String("listen"))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	addr := ln.Addr().(*net.TCPAddr)
	fmt.Printf("dictionary-server listening on http://localhost:%d\n", addr.Port)

	shutdown, err := dictionary.BroadcastMDNS(addr.Port, dictionary.DefaultHealth.Name)
	if err != nil {
		log.Printf("mDNS broadcast failed: %v", err)
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		if shutdown != nil {
			shutdown()
		}
		os.Exit(0)
	}()

	return http.Serve(ln, mux)
}
