/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

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
	"context"
	"directus-migrator/pkg/directus"
	"directus-migrator/pkg/shared"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/thisisdevelopment/go-dockly/xerrors/iferr"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restores a backup file made by apply command to given env",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		iferr.Panic(checkTargetEnv(targetEnv))

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		file, err := cmd.Flags().GetString("backup-file")
		iferr.Panic(err)
		schema, err := os.ReadFile(file)
		iferr.Panic(err)

		ctx = context.WithValue(ctx, shared.VersionForcer("forceVersion"), forceVersion)
		repo := directus.NewDirectusRepo(cfg, log)
		log.Infof("applying diff on env '%s'", targetEnv)

		response, err := repo.ApplyDiff(ctx, targetEnv, string(schema))
		iferr.Panic(err)

		fmt.Println(response)

		log.Info(`
--------		
		Sadly direct restore is not possible (yet) using this tool. but you can use the backup
		files to restore the db schema on the cli: 

		directus># node cli.js schema apply <backupfile>
				
--------
				`)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	restoreCmd.Flags().StringP("backup-file", "b", "", "backup file to restore, check your backup dir for files to restore")
	restoreCmd.MarkFlagRequired("backup-file")
}
