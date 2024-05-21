package ast

import "shipgo/src/lexer"

type NumberExpr struct {
	Value float64
}

func (n NumberExpr) expr() {}

type StringExpr struct {
	Value string
}

func (n StringExpr) expr() {}

type SymbolExpr struct {
	Value string
}

func (n SymbolExpr) expr() {}

// --

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (n BinaryExpr) expr() {}

type PrefixExpr struct {
	Operator  lexer.Token
	RightExpr Expr
}

func (n PrefixExpr) expr() {}

type AssignmentExpr struct {
	Assigne  Expr
	Operator lexer.Token
	Value    Expr
}

func (n AssignmentExpr) expr() {}

type StructInstantiationExpr struct {
	StructName string
	Properties map[string]Expr
}

func (n StructInstantiationExpr) expr() {}

type ArrayInstantiationExpr struct {
	Underlying Type
	Contents   []Expr
}

func (n ArrayInstantiationExpr) expr() {}

type MemberAccessExpr struct {
	Struct Expr
	Member string
}

func (n MemberAccessExpr) expr() {}

type ArrayAccessExpr struct {
	Array Expr
	Index Expr
	Prev  bool
	Rest  bool
}

func (n ArrayAccessExpr) expr() {}

type CallExpr struct {
	FunctionName string
	Arguments    []Expr
}

func (n CallExpr) expr() {}
