package parser

import (
	"fmt"
	"shipgo/src/ast"
	"shipgo/src/helpers"
	"shipgo/src/lexer"
	"strconv"
)

func parse_expr(p *parser, bp binding_power) ast.Expr {
	tkind := p.currentTokenKind()

	nud_fn, exists := nud_lu[tkind]

	if !exists {
		panic(fmt.Sprintf("Not impl nud handler for %s", lexer.TokenKindString(tkind)))
	}

	left := nud_fn(p)

	for bp_lu[p.currentTokenKind()] > bp {
		tkind = p.currentTokenKind()
		led_fn, exists := led_lu[tkind]

		if !exists {
			panic("Not impl led handler")
		}

		left = led_fn(p, left, bp_lu[p.currentTokenKind()])
	}

	return left
}

func parse_primary_expr(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumberExpr{Value: number}

	case lexer.STRING:
		return ast.StringExpr{Value: p.advance().Value}

	case lexer.IDENTIFIER:
		return ast.SymbolExpr{Value: p.advance().Value}

	default:
		panic(fmt.Sprintf("Cannot create primary expression from %s\n", lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parse_binary_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	optk := p.advance()
	right := parse_expr(p, bp)

	return ast.BinaryExpr{
		Left:     left,
		Operator: optk,
		Right:    right,
	}
}

func parse_assignment_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operatorToken := p.advance()

	switch operatorToken.Kind {
	case lexer.PLUS_PLUS:
		return ast.AssignmentExpr{
			Operator: lexer.NewToken(lexer.PLUS_EQUALS, "+="),
			Value:    ast.NumberExpr{Value: 1},
			Assigne:  left,
		}
	case lexer.MINUS_MINUS:
		return ast.AssignmentExpr{
			Operator: lexer.NewToken(lexer.MINUS_EQUALS, "-="),
			Value:    ast.NumberExpr{Value: 1},
			Assigne:  left,
		}
	default:
		rhs := parse_expr(p, bp)

		return ast.AssignmentExpr{
			Operator: operatorToken,
			Value:    rhs,
			Assigne:  left,
		}
	}
}

func parse_prefix_expr(p *parser) ast.Expr {
	operatorToken := p.advance()
	rhs := parse_expr(p, primary)

	return ast.PrefixExpr{
		Operator:  operatorToken,
		RightExpr: rhs,
	}
}

func parse_struct_instantiation_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {

	var structName = helpers.ExpectType[ast.SymbolExpr](left).Value
	var properties = map[string]ast.Expr{}
	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		var propertyName = p.expect(lexer.IDENTIFIER).Value
		p.expect(lexer.COLON)
		expr := parse_expr(p, logical)

		properties[propertyName] = expr

		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.StructInstantiationExpr{
		StructName: structName,
		Properties: properties,
	}
}

func parse_array_instantiation_expr(p *parser) ast.Expr {
	var underlyingType ast.Type
	var contents = []ast.Expr{}

	p.expect(lexer.OPEN_BRACKET)
	p.expect(lexer.CLOSE_BRACKET)

	underlyingType = parse_type(p, default_bp)

	p.expect(lexer.OPEN_CURLY)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		contents = append(contents, parse_expr(p, logical))

		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.ArrayInstantiationExpr{
		Underlying: underlyingType,
		Contents:   contents,
	}
}

func parse_member_access_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.expect(lexer.DOT)
	memberName := p.expect(lexer.IDENTIFIER).Value

	return ast.MemberAccessExpr{
		Struct: left,
		Member: memberName,
	}
}

func parse_array_access_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.expect(lexer.OPEN_BRACKET)
	var prev = false
	var rest = false

	if p.currentTokenKind() == lexer.CLOSE_BRACKET {
		panic("Index expected after open bracket.")
	}

	if p.currentTokenKind() == lexer.COLON {
		p.expect(lexer.COLON)
		prev = true

	}

	index := parse_expr(p, bp)

	if p.currentTokenKind() == lexer.COLON {
		p.expect(lexer.COLON)
		rest = true

	}

	p.expect(lexer.CLOSE_BRACKET)

	return ast.ArrayAccessExpr{
		Array: left,
		Index: index,
		Rest:  rest,
		Prev:  prev,
	}

}

func parse_grouping_expr(p *parser) ast.Expr {
	p.advance()
	expr := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)
	return expr
}

func parse_call_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {

	functionName := left.(ast.SymbolExpr).Value
	p.expect(lexer.OPEN_PAREN) // Consume '('

	// Parse arguments
	args := parse_call_params_list(p)

	p.expect(lexer.CLOSE_PAREN) // Consume ')'

	return ast.CallExpr{
		FunctionName: functionName,
		Arguments:    args,
	}
}

// parse_expr_list parses a list of expressions separated by commas
func parse_call_params_list(p *parser) []ast.Expr {
	exprs := make([]ast.Expr, 0)

	for p.currentTokenKind() != lexer.CLOSE_PAREN && p.hasTokens() {
		expr := parse_expr(p, default_bp)
		exprs = append(exprs, expr)

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expect(lexer.COMMA) // Consume ',' between expressions
		}
	}

	return exprs
}
