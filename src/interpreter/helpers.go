package interpreter

import (
	"fmt"
	"shipgo/src/ast"
)

func AstToRuntimeParams(astParam ast.Parameter) Parameter {
	return Parameter{
		ParamName: astParam.Name,
		ParamType: ValueType(astParam.Type.(ast.SymbolType).Name),
	}
}

func AstToRuntimeProps(astProp ast.StructProperty, propName string) Property {
	return Property{PropName: propName, PropType: ValueType(astProp.Type.(ast.SymbolType).Name)}
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

func equals(lhs RuntimeVal, rhs RuntimeVal) bool {
	switch lhs := lhs.(type) {
	case NumberVal:
		rhs, ok := rhs.(NumberVal)
		return ok && lhs.Value == rhs.Value
	case BooleanVal:
		rhs, ok := rhs.(BooleanVal)
		return ok && lhs.Value == rhs.Value
	case StringVal:
		rhs, ok := rhs.(StringVal)
		return ok && lhs.Value == rhs.Value
	default:
		return false
	}
}

func truthify(val RuntimeVal) bool {
	switch v := val.(type) {
	case NumberVal:
		return v.Value != 0
	case BooleanVal:
		return v.Value
	case StringVal:
		return v.Value != ""
	default:
		return false
	}
}

func negate(r RuntimeVal) RuntimeVal {
	switch r.(type) {
	case NumberVal:
		return MKNUM(-r.(NumberVal).Value)
	case BooleanVal:
		return MKBOOL(!r.(BooleanVal).Value)
	default:
		panic(fmt.Sprintf("unexpected prefix operator for: %s", r.Type()))
	}
}

func assignNestedMember(properties map[string]RuntimeVal, parts []string, value RuntimeVal) {

	if len(parts) == 1 {
		// Base case: Assign value to the last member
		properties[parts[0]] = value
	} else {
		// Recursive case: Traverse through nested members
		member := parts[0]

		nestedVal := properties[member]
		if nestedVal.Type() != StructType {
			panic(fmt.Sprintf("Cannot access member %s of non-struct type.", member))
		}
		assignNestedMember(nestedVal.(StructInstance).Properties, parts[1:], value)
	}
}

func getBaseVariableAndNestedMembers(expr ast.MemberAccessExpr) (string, []string) {
	var baseVarName string
	var nestedMembers []string

	// Handle the base variable
	switch e := expr.Struct.(type) {
	case ast.SymbolExpr:
		baseVarName = e.Value
	case ast.MemberAccessExpr:
		baseVarName, nestedMembers = getBaseVariableAndNestedMembers(e)
	default:
		panic("Invalid member access expression.")
	}

	// Append the current member to the nested members
	nestedMembers = append(nestedMembers, expr.Member)
	return baseVarName, nestedMembers
}
