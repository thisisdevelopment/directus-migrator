package directus

import (
	"context"
	"directus-migrator/pkg/config"
	"fmt"
	"os"
	"time"

	"github.com/thisisdevelopment/go-dockly/xlogger"
)

type IDirectusRepo interface {
	GetDiff(ctx context.Context, srcEnv, targetEnv string) (string, error)
	ApplyDiff(ctx context.Context, targetEnv, diff string) (string, error)
	BackupSchema(ctx context.Context, targetEnv string) error
}

type directusRepo struct {
	cfg *config.Config
	log *xlogger.Logger
}

func NewDirectusRepo(cfg *config.Config, log *xlogger.Logger) IDirectusRepo {
	return &directusRepo{
		cfg: cfg,
		log: log,
	}
}

func (r *directusRepo) GetDiff(ctx context.Context, srcEnv, targetEnv string) (string, error) {

	srcDirectus, err := NewDirectus(ctx, r.cfg, r.log, srcEnv)
	if err != nil {
		return "", err
	}

	targetDirectus, err := NewDirectus(ctx, r.cfg, r.log, targetEnv)
	if err != nil {
		return "", err
	}

	schema, err := srcDirectus.GetSchema(ctx)
	if err != nil {
		return "", err
	}

	diff, err := targetDirectus.Diff(ctx, schema)
	if err != nil {
		return "", err
	}

	return string(diff), nil

}
func (r *directusRepo) ApplyDiff(ctx context.Context, targetEnv, diff string) (string, error) {

	targetDirectus, err := NewDirectus(ctx, r.cfg, r.log, targetEnv)
	if err != nil {
		return "", err
	}

	applyRes, err := targetDirectus.Apply(ctx, diff)
	if err != nil {
		return "", err
	}
	if applyRes == "" {
		applyRes = "Diff successfully applied"
	}
	return string(applyRes), nil

}

func (r *directusRepo) BackupSchema(ctx context.Context, targetEnv string) error {

	targetDirectus, err := NewDirectus(ctx, r.cfg, r.log, targetEnv)
	if err != nil {
		return err
	}
	_ = targetDirectus

	var backupFile = fmt.Sprintf("%s/schema.%s.%s.backup.json", r.cfg.AppConfig.BackupSchemaPath, targetEnv, time.Now().Format(time.RFC3339))

	r.log.Infof("making backup of schema from %s", targetEnv)
	schema, err := targetDirectus.GetSchema(ctx)
	if err != nil {
		return err
	}

	if err := os.WriteFile(backupFile, []byte(schema), 0644); err != nil {
		return err
	}
	r.log.Infof("backup created, file saved as '%s'", backupFile)

	return nil
}
