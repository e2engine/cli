package errors

import (
	"github.com/ygrebnov/errorc"
)

var (
	ErrCannotLoadConfig          = errorc.New("cannot load config")
	ErrCannotResolveUserHomeDir  = errorc.New("cannot resolve user home dir")
	ErrCannotInitializeLogger    = errorc.New("cannot initialize logger")
	ErrCannotConfigureLogger     = errorc.New("cannot configure logger")
	ErrCannotInitializeService   = errorc.New("cannot initialize service")
	ErrFailedToCloseService      = errorc.New("failed to close service")
	ErrCannotRenderOutput        = errorc.New("cannot render output")
	ErrFailedToStartRunner       = errorc.New("failed to start runner")
	ErrExecutionNotPassed        = errorc.New("execution did not pass")
	ErrUnexpectedExecutionStatus = errorc.New("unexpected execution status")
	ErrExecutionCheckTimeout     = errorc.New("execution check timed out")
)
