package router_test

import (
	"fmt"
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	config.DisableStacktrace = true
	zapLogger, err := config.Build()
	if err != nil {
		fmt.Println("failed to init db flags: ", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(zapLogger)

	os.Exit(m.Run())
}
