package exthttp_test

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"testing"

	"github.com/BorisPlus/exthttp"
	"github.com/BorisPlus/leveledlogger"
)

func TestWeb(t *testing.T) {
	ExtHTTPServer := exthttp.NewInternalTestHTTPServer(
		"localhost",
		8099,
		leveledlogger.NewLogger(leveledlogger.DEBUG, os.Stdout),
		"./logs/",
	)
	var once sync.Once
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTSTP)
	defer once.Do(stop)
	wg := sync.WaitGroup{}
	// Stop
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		_ = ExtHTTPServer.Stop(ctx)
	}()
	// Start
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = ExtHTTPServer.Start()
	}()
	// Alive
	<-ctx.Done()
	wg.Wait()
}
