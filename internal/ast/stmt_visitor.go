package ast

type StmtVisitor interface {
	VisitBlockStmt(stmt *BlockStmt) any
	VisitClassStmt(stmt *ClassStmt) any
	VisitExpressionStmt(stmt *ExpressionStmt) any
	VisitFunctionStmt(stmt *FunctionStmt) any
	VisitIfStmt(stmt *IfStmt) any
	VisitPrintStmt(stmt *PrintStmt) any
	VisitReturnStmt(stmt *ReturnStmt) any
	VisitVarStmt(stmt *VarStmt) any
	VisitWhileStmt(stmt *WhileStmt) any
	VisitForInStmt(stmt *ForInStmt) any
	VisitBreakStmt(stmt *BreakStmt) any
	VisitContinueStmt(stmt *ContinueStmt) any
	VisitWithStmt(stmt *WithStmt) any
	VisitImportStmt(stmt *ImportStmt) any
	VisitExportStmt(stmt *ExportStmt) any
}

func (b *BlockStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitBlockStmt(b)
}

func (c *ClassStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitClassStmt(c)
}

func (e *ExpressionStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitExpressionStmt(e)
}

func (f *FunctionStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitFunctionStmt(f)
}

func (i *IfStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitIfStmt(i)
}

func (p *PrintStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitPrintStmt(p)
}

func (r *ReturnStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitReturnStmt(r)
}

func (v *VarStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitVarStmt(v)
}

func (w *WhileStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitWhileStmt(w)
}

func (f *ForInStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitForInStmt(f)
}

func (b *BreakStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitBreakStmt(b)
}

func (c *ContinueStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitContinueStmt(c)
}

func (w *WithStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitWithStmt(w)
}

func (i *ImportStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitImportStmt(i)
}

func (e *ExportStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitExportStmt(e)
}
