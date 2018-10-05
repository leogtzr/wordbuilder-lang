package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	String = "STRING"
	Colon  = ":"

	Illegal = "ILLEGAL"
	EOF     = "EOF"

	// Ident Identifiers + literals
	Ident = "IDENT" // add, foobar, x, y, ...
	Int   = "INT"   // 1343456

	// Assign ... Operators
	Assign = "="
	Plus   = "+"

	// Comma ...
	Comma     = ","
	Semicolon = ";"

	LeftParen  = "("
	RightParen = ")"
	LeftBrace  = "{"
	RightBrace = "}"

	// Keywords
	Function = "FUNCTION"
	Let      = "LET"

	Word  = "WORD"
	Ref   = "REF"
	Cpt   = "CPT"
	Tr    = "TR"
	Me    = "ME"
	Quote = "QUOTE"

	True   = "TRUE"
	False  = "FALSE"
	If     = "IF"
	Else   = "ELSE"
	Return = "RETURN"

	// Operators
	// MINUS ...
	Minus    = "-"
	Bang     = "!"
	Asterisk = "*"
	Slash    = "/"

	Lt = "<"
	Gt = ">"

	Eq    = "=="
	NotEq = "!="

	LeftBracket  = "["
	RightBracket = "]"
)

var keywords = map[string]Type{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
	"word":   Word,
	"ref":    Ref,
	"cpt":    Cpt,
	"tr":     Tr,
	"me":     Me,
	"quote":  Quote,
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
