package parser

import (
	"wolf404/compiler/ast"
	"wolf404/compiler/lexer"
)

func (p *Parser) parseTrackStatement() *ast.TrackStatement {
	stmt := &ast.TrackStatement{Token: p.curToken}

	p.nextToken() // consume 'track'

	stmt.Condition = p.parseExpression(LOWEST)

	if p.peekToken.Type == lexer.TOKEN_NEWLINE {
		p.nextToken()
	}

	// Handle indentation (Python style block)
	if p.peekToken.Type == lexer.TOKEN_INDENT {
		p.nextToken()
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseClassStatement() *ast.ClassStatement {
	stmt := &ast.ClassStatement{Token: p.curToken}

	p.nextToken() // consume 'mold'
	
	// Expect Identifier (Class Name)
	if p.curToken.Type != lexer.TOKEN_IDENT {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	
	p.nextToken()

	// Parse Inheritance (mold Child : Parent)
	if p.curToken.Type == lexer.TOKEN_COLON {
		p.nextToken() // consume :
		if p.curToken.Type != lexer.TOKEN_IDENT {
			// Should add error here
			return nil
		}
		stmt.SuperClass = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
	}

	if p.peekToken.Type == lexer.TOKEN_NEWLINE {
		p.nextToken()
	}

	if p.peekToken.Type == lexer.TOKEN_INDENT {
		p.nextToken()
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}
