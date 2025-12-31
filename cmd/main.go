package main

import (
	"log"
	"os"
	"time"

	"github.com/xeniasokk/field-switcher/internal/app"
	"github.com/xeniasokk/field-switcher/pkg/lifecycle"
)

const (
	shutdownTimeout = 30 * time.Second
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	a, err := app.NewApp()
	if err != nil {
		log.Printf("Failed to initialize application: %v", err)
		os.Exit(1)
	}

	exitCode := lifecycle.RunWithGracefulShutdown(
		a,
		lifecycle.WithShutdownTimeout(shutdownTimeout),
		lifecycle.WithServiceName("field-switcher"),
	)
	os.Exit(exitCode)
}
