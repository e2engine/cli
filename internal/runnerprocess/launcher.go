package runnerprocess

import (
	"context"
	"net"
	"os"
	"os/exec"
	"time"

	corekeys "github.com/e2engine/core/pkg/keys"
	"github.com/e2engine/core/pkg/log"

	"github.com/e2engine/cli/pkg/keys"
)

type Launcher struct {
	address string
	logger  log.Logger
}

func NewLauncher(address string, logger log.Logger) *Launcher {
	return &Launcher{
		address: address,
		logger:  logger.With(log.String(corekeys.Component, "runnerprocess-launcher")),
	}
}

func (l *Launcher) Start(ctx context.Context) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(
		executable,
		"internal-runner",
	)

	if err := cmd.Start(); err != nil {
		return err
	}

	pid := cmd.Process.Pid

	l.logger.Debug(
		"Started internal-runner",
		log.Int(keys.PID, pid),
	)

	go func() {
		err := cmd.Wait()

		l.logger.Debug(
			"Internal-runner exited",
			log.Int(keys.PID, pid),
			log.Err("err", err),
		)
	}()

	if err := l.waitReady(ctx); err != nil {
		_ = cmd.Process.Kill()
		return err
	}

	return nil
}

// TODO: Readiness currently verifies only that the configured address accepts
// TCP connections. It does not verify that the newly started runner owns the
// listener. Consider a runner-specific readiness handshake if this distinction
// becomes important.
func (l *Launcher) waitReady(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		conn, err := (&net.Dialer{}).DialContext(
			ctx,
			"tcp",
			l.address,
		)
		if err == nil {
			return conn.Close()
		}

		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
		}
	}
}
