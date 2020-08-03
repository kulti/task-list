package cmd

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type dbFlags struct {
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	DBName   string `env:"POSTGRES_DB,required"`
	Host     string `env:"POSTGRES_DB_HOST,required"`
}

func (d dbFlags) URL() string {
	return "postgres://" + d.User + ":" + d.Password + "@" + d.Host + ":5432/" + d.DBName
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	zapLogger, err := config.Build()
	if err != nil {
		fmt.Println("failed to init db flags: ", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(zapLogger)

	var dbFlags dbFlags
	if err := env.Parse(&dbFlags); err != nil {
		zap.S().Fatalw("failed to parse db flags", zap.Error(err))
	}

	rootCmd := &cobra.Command{
		Use: "tl",
	}

	rootCmd.AddCommand(newServerCmd(dbFlags))

	if err := rootCmd.Execute(); err != nil {
		zap.S().Fatalw("failed to execute root cmd", "err", err)
	}
}
