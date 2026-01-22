package lexer

import (
	"unicode"
)

// TokenType represents the type of token
type TokenType int

const (
	// Special tokens
	TOKEN_EOF TokenType = iota
	TOKEN_ILLEGAL

	// Identifiers and literals
	TOKEN_IDENT  // variable names, function names
	TOKEN_INT    // 123
	TOKEN_FLOAT  // 123.45
	TOKEN_STRING // "hello"
	TOKEN_TRUE   // true
	TOKEN_FALSE  // false

	// Operators
	TOKEN_ASSIGN   // =
	TOKEN_PLUS     // +
	TOKEN_MINUS    // -
	TOKEN_ASTERISK // *
	TOKEN_SLASH    // /
	TOKEN_PERCENT  // %
	TOKEN_EQ       // ==
	TOKEN_NEQ      // !=
	TOKEN_LT       // <
	TOKEN_GT       // >
	TOKEN_LTE      // <=
	TOKEN_GTE      // >=
	TOKEN_AND      // and
	TOKEN_OR       // or
	TOKEN_NOT      // not

	// Delimiters
	TOKEN_COMMA    // ,
	TOKEN_COLON    // :
	TOKEN_LPAREN   // (
	TOKEN_RPAREN   // )
	TOKEN_LBRACE   // {
	TOKEN_RBRACE   // }
	TOKEN_LBRACKET // [
	TOKEN_RBRACKET // ]
	TOKEN_DOLLAR   // $
	TOKEN_NEWLINE  // \n
	TOKEN_INDENT   // indentation
	TOKEN_DEDENT   // dedentation

	// Keywords (Wolf404 specific)
	TOKEN_HUNT    // hunt (function definition)
	TOKEN_SNIFF   // sniff (if)
	TOKEN_MISSING // missing (else/elif)
	TOKEN_TRACK   // track (for/while)
	TOKEN_BRING   // bring (return)
	TOKEN_HOWL    // howl (print/log)
	TOKEN_SUMMON  // summon (import)
	TOKEN_PACK    // pack (class/module)
	TOKEN_IN      // in
	TOKEN_RANGE   // range
	TOKEN_NIL     // nil
	TOKEN_PROWL   // prowl (go routine)
	TOKEN_MOLD    // mold (class)
	TOKEN_DOT     // .
)

var keywords = map[string]TokenType{
	"garap":      TOKEN_HUNT,
	"hunt":       TOKEN_HUNT,
	"gerombolan": TOKEN_MOLD,
	"mold":       TOKEN_MOLD,
	"menowo":     TOKEN_SNIFF,
	"sniff":      TOKEN_SNIFF,
	"yenora":     TOKEN_MISSING,
	"missing":    TOKEN_MISSING,
	"baleni":     TOKEN_TRACK,
	"track":      TOKEN_TRACK,
	"balekno":    TOKEN_BRING,
	"bring":      TOKEN_BRING,
	"ketok":      TOKEN_HOWL,
	"howl":       TOKEN_HOWL,
	"undang":     TOKEN_SUMMON,
	"summon":     TOKEN_SUMMON,
	"bungkus":    TOKEN_PACK,
	"pack":       TOKEN_PACK,
	"neng":       TOKEN_IN,
	"in":         TOKEN_IN,
	"deret":      TOKEN_RANGE,
	"range":      TOKEN_RANGE,
	"bener":      TOKEN_TRUE,
	"true":       TOKEN_TRUE,
	"salah":      TOKEN_FALSE,
	"false":      TOKEN_FALSE,
	"kopong":     TOKEN_NIL,
	"nil":        TOKEN_NIL,
	"lan":        TOKEN_AND,
	"and":        TOKEN_AND,
	"utowo":      TOKEN_OR,
	"or":         TOKEN_OR,
	"ora":        TOKEN_NOT,
	"not":        TOKEN_NOT,
	"playon":     TOKEN_PROWL,
	"prowl":      TOKEN_PROWL,
}

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Lexer represents the lexical analyzer
type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // current reading position
	ch           byte // current char
	line         int
	column       int
	indentStack  []int   // track indentation levels
	tokenQueue   []Token // Buffered tokens for DEDENT generation
}

// New creates a new Lexer
func New(input string) *Lexer {
	l := &Lexer{
		input:       input,
		line:        1,
		column:      0,
		indentStack: []int{0},
		tokenQueue:  []Token{},
	}
	l.readChar()
	return l
}

// readChar reads the next character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

// peekChar looks at the next character without advancing
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespaceExceptNewline() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) handleNewline() {
	// Consume the newline
	l.line++
	l.column = 0
	l.readChar()

	// Count indentation for the next line
	indentLen := 0
	for l.ch == ' ' || l.ch == '\t' {
		switch l.ch {
		case ' ':
			indentLen++
		case '\t':
			indentLen += 4 // Assume tab = 4 spaces
		}
		l.readChar()
	}

	// If empty line (newline again), just return, recursively NextToken will handle it
	if l.ch == '\n' || l.ch == '\r' || l.ch == 0 {
		return
	}

	currentIndent := l.indentStack[len(l.indentStack)-1]

	if indentLen > currentIndent {
		l.indentStack = append(l.indentStack, indentLen)
		l.tokenQueue = append(l.tokenQueue, Token{Type: TOKEN_NEWLINE, Literal: "\n", Line: l.line - 1})
		l.tokenQueue = append(l.tokenQueue, Token{Type: TOKEN_INDENT, Literal: "INDENT", Line: l.line})
	} else if indentLen < currentIndent {
		l.tokenQueue = append(l.tokenQueue, Token{Type: TOKEN_NEWLINE, Literal: "\n", Line: l.line - 1})
		// Pop dedents until matches
		for len(l.indentStack) > 1 && indentLen < l.indentStack[len(l.indentStack)-1] {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			l.tokenQueue = append(l.tokenQueue, Token{Type: TOKEN_DEDENT, Literal: "DEDENT", Line: l.line})
		}
	} else {
		// If equal, just a newline separator
		l.tokenQueue = append(l.tokenQueue, Token{Type: TOKEN_NEWLINE, Literal: "\n", Line: l.line - 1})
	}
}

// NextToken returns the next token
func (l *Lexer) NextToken() Token {
	// 1. Check if we have buffered tokens
	if len(l.tokenQueue) > 0 {
		tok := l.tokenQueue[0]
		l.tokenQueue = l.tokenQueue[1:]
		return tok
	}

	// 2. Check for EOF and remaining indents
	if l.ch == 0 {
		if len(l.indentStack) > 1 {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			return Token{Type: TOKEN_DEDENT, Literal: "DEDENT", Line: l.line, Column: l.column}
		}
		return Token{Type: TOKEN_EOF, Literal: "", Line: l.line, Column: l.column}
	}

	l.skipWhitespaceExceptNewline()

	var tok Token
	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '\n', '\r':
		if l.ch == '\r' && l.peekChar() == '\n' {
			l.readChar()
		}
		l.handleNewline()
		return l.NextToken()
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(TOKEN_ASSIGN, l.ch, l.line, l.column)
		}
	case '+':
		tok = newToken(TOKEN_PLUS, l.ch, l.line, l.column)
	case '-':
		tok = newToken(TOKEN_MINUS, l.ch, l.line, l.column)
	case '*':
		tok = newToken(TOKEN_ASTERISK, l.ch, l.line, l.column)
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken()
		}
		tok = newToken(TOKEN_SLASH, l.ch, l.line, l.column)
	case '%':
		tok = newToken(TOKEN_PERCENT, l.ch, l.line, l.column)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_NEQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(TOKEN_ILLEGAL, l.ch, l.line, l.column)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_LTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(TOKEN_LT, l.ch, l.line, l.column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TOKEN_GTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(TOKEN_GT, l.ch, l.line, l.column)
		}
	case ',':
		tok = newToken(TOKEN_COMMA, l.ch, l.line, l.column)
	case ':':
		tok = newToken(TOKEN_COLON, l.ch, l.line, l.column)
	case '.':
		tok = newToken(TOKEN_DOT, l.ch, l.line, l.column)
	case '(':
		tok = newToken(TOKEN_LPAREN, l.ch, l.line, l.column)
	case ')':
		tok = newToken(TOKEN_RPAREN, l.ch, l.line, l.column)
	case '{':
		tok = newToken(TOKEN_LBRACE, l.ch, l.line, l.column)
	case '}':
		tok = newToken(TOKEN_RBRACE, l.ch, l.line, l.column)
	case '[':
		tok = newToken(TOKEN_LBRACKET, l.ch, l.line, l.column)
	case ']':
		tok = newToken(TOKEN_RBRACKET, l.ch, l.line, l.column)
	case '$':
		tok = newToken(TOKEN_DOLLAR, l.ch, l.line, l.column)
	case '"':
		tok.Type = TOKEN_STRING
		tok.Literal = l.readString()

	case 0:
		tok.Literal = ""
		tok.Type = TOKEN_EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			literal, isFloat := l.readNumber()
			tok.Literal = literal
			if isFloat {
				tok.Type = TOKEN_FLOAT
			} else {
				tok.Type = TOKEN_INT
			}
			return tok
		} else {
			tok = newToken(TOKEN_ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType TokenType, ch byte, line, column int) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (string, bool) {
	position := l.position
	isFloat := false

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position], isFloat
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENT
}

// TokenTypeString returns string representation of token type
func TokenTypeString(t TokenType) string {
	switch t {
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_ILLEGAL:
		return "ILLEGAL"
	case TOKEN_IDENT:
		return "IDENT"
	case TOKEN_INT:
		return "INT"
	case TOKEN_FLOAT:
		return "FLOAT"
	case TOKEN_STRING:
		return "STRING"
	case TOKEN_TRUE:
		return "TRUE"
	case TOKEN_FALSE:
		return "FALSE"
	case TOKEN_ASSIGN:
		return "="
	case TOKEN_PLUS:
		return "+"
	case TOKEN_MINUS:
		return "-"
	case TOKEN_ASTERISK:
		return "*"
	case TOKEN_SLASH:
		return "/"
	case TOKEN_PERCENT:
		return "%"
	case TOKEN_EQ:
		return "=="
	case TOKEN_NEQ:
		return "!="
	case TOKEN_LT:
		return "<"
	case TOKEN_GT:
		return ">"
	case TOKEN_LTE:
		return "<="
	case TOKEN_GTE:
		return ">="
	case TOKEN_AND:
		return "AND"
	case TOKEN_OR:
		return "OR"
	case TOKEN_NOT:
		return "NOT"
	case TOKEN_HUNT:
		return "HUNT"
	case TOKEN_SNIFF:
		return "SNIFF"
	case TOKEN_MISSING:
		return "MISSING"
	case TOKEN_TRACK:
		return "TRACK"
	case TOKEN_BRING:
		return "BRING"
	case TOKEN_HOWL:
		return "HOWL"
	case TOKEN_SUMMON:
		return "SUMMON"
	case TOKEN_COMMA:
		return ","
	case TOKEN_COLON:
		return ":"
	case TOKEN_LPAREN:
		return "("
	case TOKEN_RPAREN:
		return ")"
	case TOKEN_LBRACE:
		return "{"
	case TOKEN_RBRACE:
		return "}"
	case TOKEN_LBRACKET:
		return "["
	case TOKEN_RBRACKET:
		return "]"
	case TOKEN_DOLLAR:
		return "$"
	case TOKEN_PACK:
		return "PACK"
	case TOKEN_NEWLINE:
		return "NEWLINE"
	case TOKEN_PROWL:
		return "PROWL"
	case TOKEN_INDENT:
		return "INDENT"
	case TOKEN_DEDENT:
		return "DEDENT"
	case TOKEN_MOLD:
		return "MOLD"
	case TOKEN_DOT:
		return "."
	default:
		return "UNKNOWN"
	}
}
