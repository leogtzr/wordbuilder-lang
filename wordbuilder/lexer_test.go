package wordbuilder

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	// 	input := `let five = 5;
	// let ten = 10;
	// let add = fn(x, y) {
	// 	x + y;
	// };
	// let result = add(five, ten);
	// !-/*5;
	// 5 < 10 > 5;
	// if (5 < 10) {
	// 	return true;
	// } else {
	// 	return false;
	// }
	// 10 == 10;
	// 10 != 9;
	// "foobar"
	// "foo bar"
	// [1, 2];
	// {"foo": "bar"}
	// `

	// tests := []struct {
	// 	expectedType    TokenType
	// 	expectedLiteral string
	// }{
	// 	{LET, "let"},
	// 	{IDENT, "five"},
	// 	{ASSIGN, "="},
	// 	{INT, "5"},
	// 	{SEMICOLON, ";"},
	// 	{LET, "let"},
	// 	{IDENT, "ten"},
	// 	{ASSIGN, "="},
	// 	{INT, "10"},
	// 	{SEMICOLON, ";"},
	// 	{LET, "let"},
	// 	{IDENT, "add"},
	// 	{ASSIGN, "="},
	// 	{FUNCTION, "fn"},
	// 	{LPAREN, "("},
	// 	{IDENT, "x"},
	// 	{COMMA, ","},
	// 	{IDENT, "y"},
	// 	{RPAREN, ")"},
	// 	{LBRACE, "{"},
	// 	{IDENT, "x"},
	// 	{PLUS, "+"},
	// 	{IDENT, "y"},
	// 	{SEMICOLON, ";"},
	// 	{RBRACE, "}"},
	// 	{SEMICOLON, ";"},
	// 	{LET, "let"},
	// 	{IDENT, "result"},
	// 	{ASSIGN, "="},
	// 	{IDENT, "add"},
	// 	{LPAREN, "("},
	// 	{IDENT, "five"},
	// 	{COMMA, ","},
	// 	{IDENT, "ten"},
	// 	{RPAREN, ")"},
	// 	{SEMICOLON, ";"},
	// 	{BANG, "!"},
	// 	{MINUS, "-"},
	// 	{SLASH, "/"},
	// 	{ASTERISK, "*"},
	// 	{INT, "5"},
	// 	{SEMICOLON, ";"},
	// 	{INT, "5"},
	// 	{LT, "<"},
	// 	{INT, "10"},
	// 	{GT, ">"},
	// 	{INT, "5"},
	// 	{SEMICOLON, ";"},
	// 	{IF, "if"},
	// 	{LPAREN, "("},
	// 	{INT, "5"},
	// 	{LT, "<"},
	// 	{INT, "10"},
	// 	{RPAREN, ")"},
	// 	{LBRACE, "{"},
	// 	{RETURN, "return"},
	// 	{TRUE, "true"},
	// 	{SEMICOLON, ";"},
	// 	{RBRACE, "}"},
	// 	{ELSE, "else"},
	// 	{LBRACE, "{"},
	// 	{RETURN, "return"},
	// 	{FALSE, "false"},
	// 	{SEMICOLON, ";"},
	// 	{RBRACE, "}"},
	// 	{INT, "10"},
	// 	{EQ, "=="},
	// 	{INT, "10"},
	// 	{SEMICOLON, ";"},
	// 	{INT, "10"},
	// 	{NOT_EQ, "!="},
	// 	{INT, "9"},
	// 	{SEMICOLON, ";"},
	// 	{STRING, "foobar"},
	// 	{STRING, "foo bar"},
	// 	{LBRAKET, "["},
	// 	{INT, "1"},
	// 	{COMMA, ","},
	// 	{INT, "2"},
	// 	{RBRAKET, "]"},
	// 	{SEMICOLON, ";"},

	// 	{LBRACE, "{"},
	// 	{STRING, "foo"},
	// 	{COLON, ":"},
	// 	{STRING, "bar"},
	// 	{RBRACE, "}"},

	// 	{EOF, ""},
	// }

	input := `w: estridente -> {
		holis
		}

	ref: Sucesos de Scottsboro -> (p)
	w:lunfardo->(p)
	`
	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{WORD, "w"},
		{COLON, ":"},
		{IDENT, "estridente"},
		{BEGIN_W_REF_OR_C, "->"},
		{TEXT, "{\n\t\tholis\n\t\t}"},

		{REF, "ref"},
		{COLON, ":"},
		{IDENT, "Sucesos de Scottsboro"},
		{BEGIN_W_REF_OR_C, "->"},
		{LPAREN, "("},
		{IDENT, "p"},
		{RPAREN, ")"},

		{WORD, "w"},
		{COLON, ":"},
		{IDENT, "lunfardo"},
		{BEGIN_W_REF_OR_C, "->"},
		{LPAREN, "("},
		{IDENT, "p"},
		{RPAREN, ")"},

		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestCountWords(t *testing.T) {

	input := `w: -> estridente (p)`
	lexer := New(input)
	tokenCount := 0
	expectedTokenCount := 5

	for tok := lexer.NextToken(); tok.Type != EOF; tokenCount++ {
		tok = lexer.NextToken()
	}

	if tokenCount != expectedTokenCount {
		t.Errorf("Error, expected: %d tokens, we got: %d", expectedTokenCount, tokenCount)
	}

}
