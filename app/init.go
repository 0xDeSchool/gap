package app

var builder = NewAppBuilder()

func Build() (*App, error) {
	return builder.Build()
}

func Configure(action ConfigureFunc) {
	builder.Configure(action)
}

func OrderRun(order int, action RunFunc) {
	builder.OrderRun(order, action)
}

func Run(action RunFunc) {
	builder.OrderRun(0, action)
}

func PostRun(action RunFunc) {
	builder.PostRun(action)
}

func DefaultBuilder() *AppBuilder {
	return builder
}

func Version(version string) {
	builder.Context.Version = version
}

func Info(use string, short string, description string) {
	builder.Context.Name = use
	builder.Context.Short = short
	builder.Context.Description = description
}
