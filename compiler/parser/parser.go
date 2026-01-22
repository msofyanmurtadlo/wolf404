package parser

import (
	"fmt"
	"strconv"
	"wolf404/compiler/ast"
	"wolf404/compiler/lexer"
)

// Precedence levels
const (
	_ int = iota
	LOWEST
	ASSIGN      // =
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[lexer.TokenType]int{
	lexer.TOKEN_ASSIGN:   ASSIGN,
	lexer.TOKEN_DOT:      INDEX,
	lexer.TOKEN_LBRACKET: INDEX,
	lexer.TOKEN_EQ:       EQUALS,
	lexer.TOKEN_NEQ:      EQUALS,
	lexer.TOKEN_LT:       LESSGREATER,
	lexer.TOKEN_GT:       LESSGREATER,
	lexer.TOKEN_LTE:      LESSGREATER,
	lexer.TOKEN_GTE:      LESSGREATER,
	lexer.TOKEN_PLUS:     SUM,
	lexer.TOKEN_MINUS:    SUM,
	lexer.TOKEN_SLASH:    PRODUCT,
	lexer.TOKEN_ASTERISK: PRODUCT,
	lexer.TOKEN_LPAREN:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.TOKEN_IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.TOKEN_INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.TOKEN_STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.TOKEN_TRUE, p.parseBoolean)
	p.registerPrefix(lexer.TOKEN_FALSE, p.parseBoolean)
	p.registerPrefix(lexer.TOKEN_LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.TOKEN_DOLLAR, p.parseVariableExpression)
	p.registerPrefix(lexer.TOKEN_LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.TOKEN_DOLLAR, p.parseVariableExpression)
	p.registerPrefix(lexer.TOKEN_SNIFF, p.parseIfExpression)
	p.registerPrefix(lexer.TOKEN_HOWL, p.parseHowlExpression)

	p.registerPrefix(lexer.TOKEN_LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(lexer.TOKEN_LBRACE, p.parseHashLiteral)
	p.registerPrefix(lexer.TOKEN_HUNT, p.parseFunctionLiteral)
	p.registerPrefix(lexer.TOKEN_DOLLAR, p.parseVariableUsage)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.TOKEN_PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_ASTERISK, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_NEQ, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_LT, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_GT, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_ASSIGN, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_DOT, p.parseInfixExpression)
	p.registerInfix(lexer.TOKEN_LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.TOKEN_LBRACKET, p.parseIndexExpression)


	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	
	// Skip extra newlines if they appear where we expect statement starts, 
	// but mostly we want to consume them as separators in loop.
	// Actually, let's filter them out inside ParseProgram main loop or here?
	// If we filter here, we lose line info somewhat, but easy for MVP.
	// But let's keeping them and handle in ParseProgram.
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != lexer.TOKEN_EOF {
		if p.curToken.Type == lexer.TOKEN_NEWLINE {
			p.nextToken()
			continue
		}
		
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.TOKEN_BRING:
		return p.parseReturnStatement()
	case lexer.TOKEN_SUMMON:
		return p.parseSummonStatement()
	case lexer.TOKEN_PROWL:
		return p.parseProwlStatement()
	case lexer.TOKEN_TRACK:
		return p.parseTrackStatement()
	case lexer.TOKEN_MOLD:
		return p.parseClassStatement()
	case lexer.TOKEN_NEWLINE, lexer.TOKEN_INDENT, lexer.TOKEN_DEDENT:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// Expecting identifier after $
	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}
	
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.TOKEN_ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseVariableExpression() ast.Expression {
	// For $var usage in expressions
	// Current token is $
	// Next should be ident
	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}


func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curToken.Type == lexer.TOKEN_TRUE}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.TOKEN_RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekToken.Type == lexer.TOKEN_RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == lexer.TOKEN_COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.TOKEN_RPAREN) {
		return nil
	}

	return args
}

// Helpers

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		lexer.TokenTypeString(t), lexer.TokenTypeString(p.peekToken.Type))
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", lexer.TokenTypeString(t))
	p.errors = append(p.errors, msg)
}
