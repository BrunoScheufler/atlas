package sdk

import (
	"context"
	"flag"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func Start(atlasfile *atlasfile.Atlasfile) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fallbackLevel := "info"
	if os.Getenv("LOG_LEVEL") != "" {
		fallbackLevel = os.Getenv("LOG_LEVEL")
	}

	fallbackPort := 0
	parsedPort, _ := strconv.Atoi(os.Getenv("PORT"))
	if parsedPort != 0 {
		fallbackPort = parsedPort
	}

	port := flag.Int("port", fallbackPort, "port to serve on")
	logLevel := flag.String("loglevel", fallbackLevel, "log level")
	flag.Parse()

	parsedLevel, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	logger := logrus.New()
	logger.SetLevel(parsedLevel)

	if *port == 0 {
		logger.Fatal("port flag must be provided with non-zero value")
	}

	baseLogger := logger.WithContext(ctx).WithField("service", "atlasfile")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		logger.WithFields(logrus.Fields{
			"signal": sig,
		}).Traceln("received signal, shutting down")
		cancel()
	}()

	baseLogger.WithField("logLevel", logLevel).Traceln("starting atlasfile provider")

	err = serve(ctx, baseLogger, atlasfile, *port)
	if err != nil {
		return err
	}

	return nil
}
