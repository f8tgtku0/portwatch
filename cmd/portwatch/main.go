package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch)

func main() {
	configPath := flag.String("config", "", "path to JSON config file (optional)")
	flag.Parse()

	var cfg *config.Config
	var err error

	if *configPath != "" {
		cfg, err = config.Load(*configPath)
		if err != nil {
			log.Fatalf("portwatch: load config: %v", err)
		}
	} else {
		cfg = config.Default()
	}

	var alertWriter *os.File
	if cfg.AlertFile != "" {
		alertWriter, err = os.OpenFile(cfg.AlertFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			log.Fatalf("portwatch: open alert file: %v", err)
		}
		defer alertWriter.Close()
	} else {
		alertWriter = os.Stdout
	}

	scn := scanner.New(cfg.PortRange.Start, cfg.PortRange.End, cfg.Timeout)
	alerter := alert.New(alertWriter)
	mon := monitor.New(scn, alerter, cfg.Interval)

	fmt.Fprintf(os.Stderr, "portwatch: scanning ports %d-%d every %s\n",
		cfg.PortRange.Start, cfg.PortRange.End, cfg.Interval)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := mon.Run(ctx); err != nil {
		log.Fatalf("portwatch: %v", err)
	}
}
