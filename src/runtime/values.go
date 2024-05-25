package runtime

import (
	"fmt"
	"shiplang/src/ast"
	"strconv"

	"strings"
)

type RuntimeVal interface {
	Type() ValueType
	Inspect() string
}

type Variable struct {
	Value        RuntimeVal
	ExpectedType ValueType
	Constant     bool
}

type Number struct {
	Value float64
}

type Bool struct {
	Value bool
}

type Null struct{}

type String struct {
	Value string
}

type Array struct {
	Elements    []RuntimeVal
	ElementType ValueType
}

type StructDef struct {
	Name       string
	Properties map[string]ValueType
	Methods    map[string]Function
}

type Struct struct {
	Name       string
	Properties map[string]RuntimeVal
}

type Function struct {
	Name       string
	Parameters []Parameter
	Body       ast.BlockStmt
	Env        *environment
	NativeFn   NativeFunction
}

type Parameter struct {
	Name string
	Type ValueType
}

type Break struct{}

type Return struct {
	Value RuntimeVal
}

type NativeFunction struct {
	Call FunctionCall
}

type FunctionCall func([]RuntimeVal) RuntimeVal

func (v Variable) Type() ValueType {
	return VarType
}

func (v Variable) Inspect() string {
	return fmt.Sprintf("%s<%s>", VarType, v.ExpectedType)
}

func (n Number) Type() ValueType {
	return NumberType
}

func (n Number) Inspect() string {
	return fmt.Sprintf("%g", n.Value)
}

func (b Bool) Type() ValueType {
	return BooleanType
}

func (b Bool) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (n Null) Type() ValueType {
	return NullType
}

func (n Null) Inspect() string {
	return "null"
}

func (s String) Type() ValueType {
	return StringType
}

func (s String) Inspect() string {
	return s.Value
}

func (a Array) Type() ValueType {
	return ValueType(fmt.Sprintf("array<%s>", a.ElementType))
}

func (a Array) Inspect() string {
	return fmt.Sprintf("array<%s>", a.Elements)
}

func (sd StructDef) Type() ValueType {
	return ValueType(sd.Name)
}

func (sd StructDef) Inspect() string {
	var properties []string
	for name, propType := range sd.Properties {
		properties = append(properties, fmt.Sprintf("(%s, %s)", name, propType))
	}
	return fmt.Sprintf("%s<%s>", sd.Name, strings.Join(properties, ", "))
}

func (s Struct) Type() ValueType {
	return ValueType(s.Name)
}

func (s Struct) Inspect() string {
	var properties []string
	for name, propVal := range s.Properties {
		properties = append(properties, fmt.Sprintf("(%s, %s)", name, propVal.Type()))
	}
	return fmt.Sprintf("%s<%s>", s.Name, strings.Join(properties, ", "))
}

func (f Function) Type() ValueType {
	return ValueType(FunctionType)
}

func (f Function) Inspect() string {
	var params []string
	for _, param := range f.Parameters {
		params = append(params, fmt.Sprintf("(%s, %s)", param.Name, param.Type))
	}
	return fmt.Sprintf("%s<%s>", f.Name, strings.Join(params, ", "))
}

func (r Return) Type() ValueType {
	return ReturnType
}

func (r Return) Inspect() string {
	return fmt.Sprintf("return<%s>", r.Value)
}

func (b Break) Type() ValueType {
	return BreakType
}

func (b Break) Inspect() string {
	return "break"
}

func (n NativeFunction) Type() ValueType {
	return NativeFnType
}

func (n NativeFunction) Inspect() string {
	return "native-fn"
}

func (s String) CallMethod(methodName string, args ...RuntimeVal) RuntimeVal {
	switch methodName {
	case "length":
		return MKNUM(float64(len(s.Value)))
	case "concat":
		if len(args) != 1 {
			panic("concat method expects exactly 1 argument")
		}
		concatStr, ok := args[0].(String)
		if !ok {
			panic("concat method argument must be a string")
		}
		return MKSTR(string(s.Value) + string(concatStr.Value))
	case "split":
		if len(args) != 1 {
			panic("split method expects exactly 1 argument")
		}
		sep, ok := args[0].(String)
		if !ok {
			panic("split method argument must be a string")
		}
		strs := strings.Split(string(s.Value), string(sep.Value))
		values := make([]RuntimeVal, len(strs))
		for i, str := range strs {
			values[i] = MKSTR(str)
		}
		return Array{Elements: values, ElementType: StringType}
	default:
		panic(fmt.Sprintf("Method %s not found for type String", methodName))
	}
}

func (n Number) CallMethod(methodName string, args ...RuntimeVal) RuntimeVal {
	switch methodName {
	case "toString":
		return MKSTR(strconv.FormatFloat(n.Value, 'f', -1, 64))
	case "isEven":
		return MKBOOL(int(n.Value)%2 == 0)
	case "isOdd":
		return MKBOOL(int(n.Value)%2 != 0)
	default:
		panic(fmt.Sprintf("Method %s not found for type Number", methodName))
	}
}

func (b Bool) CallMethod(methodName string, args ...RuntimeVal) RuntimeVal {
	switch methodName {
	case "toString":
		return MKSTR(strconv.FormatBool(b.Value))
	default:
		panic(fmt.Sprintf("Method %s not found for type Boolean", methodName))
	}
}

func (arr *Array) CallMethod(methodName string, args ...RuntimeVal) RuntimeVal {
	switch methodName {
	case "length":
		return MKNUM(float64(len(arr.Elements)))
	case "append":
		if len(args) != 1 {
			panic("append method expects exactly 1 argument")
		}
		if arr.ElementType != args[0].Type() {
			panic("")
		}
		newArr := append(arr.Elements, args[0])
		arr.Elements = newArr
		arr.ElementType = args[0].Type()
		return Array{Elements: newArr, ElementType: arr.ElementType}
	case "pop":
		if len(arr.Elements) == 0 {
			panic("pop method cannot be called on an empty array")
		}
		newArr := arr.Elements[:len(arr.Elements)-1]
		arr.Elements = newArr
		return Array{Elements: newArr, ElementType: arr.ElementType}
	default:
		panic(fmt.Sprintf("Method %s not found for type Array", methodName))
	}
}
