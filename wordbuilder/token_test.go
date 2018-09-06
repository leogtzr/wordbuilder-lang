package wordbuilder

import (
	"testing"
)

func TestLookupIdent(t *testing.T) {
	tokenType := LookupIdent("ref")
	if tokenType != REF {
		t.Fatal("Expecting keyword.")
	}

	tokenType = LookupIdent("w")
	if tokenType != WORD {
		t.Fatal("Expecting keyword.")
	}

	tokenType = LookupIdent("c")
	if tokenType != CONCEPT {
		t.Fatal("Expecting keyword.")
	}
}
