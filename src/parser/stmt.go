package parser

import (
	"fmt"
	"shipgo/src/ast"
	"shipgo/src/lexer"
)

func parse_stmt(p *parser) ast.Stmt {
	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	expr := parse_expr(p, default_bp)
	p.expect(lexer.SEMI_COLON)

	return ast.ExpressionStmt{Expression: expr}
}

func parse_var_decl_stmt(p *parser) ast.Stmt {
	var explicitType ast.Type
	var assignedVal ast.Expr

	isConst := p.advance().Kind == lexer.CONST
	varName := p.expectError(lexer.IDENTIFIER, "Inside variable declaration expected to find variable name").Value

	if p.currentTokenKind() == lexer.COLON {
		p.advance()
		explicitType = parse_type(p, default_bp)
	}

	if p.currentTokenKind() != lexer.SEMI_COLON {
		p.expect(lexer.ASSIGNMENT)
		assignedVal = parse_expr(p, assignment)
	} else if explicitType == nil {
		panic("Missing either right-hand side in var declaration or explicit type.")
	}

	p.expect(lexer.SEMI_COLON)

	if isConst && assignedVal == nil {
		panic("Cannot define constant type without providing value")
	}

	return ast.VarDeclStmt{
		ExplicitType:  explicitType,
		IsConstant:    isConst,
		VarName:       varName,
		AssignedValue: assignedVal,
	}
}

func parse_if_stmt(p *parser) ast.Stmt {
	p.expect(lexer.IF)
	p.expect(lexer.OPEN_PAREN)
	condition := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)
	body := parse_block_stmt(p)

	elifBodies := make(map[ast.Expr]ast.BlockStmt)

	var elseBody ast.BlockStmt

	for p.currentTokenKind() == lexer.ELSE {
		p.expect(lexer.ELSE)

		if p.currentTokenKind() == lexer.IF {
			p.expect(lexer.IF)
			p.expect(lexer.OPEN_PAREN)
			elifCondition := parse_expr(p, default_bp)
			p.expect(lexer.CLOSE_PAREN)
			elifBody := parse_block_stmt(p)
			elifBodies[elifCondition] = elifBody.(ast.BlockStmt)
		} else {
			// Parse else
			elseBody = parse_block_stmt(p).(ast.BlockStmt)
			break // Exit loop after parsing else
		}
	}

	return ast.IfStmt{
		IfBody:     body.(ast.BlockStmt),
		Condition:  condition,
		ElifBodies: elifBodies,
		ElseBody:   elseBody,
	}
}

func parse_struct_decl_stmt(p *parser) ast.Stmt {

	p.expect(lexer.STRUCT)
	var properties = map[string]ast.StructProperty{}
	var methods = map[string]ast.StructMethod{}
	var structName = p.expect(lexer.IDENTIFIER).Value

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		// var IsStatic bool
		var propertyName string
		// if p.currentTokenKind() == lexer.STATIC {
		// 	IsStatic = true
		// 	p.expect(lexer.STATIC)
		// }

		if p.currentTokenKind() == lexer.IDENTIFIER {
			propertyName = p.expect(lexer.IDENTIFIER).Value
			p.expectError(lexer.COLON, "Expected to find colon following property name inside struct declaration")

			structType := parse_type(p, default_bp)
			p.expect(lexer.SEMI_COLON)

			_, exists := properties[propertyName]

			if exists {
				panic(fmt.Sprintf("Property %s has already been defined in struct declaration", propertyName))
			}

			properties[propertyName] = ast.StructProperty{
				// IsStatic: IsStatic,
				Type: structType,
			}

			continue
		}

		panic("Cannot currently handle methods inside struct decl")
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.StructDeclStmt{
		StructName: structName,
		Properties: properties,
		Methods:    methods,
	}
}

// fn hello(){}
// parse_fn_decl_stmt parses a function declaration statement
func parse_fn_decl_stmt(p *parser) ast.Stmt {
	p.expect(lexer.FN)

	fnName := p.expect(lexer.IDENTIFIER).Value

	p.expect(lexer.OPEN_PAREN)
	parameters := parse_fn_params(p)
	p.expect(lexer.CLOSE_PAREN)

	body := parse_block_stmt(p)

	return ast.FnDeclStmt{
		FnName:     fnName,
		Parameters: parameters,
		Body:       body.(ast.BlockStmt),
	}
}

// parse_fn_params parses function parameters
func parse_fn_params(p *parser) []ast.Parameter {
	params := make([]ast.Parameter, 0)

	for p.currentTokenKind() != lexer.CLOSE_PAREN && p.hasTokens() {
		paramName := p.expect(lexer.IDENTIFIER).Value
		var pType ast.Type = ast.SymbolType{Name: "any"}

		if p.currentTokenKind() == lexer.COLON {
			p.advance()
			pType = parse_type(p, default_bp)
		}

		params = append(params, ast.Parameter{Name: paramName, Type: pType})

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expect(lexer.COMMA) // Consume ',' between parameters
		}
	}

	return params
}

func parse_block_stmt(p *parser) ast.Stmt {
	p.expect(lexer.OPEN_CURLY)
	body := make([]ast.Stmt, 0)

	for p.currentTokenKind() != lexer.CLOSE_CURLY && p.hasTokens() {
		stmt := parse_stmt(p)
		body = append(body, stmt)
	}

	p.expect(lexer.CLOSE_CURLY)
	return ast.BlockStmt{Body: body}
}

func parse_return_stmt(p *parser) ast.Stmt {
	p.advance() // eat the return token
	returnval := parse_expr(p, assignment)
	p.expect(lexer.SEMI_COLON)

	return ast.ReturnStmt{
		Value: returnval,
	}

}

func parse_while_stmt(p *parser) ast.Stmt {
	p.expect(lexer.WHILE)
	p.expect(lexer.OPEN_PAREN)
	cond := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)

	body := parse_block_stmt(p).(ast.BlockStmt)

	return ast.WhileStmt{
		Condition: cond,
		Body:      body,
	}
}

func parse_foreach_stmt(p *parser) ast.Stmt {
	p.expect(lexer.FOREACH)
	p.expect(lexer.OPEN_PAREN)
	iterator := p.expect(lexer.IDENTIFIER).Value
	p.expect(lexer.IN)
	collection := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)
	body := parse_block_stmt(p).(ast.BlockStmt)

	return ast.ForeachStmt{
		Iterator:   iterator,
		Collection: collection,
		Body:       body,
	}
}

func parse_for_stmt(p *parser) ast.Stmt {
	p.expect(lexer.FOR)
	p.expect(lexer.OPEN_PAREN)

	var init ast.Stmt
	if p.currentTokenKind() != lexer.SEMI_COLON {
		init = parse_stmt(p)
	}

	var cond ast.Expr
	if p.currentTokenKind() != lexer.SEMI_COLON {
		cond = parse_expr(p, default_bp)
	} else {
		panic("Expected condition expression in for loop")
	}

	p.expect(lexer.SEMI_COLON)

	var post ast.Stmt
	if p.currentTokenKind() != lexer.CLOSE_PAREN {
		post = parse_stmt(p)
	}

	p.expect(lexer.CLOSE_PAREN)

	body := parse_block_stmt(p).(ast.BlockStmt)

	return ast.ForStmt{
		Init: init,
		Cond: cond,
		Post: post,
		Body: body,
	}
}
