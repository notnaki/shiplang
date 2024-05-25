package runtime

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
	BreakType        ValueType = "break"
)

func MKNULL() RuntimeVal {
	return Null{}
}

func MKSTR(s string) RuntimeVal {
	return String{Value: s}
}

func MKNUM(n float64) RuntimeVal {
	return Number{Value: n}
}

func MKBOOL(b bool) RuntimeVal {
	return Bool{Value: b}
}
