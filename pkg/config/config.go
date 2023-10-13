package config

import (
	"fmt"
	"os"
	"strings"
)

type SchemaConfig struct {
	AuthURL     string `yaml:"auth_url"`
	SnapshotURL string `yaml:"snapshot_url"`
	DiffURL     string `yaml:"diff_url"`
	ApplyURL    string `yaml:"apply_url"`
}

type EnvConfig struct {
	DirectusURL  string `yaml:"directus_url"`
	DirectusUser string `yaml:"directus_user"`
	DirectusPass string `yaml:"directus_pass"`
}

type AppConfig struct {
	BackupSchemaPath string `yaml:"backup_schema_path"`
	LogLevel         string `yaml:"log_level"`
}

type Config struct {
	AppConfig      AppConfig            `yaml:"app_config"`
	EnvConfig      map[string]EnvConfig `yaml:"envs"`
	DirectusConfig SchemaConfig         `yaml:"directus_schema"`
}

func (c *Config) Sanitize(srcEnv, targetEnv string) error {
	checkEnv := func(env string) error {
		if _, ok := c.EnvConfig[env]; !ok {
			return fmt.Errorf("env '%s' does not exist in your config file", env)
		}
		return nil
	}

	if strings.HasSuffix(strings.TrimSpace(c.AppConfig.BackupSchemaPath), string(os.PathSeparator)) {
		c.AppConfig.BackupSchemaPath = strings.TrimSpace(c.AppConfig.BackupSchemaPath)
		c.AppConfig.BackupSchemaPath = strings.TrimRight(c.AppConfig.BackupSchemaPath, string(os.PathSeparator))
	}

	// check if envs exist
	err := checkEnv(srcEnv)
	if err != nil {
		return err
	}
	if targetEnv != "" {
		err = checkEnv(targetEnv)
		if err != nil {
			return err
		}
	}

	return nil
}
