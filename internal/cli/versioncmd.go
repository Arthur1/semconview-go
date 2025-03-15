package cli

import (
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"

	semconviewgo "github.com/Arthur1/semconview-go"
)

type VersionCmd struct{}

func (c *VersionCmd) Run(globals *Globals) error {
	printVersion()
	return nil
}

func printVersion() {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintf(writer, "semconview-go is a tool to analyze Go codes and output information about the dependencies on OpenTelemetry Semantic Conventions.\n")
	fmt.Fprintf(writer, "Version:\t%s\n", semconviewgo.Version)
	fmt.Fprintf(writer, "Go version:\t%s\n", runtime.Version())
	fmt.Fprintf(writer, "Arch:\t%s\n", runtime.GOARCH)
	writer.Flush()
}
