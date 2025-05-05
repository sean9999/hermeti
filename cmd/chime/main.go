package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"runtime"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/hermeti/fmt"
)

/**
 *	Here's an example of shadowing the "fmt" package so can continue to use non-hermitc code in a hermitc way
 *	The "github.com/sean9999/hermeti/fmt" package needs to be initialized with the output stream you want to redirect to
 *	ex: fmt.SetOutput(cli.Env.OutStream)
 **/

type exe struct {
	hermeti.PassthroughInit
}

func (s *exe) Run(env hermeti.Env) {
	fmtRoot := fmt.Sprintf("%s/src/fmt", runtime.GOROOT())
	symbols := dump(fmtRoot)
	for _, symbol := range symbols {
		fmt.Printf("var %s = fmt.%s\n", symbol, symbol)
	}
}

func (s *exe) State() *exe {
	return s
}

func main() {

	cli := hermeti.NewRealCli(new(exe))
	fmt.SetOutput(cli.Env.OutStream)
	cli.Run()

}

func dump(dir string) []string {

	symbols := []string{}

	// positions are relative to fset
	fset := token.NewFileSet()

	notMain := func(f fs.FileInfo) bool {
		if f.Name() == "main.go" {
			return false
		}
		return !f.IsDir()
	}

	pkgs, err := parser.ParseDir(fset, dir, notMain, parser.SkipObjectResolution)
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.GenDecl:
					if x.Tok == token.CONST || x.Tok == token.VAR {
						for _, spec := range x.Specs {
							vspec := spec.(*ast.ValueSpec)
							for _, name := range vspec.Names {
								if name.IsExported() {
									symbols = append(symbols, name.Name)
								}
							}
						}
					}
				case *ast.FuncDecl:
					if x.Recv == nil && x.Name.IsExported() {
						symbols = append(symbols, x.Name.Name)
					}
				case *ast.TypeSpec:
					if x.Name.IsExported() {
						switch x.Type.(type) {
						case *ast.InterfaceType:
							symbols = append(symbols, x.Name.Name)
						}
					}
				}
				return true
			})
		}
	}

	return symbols

}
