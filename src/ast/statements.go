package ast

type BlockStmt struct {
	Body []Stmt
}

func (n BlockStmt) stmt() {}

type ExpressionStmt struct {
	Expression Expr
}

func (n ExpressionStmt) stmt() {}

type VarDeclStmt struct {
	VarName       string
	IsConstant    bool
	AssignedValue Expr
	ExplicitType  Type
}

func (n VarDeclStmt) stmt() {}

type StructProperty struct {
	IsStatic bool
	Type     Type
}

type StructMethod struct {
	IsStatic bool
	// Type     Type
}

type StructDeclStmt struct {
	StructName string
	Properties map[string]StructProperty
	Methods    map[string]StructMethod
}

func (n StructDeclStmt) stmt() {}

type FnDeclStmt struct {
	FnName     string
	Parameters []Parameter
	ReturnType Type
	Body       BlockStmt
}

func (n FnDeclStmt) stmt() {}

type Parameter struct {
	Name string
	Type Type
}

type ReturnStmt struct {
	Value Expr
}

func (n ReturnStmt) stmt() {}

type IfStmt struct {
	IfBody     BlockStmt
	Condition  Expr
	ElseBody   BlockStmt
	ElifBodies map[Expr]BlockStmt
}

func (i IfStmt) stmt() {}

type WhileStmt struct {
	Body      BlockStmt
	Condition Expr
}

func (i WhileStmt) stmt() {}

type ForeachStmt struct {
	Iterator   string
	Collection Expr
	Body       BlockStmt
}

func (i ForeachStmt) stmt() {}

type ForStmt struct {
	Init Stmt
	Cond Expr
	Post Stmt
	Body BlockStmt
}

func (i ForStmt) stmt() {}
