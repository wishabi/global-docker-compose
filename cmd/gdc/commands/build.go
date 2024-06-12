/*
Package commands Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package commands

import (
	"github.com/spf13/cobra"
	"github.com/wishabi/global-docker-compose/gdc"
)

var NoCache bool

// UpCmd represents the up command
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a target service",
	Long:  `Build the image for the given target service. When --no-cache is added, the image will be built without using the cache.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cobra.MinimumNArgs(1)(cmd, args)
		}
		if len(args) > 1 {
			return cobra.MaximumNArgs(1)(cmd, args)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		info := gdc.NewComposeInfo(ComposeFile, Services)
		gdc.Build(args[0], info, NoCache)
	},
}

func init() {
	BuildCmd.Flags().BoolVarP(&NoCache, "no-cache", "n", false, "Build the image without using the cache")
	rootCmd.AddCommand(BuildCmd)
}
