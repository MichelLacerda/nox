package ast

func (b *BlockStmt) String() string {
	var result string
	for _, stmt := range b.Statements {
		result += stmt.String() + "\n"
	}
	return "{\n" + result + "}"
}

func (c *ClassStmt) String() string {
	result := "class " + c.Name.Lexeme
	result += " {\n"
	for _, method := range c.Methods {
		result += method.String() + "\n"
	}
	result += "}"
	return result
}

func (e *ExpressionStmt) String() string {
	return e.Expression.String() + ";"
}

func (f *FunctionStmt) String() string {
	return "<function " + f.Name.Lexeme + ">"
}

func (i *IfStmt) String() string {
	result := "if " + i.Condition.String() + " {\n" + i.Then.String() + "\n}"
	if i.Else != nil {
		result += " else {\n" + i.Else.String() + "\n}"
	}
	return result
}

func (p *PrintStmt) String() string {
	var result string
	for _, expr := range p.Expressions {
		result += expr.String() + " "
	}
	return "print " + result + ";"
}

func (r *ReturnStmt) String() string {
	if r.Value != nil {
		return r.Keyword.Lexeme + " " + r.Value.String() + ";"
	}
	return r.Keyword.Lexeme + ";"
}

func (v *VarStmt) String() string {
	return "let " + v.Name.Lexeme + " = " + v.Initializer.String() + ";"
}

// func (w *WhileStmt) String() string {
// 	return "while " + w.Condition.String() + " {\n" + w.Body.String() + "\n}"
// }

func (f *ForInStmt) String() string {
	result := "for "
	if f.IndexVar != nil {
		result += f.IndexVar.Lexeme + " in "
	}
	if f.ValueVar != nil {
		result += f.ValueVar.Lexeme + " in "
	}
	result += f.Iterable.String() + " {\n" + f.Body.String() + "\n}"
	return result
}

func (l *ListStmt) String() string {
	var result string
	for _, stmt := range l.Statements {
		result += stmt.String() + "\n"
	}
	return result
}

func (b *BreakStmt) String() string {
	return b.Keyword.Lexeme + ";"
}

func (c *ContinueStmt) String() string {
	return c.Keyword.Lexeme + ";"
}

func (w *WithStmt) String() string {
	result := "with " + w.Resource.String()
	if w.Alias != nil {
		result += " as " + w.Alias.Lexeme
	}
	result += " {\n" + w.Body.String() + "\n}"
	return result
}

func (i *ImportStmt) String() string {
	result := "import " + i.Path.Lexeme
	if i.Alias != nil {
		result += " as " + i.Alias.Lexeme
	}
	return result + ";"
}

func (e *ExportStmt) String() string {
	if e.Declaration == nil {
		return "export;"
	}
	return "export " + e.Declaration.String() + ";"
}
