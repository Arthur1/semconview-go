package semconview

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hashicorp/go-multierror"
)

type SemconvDependencies struct {
	Attributes []SemconvAttribute
}

type SemconvAttribute struct {
	Type    string `json:"type"`
	Key     string `json:"key,omitempty"`
	Version string `json:"version,omitempty"`
}

type (
	packageIdentifier = string
)

func AnalyzeSemconvDependencies(ctx context.Context, globs []string) (SemconvDependencies, error) {
	files, err := findGoFiles(globs)
	if err != nil {
		slog.WarnContext(ctx, err.Error())
	}
	if len(files) == 0 {
		slog.WarnContext(ctx, "no files matched")
		return SemconvDependencies{}, nil
	}

	allAttrs := []SemconvAttribute{}
	var merr error
	for _, file := range files {
		attrs, parseErr := parseFile(ctx, file)
		if parseErr != nil {
			slog.WarnContext(ctx, fmt.Sprintf("Error parsing file %s: %v", file, parseErr))
			merr = multierror.Append(merr, parseErr)
		}
		allAttrs = append(allAttrs, attrs...)
	}

	uniqueAttrs := removeDuplicateAttrs(allAttrs)
	return SemconvDependencies{Attributes: uniqueAttrs}, merr
}

func removeDuplicateAttrs(attrs []SemconvAttribute) []SemconvAttribute {
	slices.SortFunc(attrs, func(a, b SemconvAttribute) int {
		return strings.Compare(a.Key+"|"+a.Version, b.Key+"|"+b.Version)
	})
	return slices.CompactFunc(attrs, func(a, b SemconvAttribute) bool {
		return a.Key == b.Key && a.Version == b.Version
	})
}

func findGoFiles(globs []string) ([]string, error) {
	var (
		goFiles []string
		merr    error
	)
	for _, glob := range globs {
		matches, err := doublestar.FilepathGlob(glob, doublestar.WithFilesOnly())
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}
		for _, match := range matches {
			if filepath.Ext(match) == ".go" {
				goFiles = append(goFiles, match)
			}
		}
	}
	return goFiles, merr
}

func parseFile(ctx context.Context, filePath string) ([]SemconvAttribute, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	v := &visitor{
		ctx:                ctx,
		imports:            map[packageIdentifier]packagePath{},
		semconvPkgResolver: newSemconvPackageResolver(),
	}
	ast.Walk(v, file)
	return v.resultAttrs, v.resultErr
}

type visitor struct {
	ctx                context.Context
	imports            map[packageIdentifier]packagePath
	semconvPkgResolver *semconvPackageResolver
	resultAttrs        []SemconvAttribute
	resultErr          error
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ImportSpec:
		return v.visitImportSpec(n)
	case *ast.CallExpr:
		return v.visitCallExpr(n)
	case *ast.SelectorExpr:
		return v.visitSelectorExpr(n)
	default:
		return v
	}
}

func (v *visitor) visitImportSpec(imp *ast.ImportSpec) ast.Visitor {
	path := strings.Trim(imp.Path.Value, "\"")
	if !strings.Contains(path, "go.opentelemetry.io/otel/semconv/v") {
		return v
	}

	var pkgIdent packageIdentifier
	if imp.Name != nil {
		// named import
		pkgIdent = imp.Name.Name
	} else {
		// standard import
		pkgIdent = filepath.Base(path)
	}

	v.imports[pkgIdent] = path
	return v
}

func (v *visitor) visitCallExpr(call *ast.CallExpr) ast.Visitor {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return v
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return v
	}
	pkgPath, ok := v.imports[ident.Name]
	if !ok {
		return v
	}

	version := filepath.Base(pkgPath)
	attrKey, ok, err := v.semconvPkgResolver.resolveAttrKeyFromFunc(v.ctx, pkgPath, sel.Sel.Name)
	if err != nil {
		v.resultErr = multierror.Append(v.resultErr, err)
		return v
	}
	if !ok {
		return v
	}

	v.resultAttrs = append(v.resultAttrs, SemconvAttribute{
		Type:    "attribute",
		Key:     attrKey,
		Version: version,
	})
	return v
}

func (v *visitor) visitSelectorExpr(sel *ast.SelectorExpr) ast.Visitor {
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return v
	}
	pkgPath, ok := v.imports[ident.Name]
	if !ok {
		return v
	}

	version := filepath.Base(pkgPath)
	attrKey, ok, err := v.semconvPkgResolver.resolveAttrKeyFromConst(v.ctx, pkgPath, sel.Sel.Name)
	if err != nil {
		v.resultErr = multierror.Append(v.resultErr, err)
		return v
	}
	if !ok {
		return v
	}

	v.resultAttrs = append(v.resultAttrs, SemconvAttribute{
		Type:    "attribute",
		Key:     attrKey,
		Version: version,
	})
	return v
}
