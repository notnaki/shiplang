package runtime

import (
	"fmt"
	"shiplang/src/ast"
)

func checkType(valType ValueType, expectedType ValueType) bool {
	if expectedType == AnyType || valType == NullType {
		return true
	}

	return valType == expectedType
}

func extractValueType(t ast.Type) ValueType {

	switch expType := t.(type) {
	case ast.SymbolType:
		return ValueType(expType.Name)
	case ast.ArrayType:

		return ValueType(fmt.Sprintf("array<%s>", extractValueType(expType.Underlying)))
	default:
		panic("Unsupported type for variable declaration")
	}
}

func truthify(val RuntimeVal) bool {
	switch v := val.(type) {
	case Number:
		return v.Value != 0
	case Bool:
		return v.Value
	case String:
		return v.Value != ""
	default:
		return false
	}
}

func equals(lhs RuntimeVal, rhs RuntimeVal) bool {
	switch lhs := lhs.(type) {
	case Number:
		rhs, ok := rhs.(Number)
		return ok && lhs.Value == rhs.Value
	case Bool:
		rhs, ok := rhs.(Bool)
		return ok && lhs.Value == rhs.Value
	case String:
		rhs, ok := rhs.(String)
		return ok && lhs.Value == rhs.Value
	default:
		return false
	}
}

func negate(r RuntimeVal) RuntimeVal {
	switch r := r.(type) {
	case Number:
		return MKNUM(-r.Value)
	case Bool:
		return MKBOOL(!r.Value)
	default:
		panic(fmt.Sprintf("unexpected prefix operator for: %s", r.Type()))
	}
}

func isPrimitive(v ValueType) bool {
	return v == NullType || v == StringType || v == NumberType || v == BooleanType || (len(v) >= 5 && v[:5] == ArrayType)
}

func getBaseVariableName(expr ast.MemberAccessExpr) string {
	if symExpr, ok := expr.Struct.(ast.SymbolExpr); ok {
		return symExpr.Value
	} else if memberExpr, ok := expr.Struct.(ast.MemberAccessExpr); ok {

		return getBaseVariableName(memberExpr)
	} else {
		return ""
	}
}

func declareNativeValues(env environment) {
	env.declareVar("true", MKBOOL(true), BooleanType, true)
	env.declareVar("false", MKBOOL(false), BooleanType, true)
	env.declareVar("null", MKNULL(), NullType, true)
}
func declareNativeFunctions(env environment) {
	env.declareNativeFn("show", showFN)
	env.declareNativeFn("time", timeFN)
	env.declareNativeFn("date", dateFN)
	env.declareNativeFn("range", rangeFN)
}
