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
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
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

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	p.registerPrefix(token.LPAREN, p.parseGroupedExpressions)
	p.registerPrefix(token.IF, p.parseIfExpression)

	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	p.registerPrefix(token.STRING, p.parseStringLiteral)

	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
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
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// no identifiers in the function signature.
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	if !p.curTokenIs(token.RBRACE) {
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
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	/*
			For having + right associative:
			if expression.Operator == "+" {
		        expression.Right = p.parseExpression(precedence - 1)
		    } else {
		        expression.Right = p.parseExpression(precedence)
			}
	*/

	expression.Right = p.parseExpression(precedence)

	return expression

}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

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
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.WORD:
		return p.parseWordStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
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

	if !p.expectPeek(token.COLON) {
		return nil
	}

	// Expecting an identifier after the :
	if !p.peekTokenIs(token.IDENT) {
		return nil
	}

	p.nextToken()

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()

		if !p.peekTokenIs(token.STRING) {
			return nil
		}
		p.nextToken()

		stmt.Definition = p.curToken.Literal

		if !p.peekTokenIs(token.RBRACE) {
			return nil
		}

		stmt.Defined = true
	} else {
		p.debug()
	}

	p.nextToken()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt

}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	// for !p.curTokenIs(token.SEMICOLON) {
	// 	p.nextToken()
	// }
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
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
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}