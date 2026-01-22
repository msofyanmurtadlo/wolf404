package parser

import (
	"wolf404/compiler/ast"
	"wolf404/compiler/lexer"
)

func (p *Parser) parseVariableUsage() ast.Expression {
	// Current token is $
	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}
	// Now current token is the variable name (IDENT)
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekToken.Type == lexer.TOKEN_IDENT {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if p.peekToken.Type == lexer.TOKEN_NEWLINE {
		p.nextToken()
	}
	
	// Ensure we are at INDENT if it exists
	if p.peekToken.Type == lexer.TOKEN_INDENT {
		p.nextToken()
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekToken.Type == lexer.TOKEN_RPAREN {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	// Logic: Parameters must start with $
	if p.curToken.Type == lexer.TOKEN_DOLLAR {
		if !p.expectPeek(lexer.TOKEN_IDENT) {
			return nil
		}
	} else if p.curToken.Type != lexer.TOKEN_IDENT {
		// Allow IDENT as well? No, enforce $ for strictness or flexibility?
		// User writes hunt($a). curToken is $.
		// If user writes hunt(a), curToken is a.
		// My parser_extensions_v3 prev logic assumed IDENT.
		// Let's support both or strictly $. Given Wolf404 is $ var, strictly $.
		// BUT wait, p.nextToken above advanced to first token.
		// If first token is $, we need to advance to ident.
		// If first token is ident ...
		return nil // Should be $
	}

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekToken.Type == lexer.TOKEN_COMMA {
		p.nextToken() // consume comma
		p.nextToken() // consume next param start ($)
		
		if p.curToken.Type == lexer.TOKEN_DOLLAR {
			if !p.expectPeek(lexer.TOKEN_IDENT) {
				return nil
			}
		}
		
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(lexer.TOKEN_RPAREN) {
		return nil
	}

	return identifiers
}
