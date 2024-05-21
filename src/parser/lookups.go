package parser

import (
	"shipgo/src/ast"
	"shipgo/src/lexer"
)

type binding_power int

const (
	default_bp binding_power = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

type stmt_handler func(p *parser) ast.Stmt
type nud_handler func(p *parser) ast.Expr
type led_handler func(p *parser, left ast.Expr, bp binding_power) ast.Expr

type stmt_lookup map[lexer.TokenKind]stmt_handler
type nud_lookup map[lexer.TokenKind]nud_handler
type led_lookup map[lexer.TokenKind]led_handler
type bp_lookup map[lexer.TokenKind]binding_power

var bp_lu = bp_lookup{}
var nud_lu = nud_lookup{}
var led_lu = led_lookup{}
var stmt_lu = stmt_lookup{}

func led(kind lexer.TokenKind, bp binding_power, led_fn led_handler) {
	bp_lu[kind] = bp
	led_lu[kind] = led_fn
}

func nud(kind lexer.TokenKind, bp binding_power, nud_fn nud_handler) {
	bp_lu[kind] = primary
	nud_lu[kind] = nud_fn
}

func stmt(kind lexer.TokenKind, bp binding_power, stmt_fn stmt_handler) {
	bp_lu[kind] = default_bp
	stmt_lu[kind] = stmt_fn
}

func createTokenLookups() {

	led(lexer.ASSIGNMENT, assignment, parse_assignment_expr)
	led(lexer.PLUS_EQUALS, assignment, parse_assignment_expr)
	led(lexer.MINUS_EQUALS, assignment, parse_assignment_expr)
	led(lexer.PLUS_PLUS, assignment, parse_assignment_expr)
	led(lexer.MINUS_MINUS, assignment, parse_assignment_expr)

	led(lexer.AND, logical, parse_binary_expr)
	led(lexer.OR, logical, parse_binary_expr)
	led(lexer.DOT_DOT, logical, parse_binary_expr)

	led(lexer.LESS, relational, parse_binary_expr)
	led(lexer.LESS_EQUALS, relational, parse_binary_expr)
	led(lexer.GREATER, relational, parse_binary_expr)
	led(lexer.GREATER_EQUALS, relational, parse_binary_expr)
	led(lexer.EQUALS, relational, parse_binary_expr)
	led(lexer.NOT_EQUALS, relational, parse_binary_expr)

	led(lexer.PLUS, additive, parse_binary_expr)
	led(lexer.DASH, additive, parse_binary_expr)

	led(lexer.STAR, multiplicative, parse_binary_expr)
	led(lexer.SLASH, multiplicative, parse_binary_expr)
	led(lexer.PERCENT, multiplicative, parse_binary_expr)

	nud(lexer.DASH, unary, parse_prefix_expr)
	nud(lexer.NOT, unary, parse_prefix_expr)

	led(lexer.OPEN_CURLY, call, parse_struct_instantiation_expr)
	led(lexer.OPEN_PAREN, call, parse_call_expr)

	led(lexer.DOT, member, parse_member_access_expr)
	led(lexer.OPEN_BRACKET, member, parse_array_access_expr)

	nud(lexer.OPEN_BRACKET, primary, parse_array_instantiation_expr)
	nud(lexer.OPEN_PAREN, primary, parse_grouping_expr)
	nud(lexer.NUMBER, primary, parse_primary_expr)
	nud(lexer.STRING, primary, parse_primary_expr)
	nud(lexer.IDENTIFIER, primary, parse_primary_expr)

	stmt(lexer.IMPORT, default_bp, parse_import_stmt)
	stmt(lexer.CONST, default_bp, parse_var_decl_stmt)
	stmt(lexer.LET, default_bp, parse_var_decl_stmt)
	stmt(lexer.STRUCT, default_bp, parse_struct_decl_stmt)
	stmt(lexer.FN, default_bp, parse_fn_decl_stmt)
	stmt(lexer.RETURN, default_bp, parse_return_stmt)
	stmt(lexer.IF, default_bp, parse_if_stmt)
	stmt(lexer.WHILE, default_bp, parse_while_stmt)
	stmt(lexer.FOREACH, default_bp, parse_foreach_stmt)
	stmt(lexer.FOR, default_bp, parse_for_stmt)

}
