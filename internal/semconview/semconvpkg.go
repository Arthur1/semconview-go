package semconview

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log/slog"
	"strings"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/tools/go/packages"
)

type (
	packagePath  = string
	constantName = string
	functionName = string
	attributeKey = string
)

type semconvPackageResolver struct {
	analyzedPkgs map[packagePath]semconvPackage
}

func newSemconvPackageResolver() *semconvPackageResolver {
	return &semconvPackageResolver{
		analyzedPkgs: map[packagePath]semconvPackage{},
	}
}

type semconvPackage struct {
	constToAttrMap map[constantName]attributeKey
	funcToAttrMap  map[functionName]attributeKey
}

func (r *semconvPackageResolver) resolveAttrKeyFromFunc(ctx context.Context, pkgPath packagePath, funcName functionName) (attributeKey, bool, error) {
	var err error
	pkgResult, ok := r.analyzedPkgs[pkgPath]
	if !ok {
		pkgResult, err = analyzePackage(ctx, pkgPath)
		if err != nil {
			return "", false, err
		}
		r.analyzedPkgs[pkgPath] = pkgResult
	}
	attrKey, ok := pkgResult.funcToAttrMap[funcName]
	return attrKey, ok, nil
}

func (r *semconvPackageResolver) resolveAttrKeyFromConst(ctx context.Context, pkgPath packagePath, constName constantName) (attributeKey, bool, error) {
	var err error
	pkgResult, ok := r.analyzedPkgs[pkgPath]
	if !ok {
		pkgResult, err = analyzePackage(ctx, pkgPath)
		if err != nil {
			return "", false, err
		}
		r.analyzedPkgs[pkgPath] = pkgResult
	}
	attrKey, ok := pkgResult.constToAttrMap[constName]
	if !ok {
		return "", false, nil
	}
	return attrKey, true, nil
}

func analyzePackage(ctx context.Context, pkgPath packagePath) (semconvPackage, error) {
	pkg, err := loadPackage(ctx, pkgPath)
	if err != nil {
		return semconvPackage{}, err
	}
	constToAttrMap := createConstToAttrMap(pkg)
	funcToAttrMap := createFuncToAttrMap(pkg, constToAttrMap)
	return semconvPackage{
		constToAttrMap: constToAttrMap,
		funcToAttrMap:  funcToAttrMap,
	}, nil
}

func loadPackage(ctx context.Context, pkgPath string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load package %s: %w", pkgPath, err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("package %s not found", pkgPath)
	}
	if len(pkgs) > 1 {
		slog.WarnContext(ctx, "Multiple package founds. First hit item will be used.")
	}

	pkg := pkgs[0]
	for _, pkgErr := range pkg.Errors {
		err = multierror.Append(err, pkgErr)
	}
	return pkg, err
}

func createConstToAttrMap(pkg *packages.Package) map[constantName]attributeKey {
	m := map[constantName]attributeKey{}
	if pkg.TypesInfo == nil {
		return m
	}
	for _, obj := range pkg.TypesInfo.Defs {
		if obj == nil {
			continue
		}
		constObj, ok := obj.(*types.Const)
		if !ok {
			continue
		}
		if constObj.Type().String() != "go.opentelemetry.io/otel/attribute.Key" {
			continue
		}
		constName := constObj.Name()
		for _, syntax := range pkg.Syntax {
			ast.Inspect(syntax, func(n ast.Node) bool {
				valueSpec, ok := n.(*ast.ValueSpec)
				if !ok {
					return true
				}
				for i, name := range valueSpec.Names {
					if name.Name != constName {
						continue
					}
					if i >= len(valueSpec.Values) {
						continue
					}
					callExpr, ok := valueSpec.Values[i].(*ast.CallExpr)
					if !ok {
						continue
					}
					selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
					if !ok {
						continue
					}
					ident, ok := selectorExpr.X.(*ast.Ident)
					if !ok {
						continue
					}
					if ident.Name == "attribute" && selectorExpr.Sel.Name == "Key" {
						if len(callExpr.Args) == 0 {
							continue
						}
						basicLit, ok := callExpr.Args[0].(*ast.BasicLit)
						if !ok || basicLit.Kind != token.STRING {
							continue
						}
						attrName := strings.Trim(basicLit.Value, "\"")
						m[constName] = attrName
					}
				}
				return true
			})
		}
	}
	return m
}

func createFuncToAttrMap(pkg *packages.Package, constToAttrMap map[constantName]attributeKey) map[functionName]attributeKey {
	m := map[functionName]attributeKey{}
	for _, syntax := range pkg.Syntax {
		for _, decl := range syntax.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if funcDecl.Body == nil || len(funcDecl.Body.List) != 1 {
				continue
			}
			stmt := funcDecl.Body.List[0]
			returnStmt, ok := stmt.(*ast.ReturnStmt)
			if !ok {
				continue
			}
			if len(returnStmt.Results) == 0 {
				continue
			}
			callExpr, ok := returnStmt.Results[0].(*ast.CallExpr)
			if !ok {
				continue
			}
			selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			ident, ok := selectorExpr.X.(*ast.Ident)
			if !ok {
				continue
			}

			key, ok := constToAttrMap[ident.Name]
			if !ok {
				continue
			}
			m[funcDecl.Name.Name] = key
		}
	}
	return m
}
