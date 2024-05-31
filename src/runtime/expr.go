package runtime

import (
	"fmt"
	"shiplang/src/ast"
	"shiplang/src/lexer"
)

func eval_expr(expr ast.Expr, env *environment) RuntimeVal {

	switch e := expr.(type) {
	case ast.NumberExpr:
		return MKNUM(e.Value)
	case ast.StringExpr:
		return MKSTR(e.Value)
	case ast.SymbolExpr:
		return eval_symbol_expr(e, env)
	case ast.PrefixExpr:
		return eval_prefix_expr(e, env)
	case ast.BinaryExpr:
		return eval_binary_expr(e, env)
	case ast.ArrayInstantiationExpr:
		return eval_array_inst_expr(e, env)
	case ast.ArrayAccessExpr:
		return eval_array_access_expr(e, env)
	case ast.StructInstantiationExpr:
		return eval_struct_inst_expr(e, env)
	case ast.CallExpr:
		return eval_call_expr(e, env)
	case ast.MemberAccessExpr:
		return eval_member_access_expr(e, env)
	case ast.AssignmentExpr:
		return eval_assignment_expr(e, env)

	default:
		panic(fmt.Sprintf("expr not set up %s", e))
	}
}

func eval_symbol_expr(sym ast.SymbolExpr, env *environment) RuntimeVal {
	return env.lookupVar(sym.Value).Value
}

func eval_prefix_expr(pr ast.PrefixExpr, env *environment) RuntimeVal {
	right := eval_expr(pr.RightExpr, env)
	switch pr.Operator.Kind {
	case lexer.NOT:
		return MKBOOL(!truthify(right))
	case lexer.DASH:
		return negate(right)
	case lexer.PLUS:
		return right
	default:
		panic("Unknown prefix operator")
	}
}

func eval_binary_expr(b ast.BinaryExpr, env *environment) RuntimeVal {
	lhs := eval_expr(b.Left, env)
	rhs := eval_expr(b.Right, env)
	lhsType, rhsType := lhs.Type(), rhs.Type()

	if lhsType == NumberType && rhsType == NumberType {
		return eval_numeric_binary_expr(lhs.(Number), rhs.(Number), b.Operator)
	}

	return eval_logical_binary_expr(lhs, rhs, b.Operator)
}

func eval_numeric_binary_expr(lhs Number, rhs Number, op lexer.Token) RuntimeVal {
	var res float64

	switch op.Kind {
	case lexer.PLUS:
		res = lhs.Value + rhs.Value
	case lexer.DASH:
		res = lhs.Value - rhs.Value
	case lexer.STAR:
		res = lhs.Value * rhs.Value
	case lexer.SLASH:
		if rhs.Value == 0 {
			panic("Division by zero")
		}
		res = lhs.Value / rhs.Value
	case lexer.PERCENT:
		if rhs.Value == 0 {
			panic("Division by zero")
		}
		res = float64(int(lhs.Value) % int(rhs.Value))
	default:
		return eval_comparison_binary_expr(lhs, rhs, op)
	}

	return MKNUM(res)
}

func eval_comparison_binary_expr(lhs Number, rhs Number, op lexer.Token) RuntimeVal {
	var res bool

	switch op.Kind {
	case lexer.GREATER:
		res = lhs.Value > rhs.Value
	case lexer.GREATER_EQUALS:
		res = lhs.Value >= rhs.Value
	case lexer.LESS:
		res = lhs.Value < rhs.Value
	case lexer.LESS_EQUALS:
		res = lhs.Value <= rhs.Value
	case lexer.EQUALS:
		res = equals(lhs, rhs)
	case lexer.NOT_EQUALS:
		res = !equals(lhs, rhs)
	default:
		panic("Unknown comparison operator")
	}

	return MKBOOL(res)
}

func eval_logical_binary_expr(lhs RuntimeVal, rhs RuntimeVal, op lexer.Token) RuntimeVal {
	var res bool

	switch op.Kind {
	case lexer.AND:
		res = truthify(lhs) && truthify(rhs)
	case lexer.OR:
		res = truthify(lhs) || truthify(rhs)
	case lexer.EQUALS:
		res = equals(lhs, rhs)
	case lexer.NOT_EQUALS:
		res = !equals(rhs, lhs)
	default:
		panic("Unknown logical operator")
	}

	return MKBOOL(res)
}

func eval_array_inst_expr(ai ast.ArrayInstantiationExpr, env *environment) RuntimeVal {
	elements := make([]RuntimeVal, 0, len(ai.Contents))
	elementType := extractValueType(ai.Underlying)

	for _, element := range ai.Contents {
		val := eval_expr(element, env)
		if !checkType(val.Type(), elementType) {
			panic("Type mismatch in array elements")
		}
		elements = append(elements, val)
	}

	return Array{Elements: elements, ElementType: elementType}
}

func eval_array_access_expr(aa ast.ArrayAccessExpr, env *environment) RuntimeVal {
	a := eval_expr(aa.Array, env)
	i := eval_expr(aa.Index, env)

	if i.Type() != NumberType {
		panic("Array index must be a number")
	}

	index := int(i.(Number).Value)

	switch a := a.(type) {
	case String:
		return eval_string_access_expr(a, index, aa.Rest, aa.Prev)
	case Array:
		if aa.Prev {
			return Array{Elements: a.Elements[:index], ElementType: a.ElementType}
		}
		if aa.Rest {
			return Array{Elements: a.Elements[index:], ElementType: a.ElementType}
		}
		return a.Elements[index]
	default:
		panic("Invalid array access")
	}
}

func eval_string_access_expr(s String, index int, rest bool, prev bool) RuntimeVal {
	if index < 0 || index >= len(s.Value) {
		panic("String index out of bounds")
	}
	if rest {
		return MKSTR(s.Value[:index])
	}
	if prev {
		return MKSTR(s.Value[index:])
	}
	return MKSTR(string(s.Value[index]))
}

func eval_struct_inst_expr(si ast.StructInstantiationExpr, env *environment) RuntimeVal {
	sd := env.lookupStruct(si.StructName)
	structDef, ok := sd.(StructDef)
	if !ok {
		panic(fmt.Sprintf("Struct %s not found", si.StructName))
	}

	evalProps := make(map[string]RuntimeVal, len(structDef.Properties))

	for name, expectedType := range structDef.Properties {
		if prop, ok := si.Properties[name]; ok {
			propVal := eval_expr(prop, env)
			if !checkType(propVal.Type(), expectedType) {
				panic(fmt.Sprintf("Type mismatch for property %s in struct %s", name, si.StructName))
			}
			evalProps[name] = propVal
		} else {
			evalProps[name] = MKNULL()
		}
	}

	return Struct{
		Name:       si.StructName,
		Properties: evalProps,
	}
}

func eval_call_expr(c ast.CallExpr, env *environment) RuntimeVal {
	if c.Struct != nil {
		return handle_method_call(c, env)
	}

	v := env.containsVar(c.FunctionName)
	if v {
		// Check if the variable value is a function reference
		v := env.Variables[c.FunctionName]
		if fnRef, ok := v.Value.(Function); ok {
			// Evaluate arguments
			callEnv := &environment{Variables: make(map[string]Variable), Parent: env}
			for i, param := range fnRef.Parameters {

				callEnv.declareVar(param.Name, eval_expr(c.Arguments[i], callEnv), param.Type, false)
			}
			// Call the referenced function
			return Evaluate(fnRef.Body, callEnv)
		} else {
			// Handle the case where the variable is not a function reference
			panic(fmt.Sprintf("%s is not a function reference", c.FunctionName))
		}
	}

	fnVal, ok := env.lookupFn(c.FunctionName).(Function)

	if !ok {
		panic(fmt.Sprintf("%s is not callable", c.FunctionName))
	}

	if fnVal.NativeFn.Call != nil {

		var evalArgs []RuntimeVal
		for _, arg := range c.Arguments {
			evalArgs = append(evalArgs, eval_expr(arg, env))
		}
		return fnVal.NativeFn.Call(evalArgs)
	}

	argCount := len(c.Arguments)
	paramCount := len(fnVal.Parameters)

	if argCount != paramCount {
		panic(fmt.Sprintf("Incorrect number of arguments for function %s: expected %d, got %d", c.FunctionName, paramCount, argCount))
	}

	callEnv := &environment{Variables: make(map[string]Variable), Parent: env}

	for i, param := range fnVal.Parameters {

		callEnv.declareVar(param.Name, eval_expr(c.Arguments[i], callEnv), param.Type, false)
	}

	return Evaluate(fnVal.Body, callEnv)

}

func handle_method_call(c ast.CallExpr, env *environment) RuntimeVal {
	v := eval_expr(c.Struct, env)
	structType := string(v.Type())

	if isPrimitive(v.Type()) {
		return handle_primitive_method_call(v, c, env)
	}

	structDef, ok := env.lookupStruct(structType).(StructDef)
	if !ok {
		panic(fmt.Sprintf("Struct %s not found", structType))
	}

	function, exists := structDef.Methods[c.FunctionName]
	if !exists {
		panic(fmt.Sprintf("Method %s not found in struct %s", c.FunctionName, structType))
	}

	callEnv := &environment{Variables: make(map[string]Variable), Parent: env}

	for i, param := range function.Parameters {
		callEnv.declareVar(param.Name, eval_expr(c.Arguments[i], callEnv), param.Type, false)
	}

	return Evaluate(function.Body, callEnv)
}

func handle_primitive_method_call(v RuntimeVal, c ast.CallExpr, env *environment) RuntimeVal {
	switch v := v.(type) {
	case String:
		args := make([]RuntimeVal, len(c.Arguments))
		for i, arg := range c.Arguments {
			args[i] = eval_expr(arg, env)
		}
		return v.CallMethod(c.FunctionName, args...)
	case Number:
		args := make([]RuntimeVal, len(c.Arguments))
		for i, arg := range c.Arguments {
			args[i] = eval_expr(arg, env)
		}
		return v.CallMethod(c.FunctionName, args...)
	case Array:
		args := make([]RuntimeVal, len(c.Arguments))
		for i, arg := range c.Arguments {
			args[i] = eval_expr(arg, env)
		}
		return v.CallMethod(c.FunctionName, args...)
	case Bool:
		args := make([]RuntimeVal, len(c.Arguments))
		for i, arg := range c.Arguments {
			args[i] = eval_expr(arg, env)
		}
		return v.CallMethod(c.FunctionName, args...)
	default:
		return MKNULL()
	}
}

func eval_member_access_expr(ma ast.MemberAccessExpr, env *environment) RuntimeVal {

	structVal := eval_expr(ma.Struct, env)

	structInstance, ok := structVal.(Struct)

	if !ok {
		panic(fmt.Sprintf("%v is not a struct instance", ma.Struct))
	}

	memberVal, exists := structInstance.Properties[ma.Member]
	if !exists {
		panic(fmt.Sprintf("Member %s not found in struct %s", ma.Member, structInstance.Name))
	}

	return memberVal
}

func eval_assignment_expr(expr ast.AssignmentExpr, env *environment) RuntimeVal {
	switch a := expr.Assigne.(type) {
	case ast.SymbolExpr:
		var val RuntimeVal
		if expr.Operator.Kind == lexer.PLUS_EQUALS {
			val = eval_binary_expr(ast.BinaryExpr{Left: expr.Assigne, Right: expr.Value, Operator: lexer.NewToken(lexer.PLUS, "+")}, env)
		} else if expr.Operator.Kind == lexer.MINUS_EQUALS {
			val = eval_binary_expr(ast.BinaryExpr{Left: expr.Assigne, Right: expr.Value, Operator: lexer.NewToken(lexer.PLUS, "-")}, env)
		} else {
			val = eval_expr(expr.Value, env)
		}
		return env.assignVar(a.Value, val)
	case ast.MemberAccessExpr:
		val := eval_expr(expr.Value, env)
		return env.assignStruct(getBaseVariableName(a), a.Member, val)
	case ast.ArrayAccessExpr:
		array := eval_expr(a.Array, env)
		index := eval_expr(a.Index, env)
		value := eval_expr(expr.Value, env)

		switch arr := array.(type) {
		case Array:
			if int(index.(Number).Value) >= len(arr.Elements) || int(index.(Number).Value) < 0 {
				panic(fmt.Sprintf("Index out of range: %d", int(index.(Number).Value)))
			}
			arr.Elements[int(index.(Number).Value)] = value
			return env.assignVar(a.Array.(ast.SymbolExpr).Value, arr)
		default:
			panic("")
		}
	default:
		panic("")
	}
}
