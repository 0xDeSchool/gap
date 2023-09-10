package app

type AppBuilder struct {
	Context    *AppContext
	CmdBuilder *CommandBuilder
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{
		Context:    NewAppContext(),
		CmdBuilder: NewCommandBuilder(),
	}
}

func (a *AppBuilder) Build() (*App, error) {
	rootCmd, err := a.CmdBuilder.Build(a.Context)
	if err != nil {
		return nil, err
	}
	app := newApp(rootCmd)
	return &app, nil
}

func (a *AppBuilder) Version(version string) *AppBuilder {
	a.Context.Version = version
	return a
}

func (a *AppBuilder) Info(use string, short string, description string) *AppBuilder {
	a.Context.Name = use
	a.Context.Short = short
	a.Context.Description = description
	return a
}

func (a *AppBuilder) Configure(action RunFunc) {
	a.Context.PreRun(action)
}

func (a *AppBuilder) Run(action RunFunc) {
	a.Context.Run(action)
}

func (a *AppBuilder) OrderRun(order int, action RunFunc) {
	a.Context.OrderRun(order, action)
}

func (a *AppBuilder) PostRun(action RunFunc) {
	a.Context.PostRun(action)
}
