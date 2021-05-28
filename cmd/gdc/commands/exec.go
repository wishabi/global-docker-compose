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

// ExecCmd represents the exec command
var ExecCmd = &cobra.Command{
	Use:    "exec",
	Short:  "Execute a command against a service",
	Long:   `
	Execute a command against a configured service.

	Usage: global_docker_compose exec {service} {command}

	Example: Start a Bash terminal on the redis container
	
	global_docker_compose exec redis bash
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		info := gdc.NewComposeInfo(ComposeFile, Services)
		gdc.Exec(info, args[0], args[1:])
		gdc.Cleanup()
	},
}

func init() {
	rootCmd.AddCommand(ExecCmd)
}
