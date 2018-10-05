package parser

import (
	"fmt"
	"strconv"
	"wordbuilder/ast"
	"wordbuilder/lexer"
	"wordbuilder/token"
)

const (
	_ int = iota
	// LOWEST ...
	LOWEST
	// EQUALS ...
	EQUALS // ==
	// LESSGREATER ...
	LESSGREATER // > or <
	// SUM ...
	SUM // +
	// PRODUCT ...
	PRODUCT // *
	// PREFIX ...
	PREFIX // -X or !X
	// CALL ...
	CALL // myFunction(X)
	// INDEX ...
	INDEX
)

var precedences = map[token.Type]int{
	token.Eq:          EQUALS,
	token.NotEq:       EQUALS,
	token.Lt:          LESSGREATER,
	token.Gt:          LESSGREATER,
	token.Plus:        SUM,
	token.Minus:       SUM,
	token.Slash:       PRODUCT,
	token.Asterisk:    PRODUCT,
	token.LeftParen:   CALL,
	token.LeftBracket: INDEX,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

type Error struct {
	Error      string
	LineNumber int
}

func (err *Error) String() string {
	return fmt.Sprintf("Line %d:%s", err.LineNumber, err.Error)
}

// Parser ...
type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	errors    []Error

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// New ...
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []Error{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.Ident, p.parseIdentifier)
	p.registerPrefix(token.Int, p.parseIntegerLiteral)

	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Slash, p.parseInfixExpression)
	p.registerInfix(token.Asterisk, p.parseInfixExpression)
	p.registerInfix(token.Eq, p.parseInfixExpression)
	p.registerInfix(token.NotEq, p.parseInfixExpression)
	p.registerInfix(token.Lt, p.parseInfixExpression)
	p.registerInfix(token.Gt, p.parseInfixExpression)

	p.registerPrefix(token.True, p.parseBoolean)
	p.registerPrefix(token.False, p.parseBoolean)

	p.registerPrefix(token.LeftParen, p.parseGroupedExpressions)
	p.registerPrefix(token.If, p.parseIfExpression)

	p.registerPrefix(token.Function, p.parseFunctionLiteral)
	p.registerInfix(token.LeftParen, p.parseCallExpression)

	p.registerPrefix(token.String, p.parseStringLiteral)
	p.registerPrefix(token.Word, p.parseStringLiteral)
	p.registerPrefix(token.Me, p.parseStringLiteral)
	p.registerPrefix(token.Tr, p.parseStringLiteral)
	p.registerPrefix(token.Ref, p.parseStringLiteral)

	p.registerPrefix(token.LeftBracket, p.parseArrayLiteral)
	p.registerInfix(token.LeftBracket, p.parseIndexExpression)

	p.registerPrefix(token.LeftBrace, p.parseHashLiteral)
	p.registerPrefix(token.Colon, p.parseStringLiteral)

	// p.prefixParseFns = make(map[token.Type]prefixParseFn)
	// p.registerPrefix(token.Ident, p.parseIdentifier)
	// p.registerPrefix(token.Int, p.parseIntegerLiteral)
	// p.registerPrefix(token.Bang, p.parsePrefixExpression)
	// p.registerPrefix(token.Minus, p.parsePrefixExpression)
	// p.registerPrefix(token.True, p.parseBoolean)
	// p.registerPrefix(token.False, p.parseBoolean)
	// p.registerPrefix(token.LeftParen, p.parseGroupedExpressions)
	// p.registerPrefix(token.If, p.parseIfExpression)
	// p.registerPrefix(token.Function, p.parseFunctionLiteral)
	// p.registerPrefix(token.String, p.parseStringLiteral)
	// p.registerPrefix(token.LeftBracket, p.parseArrayLiteral)
	// p.registerPrefix(token.LeftBrace, p.parseHashLiteral)

	// p.infixParseFns = make(map[token.Type]infixParseFn)
	// p.registerInfix(token.Plus, p.parseInfixExpression)
	// p.registerInfix(token.Minus, p.parseInfixExpression)
	// p.registerInfix(token.Slash, p.parseInfixExpression)
	// p.registerInfix(token.Asterisk, p.parseInfixExpression)
	// p.registerInfix(token.Eq, p.parseInfixExpression)
	// p.registerInfix(token.NotEq, p.parseInfixExpression)
	// p.registerInfix(token.Lt, p.parseInfixExpression)
	// p.registerInfix(token.Gt, p.parseInfixExpression)
	// p.registerInfix(token.LeftParen, p.parseCallExpression)
	// p.registerInfix(token.LeftBracket, p.parseIndexExpression)

	return p
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RightBrace) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.Colon) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RightBrace) && !p.expectPeek(token.Comma) {
			return nil
		}
	}

	if !p.expectPeek(token.RightBrace) {
		return nil
	}

	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RightBracket) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RightBracket)
	return array
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	// exp := &ast.CallExpression{Token: p.curToken, Function: function}
	// exp.Arguments = p.parseExpressionList(token.RightParen)
	// return exp

	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	exp.Arguments = p.parseExpressionList(token.RightParen)
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RightParen) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RightParen) {
		return nil
	}

	return args
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LeftParen) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LeftBrace) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// no identifiers in the function signature.
	if p.peekTokenIs(token.RightParen) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RightParen) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LeftParen) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RightParen) {
		return nil
	}

	if !p.expectPeek(token.LeftBrace) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.Else) {
		p.nextToken()

		if !p.expectPeek(token.LeftBrace) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression

}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.curToken,
	}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RightBrace) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseGroupedExpressions() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RightParen) {
		return nil
	}
	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.True),
	}
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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// Errors ...
func (p *Parser) Errors() []Error {
	return p.errors
}

func (p *Parser) peekError(t token.Type) {
	p.debug()
	msg := fmt.Sprintf("expected next token to be [%s], got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, Error{Error: msg, LineNumber: p.l.CurrentLine()})
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram ...
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
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
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	case token.Word:
		return p.parseWordStatement()
	case token.Ref:
		return p.parseReferenceStatement()
	case token.Cpt:
		return p.parseConceptStatement()
	case token.Tr:
		return p.parseTranslationStatement()
	case token.Me:
		return p.parseMeThoughtStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	// for !p.curTokenIs(token.SEMICOLON) {
	// 	p.nextToken()
	// }

	return stmt
}

func (p *Parser) debug() {
	fmt.Printf("[debug]: curToken is: %q\n", p.curToken)
	fmt.Printf("[debug]: peekToken is: %q\n", p.peekToken)
}

func (p *Parser) parseWordStatement() *ast.WordStatement {

	stmt := &ast.WordStatement{Token: p.curToken}

	if !p.expectPeek(token.Colon) {
		return nil
	}

	// Expecting an identifier after the :
	if !p.peekTokenIs(token.String) {
		return nil
	}

	p.nextToken()

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.LeftBrace) {
		p.nextToken()

		if !p.expectPeek(token.String) {
			return nil
		}

		stmt.Definition = p.curToken.Literal
		stmt.Value = p.parseExpression(LOWEST)

		if !p.expectPeek(token.RightBrace) {
			return nil
		}
		stmt.Defined = true
	}

	p.nextToken()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	fmt.Println("About to return this ... ")
	return stmt
}

func (p *Parser) parseTranslationStatement() *ast.TranslationStatement {
	stmt := &ast.TranslationStatement{Token: p.curToken}

	if !p.expectPeek(token.Colon) {
		return nil
	}

	// Expecting an identifier after the :
	if !p.peekTokenIs(token.Ident) {
		return nil
	}

	p.nextToken()

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.LeftBrace) {
		p.nextToken()

		if !p.expectPeek(token.String) {
			return nil
		}

		stmt.Definition = p.curToken.Literal
		stmt.Value = p.parseExpression(LOWEST)

		if !p.expectPeek(token.RightBrace) {
			return nil
		}
		stmt.Defined = true
	}

	p.nextToken()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseMeThoughtStatement() *ast.MeThoughtStatement {
	stmt := &ast.MeThoughtStatement{}

	if !p.expectPeek(token.Colon) {
		return nil
	}

	// Expecting an identifier after the :
	if !p.expectPeek(token.LeftBrace) {
		return nil
	}

	if p.peekTokenIs(token.String) {
		p.nextToken()
		stmt.Content = p.parseExpression(LOWEST).String()
		if !p.expectPeek(token.RightBrace) {
			return nil
		}
	}

	p.nextToken()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReferenceStatement() *ast.ReferenceStatement {
	stmt := &ast.ReferenceStatement{Token: p.curToken}

	if !p.expectPeek(token.Colon) {
		return nil
	}

	// Expecting an identifier after the :
	if !p.peekTokenIs(token.String) {
		return nil
	}

	p.nextToken()

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.LeftBrace) {
		p.nextToken()

		if !p.expectPeek(token.String) {
			return nil
		}

		stmt.Definition = p.curToken.Literal
		stmt.Value = p.parseExpression(LOWEST)

		if !p.expectPeek(token.RightBrace) {
			return nil
		}
		stmt.Defined = true
	}

	p.nextToken()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConceptStatement() *ast.ConceptStatement {
	stmt := &ast.ConceptStatement{Token: p.curToken}

	if !p.expectPeek(token.Colon) {
		return nil
	}

	// Expecting an identifier after the :
	if !p.peekTokenIs(token.String) {
		return nil
	}

	p.nextToken()

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.LeftBrace) {
		p.nextToken()

		if !p.expectPeek(token.String) {
			return nil
		}

		stmt.Definition = p.curToken.Literal
		stmt.Value = p.parseExpression(LOWEST)

		if !p.expectPeek(token.RightBrace) {
			return nil
		}
		stmt.Defined = true
	}

	p.nextToken()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt

}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.Assign) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as Integer", p.curToken.Literal)
		p.errors = append(p.errors, Error{Error: msg, LineNumber: p.l.CurrentLine()})
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, Error{Error: msg, LineNumber: p.l.CurrentLine()})
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}
