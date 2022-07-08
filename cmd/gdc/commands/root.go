package commands

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
// Services list of services (space delimited)
var Services string
// ComposeFile optional additional docker-compose.yml file
var ComposeFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.6.3",
	Use:     "global_docker_compose (command) --services service1 service2 --compose_file ../docker-compose.yml",
	Short:   "Generate JSON files to use with the Flipp platform deploy scripts",
	Long: `
         global_docker_compose can be used to centralize and standardize Docker dependencies used within Flipp.
         The idea is to have one tool that can spin up whatever services are needed and keep that tool updated
         with fixes and improvements, rather than having a separate docker-compose.yml file in every project.
         The project can have a shell script (denoted "gdc") which calls this tool with the services it needs
         for itself.

         For more information, please see https://github.com/wishabi/global-docker-compose .
	`,
	Args:    cobra.ExactArgs(1),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&Services, "services", "s", "", "Services to perform actions for (required)")
	rootCmd.MarkFlagRequired("input")

	rootCmd.PersistentFlags().StringVarP(&ComposeFile, "compose_file", "c", "", "Additional docker-compose file to use")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".datadog" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gdc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
