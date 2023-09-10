package app

import "github.com/spf13/cobra"

var builder = NewAppBuilder()

func Build() (*App, error) {
	return builder.Build()
}

// Configure is a function that will be called before the application run.
func Configure(action ConfigureFunc) {
	builder.Configure(action)
}

// OrderRun is a function that will be called by order when the application starts run.
func OrderRun(order int, action RunFunc) {
	builder.OrderRun(order, action)
}

// Run is a function that will be called when the application starts run. order is 0.
func Run(action RunFunc) {
	builder.OrderRun(0, action)
}

// PostRun is a function that will be called before the application exit.
func PostRun(action RunFunc) {
	builder.PostRun(action)
}

// AddCommand Add a command to the application.
func AddCommand(cmd *cobra.Command) {
	builder.CmdBuilder.AddCommand(cmd)
}

// AddCommandFunc Add a command function to the application.
func AddCommandFunc(use string, short string, run func()) {
	builder.CmdBuilder.AddRun(use, short, run)
}

func DefaultBuilder() *AppBuilder {
	return builder
}

// Version set the application version.
func Version(version string) {
	builder.Context.Version = version
}

// Info set the application info.
func Info(use string, short string, description string) {
	builder.Context.Name = use
	builder.Context.Short = short
	builder.Context.Description = description
}
