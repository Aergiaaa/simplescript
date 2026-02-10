package evaluator

import (
	"testing"

	"github.com/Aergiaaa/simplescript/object"
)

func TestQuoteUnqoute(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`let foobar = 8;
            quote(foobar)`,
			`foobar`,
		},
		{
			`let foobar = 8;
            quote(unquote(foobar))`,
			`8`,
		},
		{
			`quote(unquote(true))`,
			`true`,
		},
		{
			`quote(unquote(true == false))`,
			`false`,
		},
		{
			`quote(unquote(quote(4 + 4)))`,
			`(4 + 4)`,
		},
		{
			`let quotedInfixExpression = quote(4 + 4);
            quote(unquote(4 + 4) + unquote(quotedInfixExpression))`,
			`(8 + (4 + 4))`,
		},
		{
			`quote(unquote(5))`,
			`5`,
		},
		{
			`quote(unquote(6+7))`,
			`13`,
		},
		{
			`quote(6 + unquote(3+4))`,
			`(6 + 7)`,
		},
		{
			`quote(unquote(2+4) + 7)`,
			`(6 + 7)`,
		},
	}

	for _, tt := range tests {
		evaled := testEval(tt.input)
		quote, ok := evaled.(*object.Quote)
		if !ok {
			t.Fatalf("expected quote, got=%T (%+v)", evaled, evaled)
		}

		if quote.Node == nil {
			t.Fatalf("node is nil")
		}

		if quote.Node.String() != tt.expected {
			t.Errorf("no equal. expected=%s got=%s", tt.expected, quote.Node.String())
		}
	}
}

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote(5)`,
			`5`,
		},
		{
			`quote(5 + 8)`,
			`(5 + 8)`,
		},
		{
			`quote(sixseven)`,
			`sixseven`,
		},
		{
			`quote(six + seven)`,
			`(six + seven)`,
		},
	}

	for _, tt := range tests {
		evaled := testEval(tt.input)
		quote, ok := evaled.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote, got=%T (%+v)", evaled, evaled)
		}

		if quote.Node == nil {
			t.Fatalf("node is empty")
		}

		if quote.Node.String() != tt.expected {
			t.Errorf("quote not equal, expected=%q got=%q", tt.expected, quote.Node.String())
		}
	}
}
