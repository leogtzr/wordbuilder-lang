package wordbuilder

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	lex := &Lexer{input: input}
	lex.readChar()
	return lex
}

func (lex *Lexer) readChar() {
	if lex.readPosition >= len(lex.input) {
		lex.ch = 0
	} else {
		lex.ch = lex.input[lex.readPosition]
	}

	lex.position = lex.readPosition
	lex.readPosition++
}

func (lex *Lexer) NextToken() Token {
	var tok Token

	lex.skipWhitespace()

	switch lex.ch {
	case '=':
		if lex.peekChar() == '=' {
			ch := lex.ch
			lex.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(lex.ch)}
		} else {
			tok = newToken(ASSIGN, lex.ch)
		}
	case '+':
		tok = newToken(PLUS, lex.ch)
	case '-':
		if lex.peekChar() == '>' {
			ch := lex.ch
			lex.readChar()
			tok = Token{Type: BEGIN_W_REF_OR_C, Literal: string(ch) + string(lex.ch)}
		} else {
			tok = newToken(MINUS, lex.ch)
		}
	case '!':
		if lex.peekChar() == '=' {
			ch := lex.ch
			lex.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(lex.ch)}
		} else {
			tok = newToken(BANG, lex.ch)
		}
	case '{':
		tok.Type = TEXT
		tok.Literal = lex.readText()
	case '/':
		tok = newToken(SLASH, lex.ch)
	case '*':
		tok = newToken(ASTERISK, lex.ch)
	case '<':
		tok = newToken(LT, lex.ch)
	case '>':
		tok = newToken(GT, lex.ch)
	case ';':
		tok = newToken(SEMICOLON, lex.ch)
	case '(':
		tok = newToken(LPAREN, lex.ch)
	case ')':
		tok = newToken(RPAREN, lex.ch)
	case ',':
		tok = newToken(COMMA, lex.ch)
	case '}':
		tok = newToken(RBRACE, lex.ch)
	case '[':
		tok = newToken(LBRAKET, lex.ch)
	case ']':
		tok = newToken(RBRAKET, lex.ch)
	case ':':
		tok = newToken(COLON, lex.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = lex.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(lex.ch) {
			tok.Literal = lex.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(lex.ch) {
			tok.Type = INT
			tok.Literal = lex.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, lex.ch)
		}
	}

	lex.readChar()
	return tok
}

func (lex *Lexer) peekChar() byte {
	if lex.readPosition >= len(lex.input) {
		return 0
	}

	return lex.input[lex.readPosition]
}

func (lex *Lexer) readIdentifier() string {
	pos := lex.position
	for isLetter(lex.ch) {
		lex.readChar()
	}

	return lex.input[pos:lex.position]
}

func (lex *Lexer) readText() string {
	pos := lex.position
	for lex.ch != '}' && lex.ch != 0 {
		lex.readChar()
	}
	if lex.ch == '}' {
		return lex.input[pos:(lex.position + 1)]
	}
	return lex.input[pos:(lex.position)]
}

func (lex *Lexer) readNumber() string {
	position := lex.position
	for isDigit(lex.ch) {
		lex.readChar()
	}
	return lex.input[position:lex.position]
}

func (lex *Lexer) skipWhitespace() {
	for lex.ch == ' ' || lex.ch == '\t' || lex.ch == '\n' || lex.ch == '\r' {
		lex.readChar()
	}
}

func (lex *Lexer) readString() string {
	position := lex.position + 1
	for {
		lex.readChar()
		if lex.ch == '"' {
			break
		}
	}

	return lex.input[position:lex.position]
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
