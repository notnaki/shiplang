package interpreter

type ValueType string

const (
	AnyType          ValueType = "any"
	NullType         ValueType = "null"
	StringType       ValueType = "string"
	NumberType       ValueType = "number"
	BooleanType      ValueType = "boolean"
	StructType       ValueType = "struct"
	NativeFnType     ValueType = "native-fn"
	FunctionType     ValueType = "function"
	ArrayType        ValueType = "array"
	ReturnType       ValueType = "return"
	VarType          ValueType = "variable"
	ArrayElementType ValueType = "array-element"
)

func MKNULL() RuntimeVal {
	return NullVal{}
}

func MKSTR(s string) RuntimeVal {
	return StringVal{Value: s}
}

func MKNUM(n float64) RuntimeVal {
	return NumberVal{Value: n}
}

func MKBOOL(b bool) RuntimeVal {
	return BooleanVal{Value: b}
}

func MKNATIVEFN(c FunctionCall) RuntimeVal {
	return NativeFn{Call: c}
}
