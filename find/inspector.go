package find

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"
)

// Call is used as argument to AllCalls
type Call struct {
	ImportPath string
	TypeName   string
	MethodName string
}

// Range TODO: Add docs
type Range struct {
	Name             string
	BeginPos, EndPos token.Position
}

func (tr Range) String() string {
	return fmt.Sprintf("%s:line-%d:col-%d:line-%d:col-%d", tr.Name, tr.BeginPos.Line, tr.BeginPos.Column, tr.EndPos.Line, tr.EndPos.Column)
}

// Finder TODO: add docs
type Finder struct {
	finderFunc func(c Call, n ast.Node, fset *token.FileSet) []*Range
	targetCall Call
}

// InFile TODO: add docs
func (f Finder) InFile(filename string) ([]*Range, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil /* src */, 0)
	// ast.Fprint(os.Stdout, fset, file, nil /* FieldFilter */)
	if err != nil {
		return nil, err
	}

	return f.finderFunc(f.targetCall, file, fset), nil
}

// AllCalls returns a slice Ranges of the queried call, given a filename.
func AllCalls(call Call) Finder {
	return Finder{
		finderFunc: func(c Call, n ast.Node, fset *token.FileSet) []*Range {
			subTestCalls := getAllCalls(call /* target */, n)

			for _, subTest := range subTestCalls {
				getRange(subTest)
			}
			testRanges := GetRangesOfTests(subTestCalls, fset)
			transformNames(testRanges)
			return testRanges
		},
	}
}

func getRange(testCall *ast.CallExpr) {
	switch firstArg := testCall.Args[0].(type) {
	case *ast.BasicLit:
		if !(firstArg.Kind == token.STRING) {
			log.Printf("Unrecognized literal kind %v", firstArg.Kind)
		}
	case *ast.SelectorExpr: // attribute access on a struct/pkg

	}
}

// GetRangesOfTests TODO: Add docs
func GetRangesOfTests(testCalls []*ast.CallExpr, fset *token.FileSet) []*Range {
	var ranges []*Range
	for _, call := range testCalls {
		switch expr := call.Args[0].(type) {
		case *ast.BasicLit:
			if !(expr.Kind == token.STRING) {
				fmt.Fprintf(os.Stderr, "Unrecognized literal kind %v", expr.Kind)
				continue
			}
			ranges = append(ranges, &Range{
				Name:     strings.Trim(expr.Value, "\""),
				BeginPos: fset.Position(call.Lparen),
				EndPos:   fset.Position(call.Rparen),
			})
		case *ast.SelectorExpr: // attribute access on a struct/pkg
			// objName
			// attributeGettingAccessed := expr.Sel.Name

			identifier, ok := expr.X.(*ast.Ident)
			if !ok {
				continue
			}
			targetAttribute := expr.Sel.Name
			if identifier.Obj == nil || identifier.Obj.Kind != ast.Var || identifier.Obj.Decl == nil {
				continue
			}
			assignment, ok := identifier.Obj.Decl.(*ast.AssignStmt)
			if !ok || len(assignment.Rhs) != 1 {
				continue
			}
			rangeExpr, ok := assignment.Rhs[0].(*ast.UnaryExpr)
			if !ok || rangeExpr.Op != token.RANGE {
				continue
			}
			variable, ok := rangeExpr.X.(*ast.Ident)
			if !ok || variable.Obj.Kind != ast.Var || variable.Obj.Decl == nil {
				continue
			}
			// ast.Print(token.NewFileSet(), variable.Obj)
			// Inspect declaration of test table
			variableDeclaration, ok := variable.Obj.Decl.(*ast.AssignStmt)
			if !ok {
				continue
			}
			valueCompositeLiteral, ok := variableDeclaration.Rhs[0].(*ast.CompositeLit)
			if !ok {
				continue
			}
			if _, ok := valueCompositeLiteral.Type.(*ast.ArrayType); !ok {
				continue
			}
			for _, v := range valueCompositeLiteral.Elts {
				composite, ok := v.(*ast.CompositeLit)
				if !ok {
					continue
				}

				for _, field := range composite.Elts {
					kv, ok := field.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					keyIdent, ok := kv.Key.(*ast.Ident)
					if !ok || keyIdent.Name != targetAttribute {
						continue
					}
					literal, ok := kv.Value.(*ast.BasicLit)
					if !ok || literal.Kind != token.STRING {
						continue
					}
					ranges = append(ranges, &Range{
						Name:     strings.Trim(literal.Value, "\""),
						BeginPos: fset.Position(v.Pos()),
						EndPos:   fset.Position(v.End()),
					})
				}
			}
		}
	}
	return ranges
}

// GetImportIdentifiers returns the tokens which correspond to `packageName`.
// e.g. the alias identifier in aliased imports.
func GetImportIdentifiers(importPath string, n ast.Node) []string {
	var names []string
	ast.Inspect(n, func(n ast.Node) bool {
		importStmt, ok := n.(*ast.GenDecl)
		if !ok || importStmt.Tok != token.IMPORT {
			return true
		}
		for _, spec := range importStmt.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}

			currImportPath := strings.Trim(importSpec.Path.Value, "\"")
			if currImportPath != importPath {
				continue
			}

			if importSpec.Name != nil {
				names = append(names, importSpec.Name.Name)
			} else {
				name := strings.Trim(importSpec.Path.Value, "\"")
				tokens := strings.Split(name, "/")

				names = append(names, tokens[len(tokens)-1])
			}
		}
		return true
	})
	return names
}

// getAllCalls returns a slice of *ast.CallExpr which represent the calls to <importPath>.<typeName>.<methodName>(...)
// The argument testingPackageFQNs represents the valid identifiers for the go/testing package
// example, given:
// 		import "testing"
// 		import aliasedTesting "testing"
//
// "testing" and "aliasedTesting" would be the FQNs
func getAllCalls(c Call, n ast.Node) []*ast.CallExpr {
	var subTestCalls []*ast.CallExpr

	importIdentifiers := GetImportIdentifiers(c.ImportPath, n)

	ast.Inspect(n, func(n ast.Node) bool {
		var ok bool
		var (
			call          *ast.CallExpr
			function      *ast.SelectorExpr
			id            *ast.Ident
			field         *ast.Field
			idType        *ast.StarExpr
			tMethodAccess *ast.SelectorExpr
			tReceiver     *ast.Ident
		)

		if call, ok = n.(*ast.CallExpr); !ok { // we are interested only in function calls
			return true
		}
		if function, ok = call.Fun.(*ast.SelectorExpr); !ok { //looking for receiver.Method
			return true
		}
		if function.Sel.Name != c.MethodName { // if the method name is not Run
			return true
		}
		if id, ok = function.X.(*ast.Ident); !ok || id.Obj == nil || id.Obj.Kind != ast.Var || id.Obj.Decl == nil {
			return true
		}
		if field, ok = id.Obj.Decl.(*ast.Field); !ok {
			return true
		}
		if idType, ok = field.Type.(*ast.StarExpr); !ok { // must be pointer to struct
			return true
		}
		if tMethodAccess, ok = idType.X.(*ast.SelectorExpr); !ok {
			return true
		}
		if tReceiver, ok = tMethodAccess.X.(*ast.Ident); !ok {
			return true
		}
		typePackageName, typeName := tReceiver.Name, tMethodAccess.Sel.Name
		if typeName != c.TypeName {
			return true
		}
		if !contains(importIdentifiers, typePackageName) {
			return true
		}

		subTestCalls = append(subTestCalls, call)

		// fmt.Println(tReceiver.Name, tMethodAccess.Sel.Name, function.Sel.Name)
		// fmt.Printf("\n\n\n")
		// ast.Fprint(os.Stdout, fset, tMethodAccess, nil /* FieldFilter */)
		return true
	})

	return subTestCalls
}

func contains(strings []string, s string) bool {
	for _, curr := range strings {
		if curr == s {
			return true
		}
	}
	return false
}

// transformNames transforms the names of the tests to the format understood by the go test command.
// Quotes regex meta characters and replaces spaces with underscores.
func transformNames(ts []*Range) {
	for _, t := range ts {
		t.Name = regexp.QuoteMeta(strings.ReplaceAll(t.Name, " ", "_"))
	}
}

func definitionOfTestTable(ttVar *ast.Ident) ast.Expr {
	if ttVar.Obj.Kind != ast.Var || ttVar.Obj.Decl == nil {
		return nil
	}
	switch decl := ttVar.Obj.Decl.(type) {
	case *ast.AssignStmt:
		return decl.Rhs[0]
	case *ast.ValueSpec:
		for i, ident := range decl.Names {
			if ident.Name != ttVar.Name {
				continue
			}
			return decl.Values[i]
		}
		log.Println(ttVar.Name)
	}
	return nil
}
