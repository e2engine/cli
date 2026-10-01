package service

import (
	"context"
	"sync"

	"github.com/e2engine/core/pkg/keys"
	"github.com/ygrebnov/errorc"

	"github.com/e2engine/cli/pkg/errors"
)

type Provider struct {
	once sync.Once

	buildRunner bool

	app *Service
	err error
}

func NewProvider() *Provider {
	return &Provider{}
}

func NewProviderWithRunner() *Provider {
	return &Provider{
		buildRunner: true,
	}
}

func (p *Provider) Get(ctx context.Context) (*Service, error) {
	p.once.Do(func() {
		var opts []option
		if p.buildRunner {
			opts = append(opts, withRunner())
		}
		p.app, p.err = build(ctx, opts...)
	})

	if p.err != nil {
		return nil, p.err
	}

	return p.app, nil
}

func (p *Provider) Close() error {
	if p == nil {
		return nil
	}
	if p.app == nil {
		return nil
	}

	if err := p.app.Close(); err != nil {
		return errorc.With(errors.ErrFailedToCloseService, errorc.Error(keys.Cause, err))
	}

	return nil
}
