package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MichelLacerda/nox/internal/ast"
	"github.com/MichelLacerda/nox/internal/parser"
	"github.com/MichelLacerda/nox/internal/scanner"
	"github.com/MichelLacerda/nox/internal/signal"
)

func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) any {
	result := i.evaluate(stmt.Expression)
	if result != nil {
		i.Runtime.HadRuntimeError = false
	}
	return result
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.FunctionStmt) any {
	function := NewFunction(i.Runtime, stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) any {
	condition := i.evaluate(stmt.Condition)

	if i.isTruthy(condition) {
		i.execute(stmt.Then)
	} else if stmt.Else != nil {
		i.execute(stmt.Else)
	}

	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) any {
	var parts []string
	for _, expr := range stmt.Expressions {
		value := i.evaluate(expr)
		parts = append(parts, i.stringify(value))
	}
	fmt.Println(strings.Join(parts, " "))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) any {
	var value any

	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
		if value != nil {
			i.Runtime.HadRuntimeError = false
		}
	}

	// Encerra a execução da função com um "Return" (que será capturado via recover)
	panic(Return{Value: value})
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) any {
	var value any
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
		if value != nil {
			i.Runtime.HadRuntimeError = false
		}
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) any {
	i.ExecuteBlock(stmt.Statements, NewEnvironment(i.Runtime, i.environment))
	return nil
}

func (i *Interpreter) VisitClassStmt(stmt *ast.ClassStmt) any {
	var superclass *Class

	if stmt.Superclass != nil {
		evaluated := i.evaluate(stmt.Superclass)
		if sc, ok := evaluated.(*Class); ok {
			superclass = sc
		} else {
			i.Runtime.ReportRuntimeError(stmt.Name, "Superclass must be a class.")
			return nil
		}
	}

	i.environment.Define(stmt.Name.Lexeme, nil) // Define a classe antes de instanciá-la

	if stmt.Superclass != nil {
		i.environment = NewEnvironment(i.Runtime, i.environment) // Cria um novo ambiente para a classe
		i.environment.Define("super", superclass)                // Define a variável 'super' no ambiente da classe
	}

	methods := MethodType{}
	for _, method := range stmt.Methods {
		fn := NewFunction(i.Runtime, method, i.environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = fn
	}

	class := NewClass(stmt.Name.Lexeme, superclass, methods)
	if stmt.Superclass != nil {
		i.environment = i.environment.Enclosing // Retorna ao ambiente anterior após definir a classe
	}
	i.environment.Assign(stmt.Name, class)
	return nil
}

func (i *Interpreter) VisitBreakStmt(stmt *ast.BreakStmt) any {
	panic(signal.BreakSignal{})
}

func (i *Interpreter) VisitContinueStmt(stmt *ast.ContinueStmt) any {
	panic(signal.ContinueSignal{})
}

func (i *Interpreter) VisitWithStmt(stmt *ast.WithStmt) any {
	resource := i.evaluate(stmt.Resource)
	env := NewEnvironment(i.Runtime, i.environment)
	env.Define(stmt.Alias.Lexeme, resource)

	defer func() {
		if file, ok := resource.(*FileObject); ok {
			if closeFn := file.GetMethod("close"); closeFn != nil {
				if callable, ok := closeFn.(Callable); ok {
					defer func() { recover() }()
					callable.Call(i, []any{})
				}
			}
		}
	}()

	i.ExecuteBlock([]ast.Stmt{stmt.Body}, env)
	return nil
}

// func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) any {
// 	for i.isTruthy(i.evaluate(stmt.Condition)) {
// 		func() {
// 			defer func() {
// 				if r := recover(); r != nil {
// 					switch r.(type) {
// 					case signal.BreakSignal:
// 						panic(r) // quebra o laço externo
// 					case signal.ContinueSignal:
// 						// simplesmente ignora, continua o loop
// 					default:
// 						panic(r)
// 					}
// 				}
// 			}()
// 			i.ExecuteBlock([]ast.Stmt{stmt.Body}, i.environment)
// 		}()
// 	}
// 	return nil
// }

func (i *Interpreter) VisitForInStmt(stmt *ast.ForInStmt) any {
	iterable := i.evaluate(stmt.Iterable)

	switch coll := iterable.(type) {
	case *ListInstance: // lista personalizada
		for index, value := range coll.Elements {
			env := NewEnvironment(i.Runtime, i.environment)
			if stmt.IndexVar != nil {
				env.Define(stmt.IndexVar.Lexeme, float64(index))
			}
			env.Define(stmt.ValueVar.Lexeme, value)
			func() {
				defer func() {
					if r := recover(); r != nil {
						switch r.(type) {
						case signal.BreakSignal:
							panic(r) // repassa pro loop pai
						case signal.ContinueSignal:
							// ignora, continua o próximo item

						default:
							panic(r) // repassa qualquer outro erro
						}
					}
				}()
				i.ExecuteBlock([]ast.Stmt{stmt.Body}, env)
			}()
		}

	case *DictInstance: // dicionário personalizado
		for key, value := range coll.Entries {
			env := NewEnvironment(i.Runtime, i.environment)
			if stmt.IndexVar != nil {
				env.Define(stmt.IndexVar.Lexeme, key)
			}
			env.Define(stmt.ValueVar.Lexeme, value)
			func() {
				defer func() {
					if r := recover(); r != nil {
						switch r.(type) {
						case signal.BreakSignal:
							panic(r) // repassa pro loop pai
						case signal.ContinueSignal:
							// ignora, continua o próximo item
						default:
							panic(r) // repassa qualquer outro erro
						}
					}
				}()
				i.ExecuteBlock([]ast.Stmt{stmt.Body}, env)
			}()
		}

	case []any: // list
		for index, value := range coll {
			env := NewEnvironment(i.Runtime, i.environment)

			if stmt.IndexVar != nil {
				env.Define(stmt.IndexVar.Lexeme, float64(index))
			}
			env.Define(stmt.ValueVar.Lexeme, value)

			func() {
				defer func() {
					if r := recover(); r != nil {
						switch r.(type) {
						case signal.BreakSignal:
							panic(r) // repassa pro loop pai
						case signal.ContinueSignal:
							// ignora, continua o próximo item
						default:
							panic(r)
						}
					}
				}()

				i.ExecuteBlock([]ast.Stmt{stmt.Body}, env)
			}()
		}

	case map[string]any: // dict
		for key, value := range coll {
			env := NewEnvironment(i.Runtime, i.environment)

			if stmt.IndexVar != nil {
				env.Define(stmt.IndexVar.Lexeme, key)
			}
			env.Define(stmt.ValueVar.Lexeme, value)

			func() {
				defer func() {
					if r := recover(); r != nil {
						switch r.(type) {
						case signal.BreakSignal:
							panic(r)
						case signal.ContinueSignal:
							// ignora
						default:
							panic(r)
						}
					}
				}()

				i.ExecuteBlock([]ast.Stmt{stmt.Body}, env)
			}()
		}

	case *StringInstance: // string
		for index, char := range coll.Value {
			env := NewEnvironment(i.Runtime, i.environment)
			if stmt.IndexVar != nil {
				env.Define(stmt.IndexVar.Lexeme, float64(index))
			}
			env.Define(stmt.ValueVar.Lexeme, string(char))
			func() {
				defer func() {
					if r := recover(); r != nil {
						switch r.(type) {
						case signal.BreakSignal:
							panic(r) // repassa pro loop pai
						case signal.ContinueSignal:
							// ignora, continua o próximo item
						default:
							panic(r) // repassa qualquer outro erro
						}
					}
				}()
				i.ExecuteBlock([]ast.Stmt{stmt.Body}, env)
			}()
		}

	case string: // string
		for index, char := range coll {
			env := NewEnvironment(i.Runtime, i.environment)
			if stmt.IndexVar != nil {
				env.Define(stmt.IndexVar.Lexeme, float64(index))
			}
			env.Define(stmt.ValueVar.Lexeme, string(char))
			func() {
				defer func() {
					if r := recover(); r != nil {
						switch r.(type) {
						case signal.BreakSignal:
							panic(r)
						case signal.ContinueSignal:
							// ignora, continua o próximo item
						default:
							panic(r) // repassa qualquer outro erro
						}
					}
				}()
				i.ExecuteBlock([]ast.Stmt{stmt.Body}, env)
			}()
		}

	case bool: // for { ... } → loop infinito
		if coll {
			for {
				env := NewEnvironment(i.Runtime, i.environment)
				func() {
					defer func() {
						if r := recover(); r != nil {
							switch r.(type) {
							case signal.BreakSignal:
								panic(r)
							case signal.ContinueSignal:
								// ignora
							default:
								panic(r)
							}
						}
					}()

					i.ExecuteBlock([]ast.Stmt{stmt.Body}, env)
				}()
			}
		}

	default:
		i.Runtime.ReportRuntimeError(stmt.ValueVar, "Object is not iterable.")
	}

	return nil
}

func (i *Interpreter) VisitImportStmt(stmt *ast.ImportStmt) any {
	path := stmt.Path.Literal.(string)
	absPath := filepath.Join(i.Runtime.WorkingDir, path)
	if !strings.HasSuffix(absPath, ".nox") {
		absPath += ".nox"
	}
	absPath, err := filepath.Abs(absPath)
	if err != nil {
		i.Runtime.ReportRuntimeError(stmt.Path, "Invalid import path.")
		return nil
	}

	if mod, ok := i.Runtime.Modules[absPath]; ok {
		if stmt.Alias != nil {
			i.environment.Define(stmt.Alias.Lexeme, mod)
		} else if wrapper, ok := mod.(*EnvironmentWrapper); ok {
			for name, val := range wrapper.Env.Values {
				i.environment.Define(name, val)
			}
		}
		return nil
	}

	// carrega o módulo normalmente
	source, err := os.ReadFile(absPath)
	if err != nil {
		i.Runtime.ReportRuntimeError(stmt.Path, "Failed to read module: "+err.Error())
		return nil
	}

	tokens, err := scanner.NewScanner([]rune(string(source))).ScanTokens()
	if err != nil {
		i.Runtime.ReportRuntimeError(stmt.Path, "Failed to tokenize module: "+err.Error())
		return nil
	}

	stmts, err := parser.NewParser(tokens).Parse()
	if err != nil {
		i.Runtime.ReportRuntimeError(stmt.Path, "Parse error in module.")
		return nil
	}

	modEnv := NewEnvironment(i.Runtime, nil)

	// Executa no escopo isolado
	prevEnv := i.environment
	i.environment = modEnv
	resolver := NewResolver(i)
	resolver.ResolveStatements(stmts)
	for _, stmt := range stmts {
		if export, ok := stmt.(*ast.ExportStmt); ok {
			i.execute(export)
		}
	}
	i.environment = prevEnv

	moduleObj := &EnvironmentWrapper{Env: modEnv}
	i.Runtime.Modules[absPath] = moduleObj

	// define no escopo original
	if stmt.Alias != nil {
		i.environment.Define(stmt.Alias.Lexeme, moduleObj)
	} else {
		for name, val := range modEnv.Values {
			i.environment.Define(name, val)
		}
	}

	return nil
}

func (i *Interpreter) VisitExportStmt(stmt *ast.ExportStmt) any {
	return i.execute(stmt.Declaration)
}
