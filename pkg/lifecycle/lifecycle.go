package lifecycle

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	appErrors "github.com/xeniasokk/field-switcher/pkg/errors"
)

type App interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type Logger interface {
	Println(v ...any)
	Printf(format string, v ...any)
}

type defaultLogger struct{}

func (l *defaultLogger) Println(v ...any) {
	log.Println(v...)
}

func (l *defaultLogger) Printf(format string, v ...any) {
	log.Printf(format, v...)
}

type serviceLogger struct {
	serviceName string
	logger      Logger
}

func (l *serviceLogger) Println(v ...any) {
	prefix := fmt.Sprintf("[%s]", l.serviceName)
	l.logger.Println(append([]any{prefix}, v...)...)
}

func (l *serviceLogger) Printf(format string, v ...any) {
	prefix := fmt.Sprintf("[%s] ", l.serviceName)
	l.logger.Printf(prefix+format, v...)
}

type config struct {
	shutdownTimeout time.Duration
	logger          Logger
	serviceName     string
}

const defaultShutdownTimeout = 30 * time.Second

type Option func(*config)

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.shutdownTimeout = timeout
	}
}

func WithLogger(logger Logger) Option {
	return func(c *config) {
		c.logger = logger
	}
}

func WithServiceName(name string) Option {
	return func(c *config) {
		c.serviceName = name
	}
}

func RunWithGracefulShutdown(
	application App,
	opts ...Option,
) int {
	cfg := &config{
		shutdownTimeout: defaultShutdownTimeout,
		logger:          &defaultLogger{},
		serviceName:     "app",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	logger := cfg.logger
	if logger == nil {
		logger = &defaultLogger{}
	}

	// Обёртываем логгер с названием сервиса, если оно указано
	if cfg.serviceName != "" {
		logger = &serviceLogger{
			serviceName: cfg.serviceName,
			logger:      logger,
		}
	}

	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	ctxSig, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Println("Application starting...")

	appCtx, cancel := context.WithCancel(ctxSig)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Println("panic in Run:", r)
				errCh <- appErrors.New(appErrors.CodeInternal, fmt.Sprintf("panic in Run: %v", r))
			}
		}()
		if err := application.Run(appCtx); err != nil {
			logger.Printf("Application run error: %v (code: %s)", err, appErrors.CodeOf(err))
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	logger.Println("Application started.")

	// ждём сигнал или завершение Run
	select {
	case <-ctxSig.Done():
		logger.Println("Shutdown signal received. reason:", ctxSig.Err())
	case err := <-errCh:
		if err != nil {
			logger.Printf("Application run error: %v (code: %s)", err, appErrors.CodeOf(err))
		}
		cancel()
	}

	logger.Println("Application stopping...")

	stop()
	forceExit := make(chan os.Signal, 1)
	signal.Notify(forceExit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(forceExit)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.shutdownTimeout)
	defer shutdownCancel()

	done := make(chan error, 1)
	go func() {
		done <- application.Shutdown(shutdownCtx)
	}()

	select {
	case err := <-done:
		if err != nil {
			logger.Printf("Graceful shutdown failed: %v (code: %s)", err, appErrors.CodeOf(err))
			exitCode = 1
		} else {
			logger.Println("Graceful shutdown complete.")
		}
	case <-forceExit:
		logger.Println("Second signal received. Forcing exit.")
		exitCode = 1
	case <-shutdownCtx.Done():
		logger.Println("Shutdown timeout exceeded. Forcing exit.")
		exitCode = 1
	}

	logger.Println("Application stopped.")
	return exitCode
}
