package interpreter

import (
	"fmt"
	"shipgo/src/ast"
	"shipgo/src/lexer"
)

func eval_expr(expr ast.Expr, env *environment) RuntimeVal {
	// println("--------------- Current Env ---------------")
	// litter.Dump(env)
	switch e := expr.(type) {
	case ast.NumberExpr:
		return MKNUM(e.Value)
	case ast.StringExpr:
		return MKSTR(e.Value)
	case ast.AssignmentExpr:
		return eval_assignment_expr(e, env)
	case ast.BinaryExpr:
		return eval_binary_expr(e, env)
	case ast.SymbolExpr:
		return eval_symbol_expr(e, env)
	case ast.CallExpr:
		return eval_call_expr(e, env)
	case ast.StructInstantiationExpr:
		return eval_struct_inst_expt(e, env)
	case ast.MemberAccessExpr:
		return eval_member_access_expr(e, env)
	case ast.ArrayInstantiationExpr:
		return eval_array_inst_expr(e, env)
	case ast.ArrayAccessExpr:
		return eval_array_access_expr(e, env)
	case ast.PrefixExpr:
		return eval_prefix_expr(e, env)

	default:
		panic(fmt.Sprintf("This Expr Node has not yet been set up for interpretation.\n %s", e))
	}

}

func eval_assignment_expr(expr ast.AssignmentExpr, env *environment) RuntimeVal {
	switch e := expr.Assigne.(type) {
	case ast.SymbolExpr:

		var val RuntimeVal
		if expr.Operator.Kind == lexer.PLUS_EQUALS {
			val = eval_binary_expr(ast.BinaryExpr{Left: expr.Assigne, Right: expr.Value, Operator: lexer.NewToken(lexer.PLUS, "+")}, env)
		} else if expr.Operator.Kind == lexer.MINUS_EQUALS {
			val = eval_binary_expr(ast.BinaryExpr{Left: expr.Assigne, Right: expr.Value, Operator: lexer.NewToken(lexer.PLUS, "-")}, env)
		} else {
			val = Evaluate(ast.ExpressionStmt{Expression: expr.Value}, env)
		}
		return env.assign_var(e.Value, val)
	case ast.MemberAccessExpr:

		val := Evaluate(ast.ExpressionStmt{Expression: expr.Value}, env)

		baseVarName, nestedMembers := getBaseVariableAndNestedMembers(e)

		baseVar := env.lookup_var(baseVarName).(Variable)
		if baseVar.Value.Type() != StructType {
			panic(fmt.Sprintf("Variable %s is not a struct.", baseVarName))
		}

		assignNestedMember(baseVar.Value.(StructInstance).Properties, nestedMembers, val)

		return val

	default:
		panic(fmt.Sprintf("Invalid LHS inside assignment expr: %s", e))
	}
}

func eval_binary_expr(binop ast.BinaryExpr, env *environment) RuntimeVal {
	lhs := eval_expr(binop.Left, env)  // Evaluate the left-hand side expression
	rhs := eval_expr(binop.Right, env) // Evaluate the right-hand side expression

	if lhs.Type() == "number" && rhs.Type() == "number" {

		return eval_numeric_binary_expr(lhs.(NumberVal), rhs.(NumberVal), binop.Operator.Value)
	}
	return eval_logical_binary_expr(lhs, rhs, binop.Operator.Value)

}

func eval_symbol_expr(expr ast.SymbolExpr, env *environment) RuntimeVal {
	return env.lookup_var(expr.Value).(Variable).Value
}

func eval_numeric_binary_expr(lhs NumberVal, rhs NumberVal, operator string) RuntimeVal {
	var res float64

	switch operator {
	case "+":
		res = lhs.Value + rhs.Value
	case "-":
		res = lhs.Value - rhs.Value
	case "*":
		res = lhs.Value * rhs.Value
	case "/":
		// Check for division by zero
		if rhs.Value != 0 {
			res = (lhs.Value / rhs.Value)
		} else {
			fmt.Println("Error: Division by zero")
			return MKNULL()
		}
	case "%":
		// Check if the right value is not zero
		if rhs.Value != 0 {
			// Modulus operation only supports integers, so convert to integers
			res = float64(int(lhs.Value) % int(rhs.Value))
		} else {
			fmt.Println("Error: Modulus by zero")
			return MKNULL()
		}
	default:

		return eval_comparison_binary_expr(lhs, rhs, operator)
	}

	return MKNUM(res)
}

func eval_logical_binary_expr(lhs RuntimeVal, rhs RuntimeVal, operator string) RuntimeVal {
	var res bool

	switch operator {
	case "&&":
		res = truthify(lhs) && truthify(rhs)
	case "||":
		res = truthify(lhs) || truthify(rhs)
	case "==":
		res = equals(lhs, rhs)
	case "!=":
		res = !equals(lhs, rhs)
	default:
		fmt.Println("Error: Unsupported logical operator:", operator, "for types:", lhs.Type(), "-", rhs.Type())
		return MKNULL()
	}

	return MKBOOL(res)
}

func eval_comparison_binary_expr(lhs RuntimeVal, rhs RuntimeVal, operator string) RuntimeVal {
	lhsNum, lhsIsNum := lhs.(NumberVal)
	rhsNum, rhsIsNum := rhs.(NumberVal)

	if lhsIsNum && rhsIsNum {
		var res bool

		switch operator {
		case ">":
			res = lhsNum.Value > rhsNum.Value
		case "<":
			res = lhsNum.Value < rhsNum.Value
		case ">=":
			res = lhsNum.Value >= rhsNum.Value
		case "<=":
			res = lhsNum.Value <= rhsNum.Value
		case "==":
			res = equals(lhs, rhs)
		case "!=":
			res = !equals(lhs, rhs)
		default:
			fmt.Println("Error: Unsupported operator:", operator)
			return MKNULL()
		}

		return MKBOOL(res)
	}

	// Handle other types if needed (e.g., strings)
	fmt.Println("Error: Unsupported types for comparison")
	return MKNULL()
}

func eval_call_expr(call ast.CallExpr, env *environment) RuntimeVal {

	switch fnVal := env.lookup_var(call.FunctionName).(Variable).Value.(type) {
	case Function:

		if len(call.Arguments) != len(fnVal.Parameters) {

			if len(call.Arguments) < len(fnVal.Parameters) {
				panic(fmt.Sprintf("Missing arguments for function %s", call.FunctionName))
			} else {
				panic(fmt.Sprintf("Too many arguments for function %s", call.FunctionName))
			}
		}

		callEnv := &environment{Variables: map[string]Variable{}, Parent: env}
		for i, param := range fnVal.Parameters {
			callEnv.declare_var(param.ParamName, eval_expr(call.Arguments[i], callEnv), param.ParamType, false)
		}

		return Evaluate(fnVal.Body, callEnv)

	case NativeFn:
		eval_args := make([]RuntimeVal, 0)
		for _, elem := range call.Arguments {
			eval_args = append(eval_args, eval_expr(elem, env))
		}
		return fnVal.Call(eval_args)

	default:
		panic(fmt.Sprintf("%s is not callable", call.FunctionName))
	}

}

func eval_struct_inst_expt(s ast.StructInstantiationExpr, env *environment) RuntimeVal {
	structDef, ok := env.lookup_struct(s.StructName)

	if !ok {
		panic(fmt.Sprintf("Struct %s not defined", s.StructName))
	}

	evalProps := make(map[string]RuntimeVal)

	for _, structProp := range structDef.(Struct).Properties {
		prop, ok := s.Properties[structProp.PropName]
		if ok {
			propVal := eval_expr(prop, env)

			expectedType := structProp.PropType

			if expectedType != propVal.Type() && expectedType != AnyType && propVal.Type() != StructType {
				panic(fmt.Sprintf("Type mismatch for property struct %s of %s : expected %s, got %s", structProp.PropName, structDef.(Struct).Name, expectedType, propVal.Type()))
			} else if propVal.Type() == StructType && expectedType != AnyType {
				if propVal.(StructInstance).StructName != string(expectedType) {
					panic(fmt.Sprintf("Property %s of struct %s is of type %s, got %s", structProp.PropName, structDef.(Struct).Name, expectedType, propVal.(StructInstance).StructName))
				}
			}

			evalProps[structProp.PropName] = propVal
		} else {
			evalProps[structProp.PropName] = MKNULL()
		}

	}

	return StructInstance{
		StructName: s.StructName,
		Properties: evalProps,
	}
}

func eval_member_access_expr(ma ast.MemberAccessExpr, env *environment) RuntimeVal {

	structVal := eval_expr(ma.Struct, env)

	// Check if the evaluated value is a struct instance
	structInstance, ok := structVal.(StructInstance)

	if !ok {
		panic(fmt.Sprintf("%v is not a struct instance", ma.Struct))
	}

	memberVal, exists := structInstance.Properties[ma.Member]
	if !exists {
		panic(fmt.Sprintf("Member %s not found in struct %s", ma.Member, structInstance.StructName))
	}

	return memberVal
}

func eval_array_inst_expr(aie ast.ArrayInstantiationExpr, env *environment) RuntimeVal {
	var elements []ArrayElement
	declaredType := extractValueType(aie.Underlying)

	// Evaluate each element expression and append to the elements array
	for _, elemExpr := range aie.Contents {
		elemVal := eval_expr(elemExpr, env)

		// Check if the element's type matches the declared array type
		if elemVal.Type() != declaredType {
			panic(fmt.Sprintf("Element type mismatch: expected %s, got %s", declaredType, elemVal.Type()))
		}

		elements = append(elements, ArrayElement{ElementType: elemVal.Type(), Value: elemVal})
	}

	// Create and return the array instance
	return Array{Elements: elements, ValType: declaredType}
}

func eval_array_access_expr(aa ast.ArrayAccessExpr, env *environment) RuntimeVal {
	array := eval_expr(aa.Array, env)
	index := eval_expr(aa.Index, env)

	if index.Type() != NumberType {
		panic("Index must be of type number.")
	}

	if array.Type() == StringType {
		if aa.Rest {
			return MKSTR(array.(StringVal).Value[int(index.(NumberVal).Value):])
		}
		if aa.Prev {
			return MKSTR(array.(StringVal).Value[:int(index.(NumberVal).Value)])
		}
		return MKSTR(string(array.(StringVal).Value[int(index.(NumberVal).Value)]))
	} else if string(array.Type())[:5] != "array" {
		panic(fmt.Sprintf("Out of index expected %s got %s", ArrayType, array.Type()))
	}

	if aa.Prev {
		return Array{Elements: array.(Array).Elements[:int(index.(NumberVal).Value)], ValType: array.(Array).ValType}
	}
	if aa.Rest {
		return Array{Elements: array.(Array).Elements[int(index.(NumberVal).Value):], ValType: array.(Array).ValType}
	}

	return array.(Array).Elements[int(index.(NumberVal).Value)].Value
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
		panic(fmt.Sprintf("unexpected prefix operator: %s", lexer.TokenKindString(pr.Operator.Kind)))
	}
}
