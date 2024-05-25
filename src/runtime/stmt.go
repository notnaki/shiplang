package runtime

import (
	"shiplang/src/ast"
)

func eval_block_stmt(s ast.BlockStmt, env *environment) RuntimeVal {
	last_evaluated := MKNULL()

	for _, s := range s.Body {
		last_evaluated = Evaluate(s, env)

		var leType ValueType = last_evaluated.Type()

		if leType == BreakType {
			return Break{}
		}

		if leType == ReturnType {
			return last_evaluated.(Return).Value
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

	switch expTyp := decl.ExplicitType.(type) {
	case ast.SymbolType:
		env.declareVar(decl.VarName, val, ValueType(expTyp.Name), decl.IsConstant)
	case ast.ArrayType:
		env.declareVar(decl.VarName, val, extractValueType(expTyp), decl.IsConstant)
	default:
		env.declareVar(decl.VarName, val, AnyType, decl.IsConstant)
	}

	return val
}

func eval_struct_decl_stmt(decl ast.StructDeclStmt, env *environment) RuntimeVal {
	props := make(map[string]ValueType)

	for name, prop := range decl.Properties {
		switch p := prop.Type.(type) {
		case ast.SymbolType:
			props[name] = ValueType(p.Name)
		case ast.ArrayType:
			props[name] = ValueType(extractValueType(p))

		}
	}

	return env.declareStruct(decl.StructName, props)
}

func eval_struct_impl_stmt(impl ast.ImplStmt, env *environment) RuntimeVal {
	m := eval_impl_fn(impl.Method, env)
	method, isFunc := m.(Function)

	if isFunc {
		return env.implMethod(impl.Struct, method)
	}
	panic("")
}

func eval_impl_fn(decl ast.FnDeclStmt, env *environment) RuntimeVal {
	params := make([]Parameter, len(decl.Parameters))

	for i, param := range decl.Parameters {
		switch p := param.Type.(type) {
		case ast.SymbolType:
			params[i] = Parameter{Name: param.Name, Type: ValueType(p.Name)}
		case ast.ArrayType:
			params[i] = Parameter{Name: param.Name, Type: ValueType(extractValueType(p))}

		}
	}

	fn := Function{
		Name:       decl.FnName,
		Parameters: params,
		Body:       decl.Body,
		Env:        env,
	}

	return fn
}

func eval_fn_decl_stmt(decl ast.FnDeclStmt, env *environment) RuntimeVal {
	params := make([]Parameter, len(decl.Parameters))

	for i, param := range decl.Parameters {
		switch p := param.Type.(type) {
		case ast.SymbolType:
			params[i] = Parameter{Name: param.Name, Type: ValueType(p.Name)}
		case ast.ArrayType:
			params[i] = Parameter{Name: param.Name, Type: ValueType(extractValueType(p))}

		}
	}

	fn := Function{
		Name:       decl.FnName,
		Parameters: params,
		Body:       decl.Body,
		Env:        env,
	}

	env.declareFn(fn)

	return fn
}

func eval_return_stmt(r ast.ReturnStmt, env *environment) RuntimeVal {
	value := eval_expr(r.Value, env)
	return Return{Value: value}
}

func eval_if_stmt(i ast.IfStmt, env *environment) RuntimeVal {
	condition := truthify(eval_expr(i.Condition, env))

	if condition {
		return eval_block_stmt(i.IfBody, env)
	}

	for cond, body := range i.ElifBodies {
		if truthify(eval_expr(cond, env)) {

			return eval_block_stmt(body, env)
		}
	}

	return eval_block_stmt(i.ElseBody, env)

}

func eval_while_stmt(w ast.WhileStmt, env *environment) RuntimeVal {
	for {
		if !truthify(eval_expr(w.Condition, env)) {
			break
		}
		value := eval_block_stmt(w.Body, env)
		if value.Type() == BreakType {
			break
		}
	}

	return MKNULL()
}

func eval_for_stmt(f ast.ForStmt, env *environment) RuntimeVal {
	loopEnv := &environment{Variables: make(map[string]Variable), Parent: env}

	if f.Init == nil || f.Cond == nil || f.Post == nil {
		panic("")
	}

	Evaluate(f.Init, loopEnv)
	for {
		cond := eval_expr(f.Cond, loopEnv)
		if !truthify(cond) {
			break
		}

		val := eval_block_stmt(f.Body, loopEnv)

		if val.Type() == BreakType {
			break
		}

		Evaluate(f.Post, loopEnv)
	}

	return MKNULL()
}

func eval_foreach_stmt(fe ast.ForeachStmt, env *environment) RuntimeVal {

	collection := eval_expr(fe.Collection, env)
	loopEnv := &environment{Variables: make(map[string]Variable), Parent: env}

	loopEnv.declareVar(fe.Iterator, MKNULL(), AnyType, false)

	switch collection := collection.(type) {
	case Array:
		for _, item := range collection.Elements {
			loopEnv.assignVar(fe.Iterator, item)
			val := eval_block_stmt(fe.Body, loopEnv)
			if val.Type() == BreakType {
				break
			}
		}
	case String:
		for _, char := range collection.Value {
			loopEnv.assignVar(fe.Iterator, MKSTR(string(char)))
			val := eval_block_stmt(fe.Body, loopEnv)
			if val.Type() == BreakType {
				break
			}
		}
	default:
		panic("")
	}

	return MKNULL()
}
