/*
Copyright © 2023 Th[is] Development

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"directus-migrator/pkg/banner"
	"directus-migrator/pkg/config"
	"directus-migrator/pkg/shared"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thisisdevelopment/go-dockly/xconfig"
	"github.com/thisisdevelopment/go-dockly/xerrors/iferr"
	"github.com/thisisdevelopment/go-dockly/xlogger"
)

var cfgFile string
var forceVersion bool
var srcEnv, targetEnv string
var cfg = new(config.Config)
var log *xlogger.Logger

var defaultConfig = &xlogger.Config{
	Level:  "debug",
	Format: "text",
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "directus-migrator",
	Short: "Helper application to migrate directus db schemas",
	Long:  `Helper application to migrate directus db schemas`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	banner.Display()
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&srcEnv, "source", "s", "dev", "source environment")
	rootCmd.PersistentFlags().StringVarP(&targetEnv, "target", "t", "", "target environment")
	rootCmd.PersistentFlags().BoolVarP(&forceVersion, "force", "f", false, "forces proceed on version mismatches, use at own risk")

	// err := rootCmd.MarkPersistentFlagRequired("t")
	// iferr.Panic(err)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".directus-migrator" (without extension).
		viper.AddConfigPath("./config")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	err := xconfig.LoadConfig(cfg, viper.GetViper().ConfigFileUsed())
	if err != nil {
		fmt.Printf("could not open config file err: '%s'", err.Error())
		os.Exit(1)
	}

	// check if backup path exists
	err = cfg.Sanitize(srcEnv, targetEnv)
	iferr.Exit(err)

	if cfg.AppConfig.LogLevel != "" {
		defaultConfig.Level = cfg.AppConfig.LogLevel
	}

	l, err := setupLogger(defaultConfig, cfg)
	iferr.Exit(err)

	log = l

	if !shared.FileExists(cfg.AppConfig.BackupSchemaPath) {
		log.Infof("backup dir '%s' does not exist. Creating now...", cfg.AppConfig.BackupSchemaPath)
		iferr.Exit(os.MkdirAll(cfg.AppConfig.BackupSchemaPath, 0755))
	}

}

func setupLogger(defaultConfig *xlogger.Config, cfg *config.Config) (*xlogger.Logger, error) {
	log, err := xlogger.New(defaultConfig)
	return log, err

}

func checkTargetEnv(env string) error {
	if env == "" {
		return fmt.Errorf("no target environment specified")
	}

	checkEnv := func(env string) error {
		if _, ok := cfg.EnvConfig[env]; !ok {
			return fmt.Errorf("env '%s' does not exist in your config file", env)
		}
		return nil
	}

	return checkEnv(targetEnv)

}
