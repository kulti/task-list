package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/kulti/task-list/internal/router"
	"github.com/kulti/task-list/internal/storages/pgstore"
	"github.com/spf13/cobra"
)

type serverCmdFlags struct {
	Port uint16 `env:"PORT" envDefault:"0"`
}

func newServerCmd(dbFlags dbFlags) *cobra.Command {
	var serverCmdFlags serverCmdFlags
	if err := env.Parse(&serverCmdFlags); err != nil {
		fmt.Println(err)
		os.Exit(1) //nolint:gomnd
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Starts task list server",
		Run: func(cmd *cobra.Command, args []string) {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", serverCmdFlags.Port))
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			log.Printf("Listen at %s", listener.Addr().String())

			taskStore, err := pgstore.New(dbFlags.URL)
			if err != nil {
				log.Fatalf("failed to connect to db: %v", err)
			}
			router := router.New(taskStore)

			log.Fatal(http.Serve(listener, router.RootHandler()))
		},
	}

	return serverCmd
}
