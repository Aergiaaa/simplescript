package ast

import (
	"testing"

	"github.com/Aergiaaa/idiotic_interpreter/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "outvar",
					},
					Value: "outvar",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "inpvar",
					},
					Value: "inpvar",
				},
			},
		},
	}

	if program.String() != "let outvar = inpvar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
