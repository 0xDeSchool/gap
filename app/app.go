package app

import (
	"github.com/rs/zerolog"
	"os"

	"github.com/0xDeSchool/gap/log"

	"github.com/spf13/cobra"
)

type AppEnvironment string

const (
	Development AppEnvironment = "development"
	Production  AppEnvironment = "production"
	Staging     AppEnvironment = "staging"
)

type AppOptions struct {
	Environment AppEnvironment
	LogLevel    zerolog.Level
}

type App struct {
	rootCmd *cobra.Command
}

func newApp(root *cobra.Command) App {
	return App{
		rootCmd: root,
	}
}

func (app App) Run() {
	if err := app.rootCmd.Execute(); err != nil {
		log.Fatal(err, "程序出错: "+err.Error())
		os.Exit(1)
	}
}
