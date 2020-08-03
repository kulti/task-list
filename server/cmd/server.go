package cmd

import (
	"fmt"
	"net"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/kulti/task-list/server/internal/router"
	"github.com/kulti/task-list/server/internal/services/calservice"
	"github.com/kulti/task-list/server/internal/services/sprinttmpl"
	"github.com/kulti/task-list/server/internal/storages/pgstore"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type serverCmdFlags struct {
	Port uint16 `env:"PORT" envDefault:"0"`
}

const (
	calendarCredentialPath = "/etc/tl/calendar/credentials.json"
	calendarIDsPath        = "/etc/tl/calendar/calendars.json"
)

func newServerCmd(dbFlags dbFlags) *cobra.Command {
	var serverCmdFlags serverCmdFlags
	if err := env.Parse(&serverCmdFlags); err != nil {
		zap.S().Fatalw("failed to parse server cmd flags", zap.Error(err))
	}

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Starts task list server",
		Run: func(cmd *cobra.Command, args []string) {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", serverCmdFlags.Port))
			if err != nil {
				zap.S().Fatalw("failed to listen", zap.Error(err), "port", serverCmdFlags.Port)
			}
			zap.S().Infow("listen at", "addr", listener.Addr().String())

			taskStore, err := pgstore.New(dbFlags.URL())
			if err != nil {
				zap.S().Fatalw("failed to connect to db", zap.Error(err))
			}

			sprintTmpl := sprinttmpl.New(taskStore, newCalendarService())
			router := router.New(taskStore, sprintTmpl)

			err = http.Serve(listener, router.RootHandler())
			if err != nil {
				zap.S().Fatalw("failed to graceful server shutdown", zap.Error(err))
			}
		},
	}

	return serverCmd
}

func newCalendarService() sprinttmpl.CalService {
	cs, err := calservice.New(calservice.Options{
		CredentialPath:  calendarCredentialPath,
		CalendarIDsPath: calendarIDsPath,
	})
	if err != nil {
		zap.L().Warn("failed to create calendar service", zap.Error(err))
		return nil
	}

	zap.L().Info("calendar service created")
	return cs
}
