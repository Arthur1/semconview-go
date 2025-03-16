package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/Arthur1/semconview-go/internal/semconview"
	"github.com/rodaine/table"
	"gopkg.in/yaml.v3"
)

type ListCmd struct {
	Pattern []string `arg:"" help:"Glob pattern to match Go files. Default is **/*.go" default:"**/*.go"`
	Output  string   `help:"Output format (json, yaml, table). Default is table mode." enum:"json,yaml,table" default:"table"`
}

func (c *ListCmd) Run(globals *Globals) error {
	ctx := context.Background()
	slog.SetDefault(newLogger(globals.Verbose))

	result, err := semconview.AnalyzeSemconvDependencies(ctx, c.Pattern)
	if err != nil {
		return err
	}

	switch c.Output {
	case "json":
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal to JSON: %w", err)
		}
		fmt.Println(string(jsonData))
		return nil
	case "yaml":
		enc := yaml.NewEncoder(os.Stdout)
		enc.SetIndent(2)
		if err := enc.Encode(result); err != nil {
			return fmt.Errorf("failed to marshal to YAML: %w", err)
		}
		return nil
	default:
		tbl := table.New("Type", "Name", "Version")
		for _, attr := range result.Attributes {
			tbl.AddRow("attribute", attr.Key, attr.Version)
		}
		tbl.Print()
		return nil
	}
}
