package interpreter

import (
	"fmt"
	"shipgo/src/ast"
	"strings"
)

type RuntimeVal interface {
	Type() ValueType
	Inspect() string
}

type NullVal struct{}

func (n NullVal) Type() ValueType {
	return NullType
}

func (n NullVal) Inspect() string {
	return "null"
}

type StringVal struct {
	Value string
}

func (s StringVal) Type() ValueType {
	return StringType
}

func (n StringVal) Inspect() string {
	return fmt.Sprintf("%s", n.Value)
}

type NumberVal struct {
	Value float64
}

func (s NumberVal) Type() ValueType {
	return NumberType
}

func (n NumberVal) Inspect() string {
	return fmt.Sprintf("%g", n.Value)
}

type BooleanVal struct {
	Value bool
}

func (s BooleanVal) Type() ValueType {
	return BooleanType
}

func (n BooleanVal) Inspect() string {
	return fmt.Sprintf("%t", n.Value)
}

type Variable struct {
	VarType ValueType
	Value   RuntimeVal
}

func (s Variable) Type() ValueType {
	return VarType
}

func (n Variable) Inspect() string {
	return fmt.Sprintf("<variable %s - %g>", n.VarType, n.Value)
}

type Parameter struct {
	ParamName string
	ParamType ValueType
}

type Function struct {
	Name       string
	Parameters []Parameter
	Body       ast.BlockStmt
	DecEnv     *environment
}

func (s Function) Type() ValueType {
	return FunctionType
}

func (f Function) Inspect() string {
	return fmt.Sprintf("<function %s>", f.Name)
}

type NativeFn struct {
	Call FunctionCall
}

func (n NativeFn) Type() ValueType {
	return NativeFnType
}

func (n NativeFn) Inspect() string {
	return "<native fn>"
}

type FunctionCall func([]RuntimeVal) RuntimeVal

type Property struct {
	PropName string
	PropType ValueType
}

type Struct struct {
	Name       string
	Properties []Property
}

func (s Struct) Type() ValueType {
	return StructType
}

func (s Struct) Inspect() string {
	return fmt.Sprintf("<struct %s>", s.Name)
}

type StructInstance struct {
	StructName string
	Properties map[string]RuntimeVal
}

func (s StructInstance) Type() ValueType {
	return StructType
}

func (s StructInstance) Inspect() string {
	var props []string
	for key, value := range s.Properties {
		props = append(props, fmt.Sprintf("%s: %s", key, value.Inspect()))
	}
	return fmt.Sprintf("<struct %s {%s}>", s.StructName, strings.Join(props, ", "))
}

type Return struct {
	Value RuntimeVal
}

func (s Return) Type() ValueType {
	return ReturnType
}

func (n Return) Inspect() string {
	return fmt.Sprintf("<return %g>", n.Value)
}

type Break struct {
}

func (b Break) Type() ValueType {
	return BreakType
}

func (n Break) Inspect() string {
	return "<break>"
}

type ArrayElement struct {
	ElementType ValueType
	Value       RuntimeVal
}

// func (ae ArrayElement) Type() ValueType {
// 	return ArrayElementType
// }

// func (ae ArrayElement) Inspect() string {
// 	return fmt.Sprintf("<return %g>", ae.Value)
// }

type Array struct {
	ValType  ValueType
	Elements []ArrayElement
}

func (v Array) Type() ValueType {
	return ValueType(v.Inspect())
}

func (v Array) Inspect() string {
	return fmt.Sprintf("array<%s>", v.ValType)
}
