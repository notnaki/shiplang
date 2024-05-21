package interpreter

import (
	"fmt"
	"shipgo/src/ast"
)

func Evaluate(node ast.Stmt, env *environment) RuntimeVal {

	switch n := node.(type) {
	case ast.ExpressionStmt:
		return eval_expr(n.Expression, env)

	case ast.BlockStmt:
		return eval_block_stmt(n, env)

	case ast.VarDeclStmt:
		return eval_var_decl_stmt(n, env)

	case ast.FnDeclStmt:
		return eval_fn_decl_stmt(n, env)

	case ast.ReturnStmt:
		return eval_return_stmt(n, env)

	case ast.StructDeclStmt:
		return eval_struct_decl_stmt(n, env)
	case ast.IfStmt:
		return eval_if_stmt(n, env)
	case ast.WhileStmt:
		return eval_while_stmt(n, env)
	case ast.ForeachStmt:
		return eval_foreach_stmt(n, env)
	case ast.ForStmt:
		return eval_for_stmt(n, env)
	case ast.ImportStmt:
		return eval_import_stmt(n, env)

	default:
		panic(fmt.Sprintf("This Stmt Node has not yet been set up for interpretation.\n %s", node))
	}
}
