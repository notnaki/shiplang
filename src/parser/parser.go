package parser

import (
	"fmt"
	"shiplang/src/ast"
	"shiplang/src/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func createParser(tokens []lexer.Token) *parser {
	createTokenTypeLookups()
	createTokenLookups()
	return &parser{
		tokens: tokens,
		pos:    0,
	}
}

func Parse(tokens []lexer.Token) ast.BlockStmt {

	Body := make([]ast.Stmt, 0)

	p := createParser(tokens)

	for p.hasTokens() {
		Body = append(Body, parse_stmt(p))
	}

	return ast.BlockStmt{
		Body: Body,
	}
}

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) currentTokenKind() lexer.TokenKind {
	return p.currentToken().Kind
}

func (p *parser) advance() lexer.Token {
	tk := p.currentToken()
	p.pos++
	return tk
}

func (p *parser) hasTokens() bool {
	return p.pos < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

func (p *parser) expectError(expected lexer.TokenKind, err any) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	if kind != expected {
		if err == nil {
			err = fmt.Sprintf("Expected %s but recieved %s instead", lexer.TokenKindString(expected), lexer.TokenKindString(kind))
		}
		panic(err)
	}

	return p.advance()
}

func (p *parser) expect(expected lexer.TokenKind) lexer.Token {
	return p.expectError(expected, nil)
}
