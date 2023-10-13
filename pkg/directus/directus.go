package directus

import (
	"context"
	"directus-migrator/pkg/config"
	"fmt"
	"net/http"

	"github.com/thisisdevelopment/go-dockly/xclient"
	"github.com/thisisdevelopment/go-dockly/xlogger"
)

type directus struct {
	cfg    *config.Config
	log    *xlogger.Logger
	env    string
	bearer string
	cli    xclient.IAPIClient
}

func NewDirectus(ctx context.Context, cfg *config.Config, log *xlogger.Logger, env string) (IDirectus, error) {

	o := &directus{
		cfg: cfg,
		log: log,
		env: env,
	}

	err := o.login(ctx)
	if err != nil {
		return nil, err
	}

	cli, err := o.newClient(map[string]string{"Authorization": "Bearer " + o.bearer})
	if err != nil {
		return nil, err
	}

	o.cli = cli

	return o, nil
}

func (d *directus) newClient(customHeaders map[string]string) (xclient.IAPIClient, error) {
	var cfg = xclient.GetDefaultConfig()
	if customHeaders != nil {
		cfg.CustomHeader = customHeaders
	}
	return xclient.New(d.log, d.cfg.EnvConfig[d.env].DirectusURL, nil, cfg)
}

func (d *directus) login(ctx context.Context) error {
	client, err := d.newClient(nil)
	if err != nil {
		return err
	}

	var authResponse = new(AuthResponse)
	var payload = AuthPayload{
		Email:    d.cfg.EnvConfig[d.env].DirectusUser,
		Password: d.cfg.EnvConfig[d.env].DirectusPass,
	}

	code, err := client.Do(ctx, "POST", d.cfg.DirectusConfig.AuthURL, payload, authResponse)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return fmt.Errorf("unexpected response code %d expected 200", code)
	}

	d.bearer = authResponse.Data.AccessToken

	return nil
}
