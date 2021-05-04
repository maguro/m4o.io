package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/golang/glog"
	"l7e.io/vanity"
	"l7e.io/vanity/cmd/vanity/server/interceptors"
	"l7e.io/vanity/pkg/toml"
	"l7e.io/yama"
)

func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	glog.Infof("port configured to listen to %s", addr)

	file := os.Getenv("TOML_FILE")
	if file == "" {
		file = "/var/vanity/config.toml"
	}
	glog.Infof("TOML file location %s", file)

	table := os.Getenv("TOML_TABLE")
	if table == "" {
		table = "vanity"
	}
	glog.Infof("TOML table name %s", table)

	be, err := toml.NewTOMLBackend(toml.InTable(table), toml.FromFile(file))
	if err != nil {
		glog.Fatal(err)
	}

	glog.Info("Vanity entries:")
	err = be.List(context.Background(), vanity.ConsumerFunc(func(context context.Context, importPath, vcs, vcsPath string) {
		glog.Infof("  %s %s %s", importPath, vcs, vcsPath)
	}))
	if err != nil {
		glog.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", interceptors.WrapHandler(vanity.NewVanityHandler(be)))

	s := &http.Server{Addr: addr, Handler: mux}

	watcher := yama.NewWatcher(
		yama.WatchingSignals(syscall.SIGINT, syscall.SIGTERM),
		yama.WithTimeout(2*time.Second), // nolint
		yama.WithClosers(s))

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		glog.Error(err)
		_ = watcher.Close()
	}

	if err := watcher.Wait(); err != nil {
		glog.Warningf("Shutdown error: %s", err)
	}
}
