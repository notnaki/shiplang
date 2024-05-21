package interpreter

import (
	"fmt"
	"shipgo/src/ast"
)

func eval_block_stmt(s ast.BlockStmt, env *environment) RuntimeVal {
	last_evaluated := MKNULL()
	for _, s := range s.Body {
		// litter.Dump(s)
		last_evaluated = Evaluate(s, env)

		if returnVal, ok := last_evaluated.(Return); ok {
			return returnVal.Value
		}
	}

	return last_evaluated
}

func eval_var_decl_stmt(decl ast.VarDeclStmt, env *environment) RuntimeVal {
	var val RuntimeVal

	if decl.AssignedValue == nil {
		val = MKNULL()
	} else {
		val = eval_expr(decl.AssignedValue, env)
	}

	switch exp := decl.ExplicitType.(type) {
	case ast.SymbolType:
		env.declare_var(decl.VarName, val, ValueType(exp.Name), decl.IsConstant)
	case ast.ArrayType:

		env.declare_var(decl.VarName, val, extractValueType(exp), decl.IsConstant)
	default:
		env.declare_var(decl.VarName, val, AnyType, decl.IsConstant)
	}

	return val
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

func eval_return_stmt(stmt ast.ReturnStmt, env *environment) RuntimeVal {
	return Return{Value: eval_expr(stmt.Value, env)}
}

func eval_fn_decl_stmt(fn ast.FnDeclStmt, env *environment) RuntimeVal {

	params := make([]Parameter, len(fn.Parameters))
	for i, p := range fn.Parameters {
		params[i] = AstToRuntimeParams(p)
	}

	fnVal := Function{
		fn.FnName,
		params,
		fn.Body,
		env,
	}

	env.declare_var(fn.FnName, fnVal, FunctionType, true)

	return fnVal
}

func AstToRuntimeParams(astParam ast.Parameter) Parameter {
	return Parameter{
		ParamName: astParam.Name,
		ParamType: ValueType(astParam.Type.(ast.SymbolType).Name),
	}
}

func eval_struct_decl_stmt(s ast.StructDeclStmt, env *environment) RuntimeVal {

	props := make([]Property, 0, len(s.Properties))
	for n, p := range s.Properties {
		props = append(props, AstToRuntimeProps(p, n))
	}

	return env.declare_struct(s.StructName, props)
}

func AstToRuntimeProps(astProp ast.StructProperty, propName string) Property {
	return Property{PropName: propName, PropType: ValueType(astProp.Type.(ast.SymbolType).Name)}
}

func eval_if_stmt(i ast.IfStmt, env *environment) RuntimeVal {
	// Evaluate the condition of the if statement
	cond := truthify(eval_expr(i.Condition, env))

	if cond {
		// If the condition is true, evaluate the if body
		eval_block_stmt(i.IfBody, env)
	} else {
		// Otherwise, iterate through the elif bodies
		for elifCond, elifBody := range i.ElifBodies {
			// Evaluate each elif condition
			elifCondVal := truthify(eval_expr(elifCond, env))
			if elifCondVal {
				// If an elif condition is true, evaluate its body and return
				eval_block_stmt(elifBody, env)
				return MKNULL()
			}
		}
		// If none of the elif conditions were true, evaluate the else body
		eval_block_stmt(i.ElseBody, env)
	}
	return MKNULL()
}

func eval_while_stmt(w ast.WhileStmt, env *environment) RuntimeVal {
	for {
		cond := truthify(eval_expr(w.Condition, env))
		if !cond {
			break
		}
		eval_block_stmt(w.Body, env)
	}
	return MKNULL()
}

func eval_foreach_stmt(f ast.ForeachStmt, env *environment) RuntimeVal {

	collectionVal := eval_expr(f.Collection, env)
	loopEnv := &environment{Variables: make(map[string]Variable), Parent: env}

	loopEnv.declare_var(f.Iterator, MKNULL(), AnyType, false)

	switch collectionVal := collectionVal.(type) {
	case Array:
		for _, item := range collectionVal.Elements {

			loopEnv.assign_var(f.Iterator, item.Value)

			// Evaluate the body of the loop
			eval_block_stmt(f.Body, loopEnv)
		}
	default:

		panic("Invalid foreach loop collection type")
	}

	return MKNULL()
}

func eval_for_stmt(f ast.ForStmt, env *environment) RuntimeVal {
	// Evaluate initialization statement
	loopEnv := &environment{Variables: make(map[string]Variable), Parent: env}
	if f.Init != nil {
		Evaluate(f.Init, loopEnv)
	}

	for {

		if f.Cond != nil {
			cond := eval_expr(f.Cond, loopEnv)
			if !truthify(cond) {
				break
			}
		}

		eval_block_stmt(f.Body, loopEnv)

		if f.Post != nil {
			Evaluate(f.Post, loopEnv)
		}
	}

	return MKNULL()
}
