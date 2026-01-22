package parser

import (
	"wolf404/compiler/ast"
	"wolf404/compiler/lexer"
)

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken() // consume 'sniff'
	expression.Condition = p.parseExpression(LOWEST)

	if p.peekToken.Type == lexer.TOKEN_NEWLINE {
		p.nextToken()
	}

	if p.peekToken.Type == lexer.TOKEN_INDENT {
		p.nextToken()
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == lexer.TOKEN_MISSING {
		p.nextToken() // consume 'missing'
		if p.peekToken.Type == lexer.TOKEN_NEWLINE {
			p.nextToken()
		}
		if p.peekToken.Type == lexer.TOKEN_INDENT {
			p.nextToken()
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	// Expect INDENT
	if p.curToken.Type == lexer.TOKEN_INDENT {
		p.nextToken() // consume INDENT
		
		for p.curToken.Type != lexer.TOKEN_DEDENT && p.curToken.Type != lexer.TOKEN_EOF {
			if p.curToken.Type == lexer.TOKEN_NEWLINE {
				p.nextToken()
				continue
			}
			stmt := p.parseStatement()
			if stmt != nil {
				block.Statements = append(block.Statements, stmt)
			}
			p.nextToken()
		}
		
		if p.curToken.Type == lexer.TOKEN_DEDENT {
			p.nextToken() // consume DEDENT
		}
	} else {
		// Fallback for single line or no-indent (MVP compliance)
		// parse single statement
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken() // consume the statement end (?)
	}

	return block
}

func (p *Parser) parseHowlExpression() ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: &ast.Identifier{Token: p.curToken, Value: "howl"}}
	
	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseProwlStatement() *ast.ProwlStatement {
	stmt := &ast.ProwlStatement{Token: p.curToken}
	
	p.nextToken() // consume 'prowl'
	
	stmt.Call = p.parseExpression(LOWEST)
	
	return stmt
}
