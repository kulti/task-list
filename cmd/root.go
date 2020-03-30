package cmd

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

type dbFlags struct {
	URL string `env:"DB_URL,required"`
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var dbFlags dbFlags
	if err := env.Parse(&dbFlags); err != nil {
		fmt.Println(err)
		os.Exit(1) //nolint:gomnd
	}

	var rootCmd = &cobra.Command{
		Use: "tl",
	}

	rootCmd.AddCommand(newServerCmd(dbFlags))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1) //nolint:gomnd
	}
}
