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

// DownCmd represents the down command
var DownCmd = &cobra.Command{
	Use:    "down",
	Short:  "Bring down Docker containers",
	Long:   `
	Bring down either specified or all Docker containers.

	Usage: global_docker_compose down {service}
				 global_docker_compose down
	`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		info := gdc.NewComposeInfo(ComposeFile, Services)
		if len(args) > 0 {
			gdc.Down(info, args[0])
		} else {
			gdc.Down(info, "")
		}
		gdc.Cleanup()
	},
}

func init() {
	rootCmd.AddCommand(DownCmd)
}
