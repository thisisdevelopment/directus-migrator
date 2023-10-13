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
	"encoding/json"
	"fmt"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/spf13/cobra"
	"github.com/thisisdevelopment/go-dockly/xerrors/iferr"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "check status of env vs other env",
	Long:  `check status of env vs other env`,
	Run: func(cmd *cobra.Command, args []string) {
		iferr.Panic(checkTargetEnv(targetEnv))

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		ctx = context.WithValue(ctx, shared.VersionForcer("forceVersion"), forceVersion)

		repo := directus.NewDirectusRepo(cfg, log)
		log.Infof("getting diff between env '%s -> %s'", srcEnv, targetEnv)

		diff, err := repo.GetDiff(ctx, srcEnv, targetEnv)
		iferr.Panic(err)

		log.Infof("Differences between %s and %s\n-----------------\n", srcEnv, targetEnv)

		var clr map[string]interface{}
		err = json.Unmarshal([]byte(diff), &clr)
		iferr.Panic(err)

		nocolor, err := cmd.Flags().GetBool("nocolor")
		iferr.Panic(err)

		f := colorjson.NewFormatter()
		f.Indent = 2
		f.DisabledColor = nocolor

		s, err := f.Marshal(clr)
		iferr.Panic(err)

		fmt.Println(string(s))

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	statusCmd.Flags().BoolP("nocolor", "n", false, "no color on diff ouput")
}
