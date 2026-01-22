package parser

import (
	"wolf404/compiler/ast"
	"wolf404/compiler/lexer"
)

// ... existing code ...

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(lexer.TOKEN_RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekToken.Type == end {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekToken.Type == lexer.TOKEN_COMMA {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for p.peekToken.Type != lexer.TOKEN_RBRACE {
		if p.peekToken.Type == lexer.TOKEN_NEWLINE || p.peekToken.Type == lexer.TOKEN_INDENT {
			p.nextToken()
			continue
		}
		
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.TOKEN_COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		// Skip newlines after value
		for p.peekToken.Type == lexer.TOKEN_NEWLINE || p.peekToken.Type == lexer.TOKEN_INDENT || p.peekToken.Type == lexer.TOKEN_DEDENT {
			p.nextToken()
		}

		if p.peekToken.Type != lexer.TOKEN_RBRACE && !p.expectPeek(lexer.TOKEN_COMMA) {
			return nil
		}
	}

	if !p.expectPeek(lexer.TOKEN_RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.TOKEN_RBRACKET) {
		return nil
	}

	return exp
}

// ... existing parseProwlStatement etc ... (I will append this to parser_extensions.go)
