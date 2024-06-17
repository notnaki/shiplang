package ast

import "encoding/gob"

type Stmt interface {
	stmt()
}

type Expr interface {
	expr()
}

type Type interface {
	_type()
}

func init() {
	gob.Register(NumberExpr{})
	gob.Register(StringExpr{})
	gob.Register(SymbolExpr{})

	gob.Register(BinaryExpr{})
	gob.Register(PrefixExpr{})
	gob.Register(AssignmentExpr{})
	gob.Register(StructInstantiationExpr{})
	gob.Register(ArrayInstantiationExpr{})
	gob.Register(MemberAccessExpr{})
	gob.Register(ArrayAccessExpr{})
	gob.Register(CallExpr{})

	gob.Register(BlockStmt{})
	gob.Register(ExpressionStmt{})
	gob.Register(VarDeclStmt{})
	gob.Register(StructDeclStmt{})
	gob.Register(FnDeclStmt{})
	gob.Register(Parameter{})
	gob.Register(ReturnStmt{})
	gob.Register(BreakStmt{})
	gob.Register(IfStmt{})
	gob.Register(WhileStmt{})
	gob.Register(ForeachStmt{})
	gob.Register(ForStmt{})
	gob.Register(ImportStmt{})

	gob.Register(SymbolType{})
	gob.Register(ArrayType{})
}
