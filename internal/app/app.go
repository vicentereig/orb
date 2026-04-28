package app

import (
	"context"

	"github.com/vicentereig/orb/internal/output"
)

type OrbClient interface {
	Ping(ctx context.Context) (interface{}, error)
}

type App struct {
	client  OrbClient
	version string
}

func New(client OrbClient, version string) *App {
	return &App{client: client, version: version}
}

func (a *App) Version() string {
	return output.Success(
		map[string]string{"version": a.version},
		map[string]string{"resource": "system", "operation": "version"},
	)
}

func (a *App) Ping(ctx context.Context) string {
	res, err := a.client.Ping(ctx)
	if err != nil {
		return output.Error(err, "api_error", map[string]string{"resource": "system", "operation": "ping"})
	}
	return output.Success(res, map[string]string{"resource": "system", "operation": "ping"})
}
