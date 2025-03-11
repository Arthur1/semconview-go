package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Arthur1/semconview-go/internal/semconview"
	"github.com/alecthomas/kong"
)

var CLI struct {
	Pattern []string `arg:"" help:"Glob pattern to match Go files" default:"**/*.go"`
	Pretty  bool     `help:"Pretty print JSON output" default:"true"`
	Verbose bool     `help:"Enable verbose output" short:"v"`
}

func main() {
	kctx := kong.Parse(&CLI)
	ctx := context.Background()

	result, err := semconview.AnalyzeSemconvDependencies(ctx, CLI.Pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing files: %v\n", err)
		kctx.Exit(1)
	}

	for _, attr := range result.Attributes {
		fmt.Printf("Type: %s, Key: %s, Version: %s\n", attr.Type, attr.Key, attr.Version)
	}
}
