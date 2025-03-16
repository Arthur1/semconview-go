package cli

import (
	"github.com/alecthomas/kong"
)

type Globals struct {
	Version VersionFlag `name:"version" short:"v" help:"Print version and quit"`
	Verbose bool        `help:"Make more talkative. Disabled by default." default:"false"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error {
	return nil
}
func (v VersionFlag) IsBool() bool {
	return true
}
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	printVersion()
	app.Exit(0)
	return nil
}

var cli struct {
	Globals
	List    ListCmd    `cmd:"list" help:"list the dependencies on OpenTelemetry Semantic Conventions"`
	Version VersionCmd `cmd:"version" help:"print version information"`
}

type Cli struct{}

func (c *Cli) Run() {
	kctx := kong.Parse(&cli,
		kong.Name("semconview-go"),
		kong.Description("semconview-go is a tool to analyze Go codes and output information about the dependencies on OpenTelemetry Semantic Conventions."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)
	err := kctx.Run(&cli.Globals)
	kctx.FatalIfErrorf(err)
}
