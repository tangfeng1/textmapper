package test_test

import (
	"testing"

	"github.com/inspirer/textmapper/tm-go/parsers/test"
	"github.com/inspirer/textmapper/tm-parsers/parsertest"
)

var parseTests = []struct {
	nt     test.NodeType
	inputs []string
}{

	{test.Test, []string{
		` «decl2 decl1(a)»`,
	}},
	{test.Block, []string{
		`«{decl2}»`,
		`«{-decl2}»`,
		`«{--decl2}»`,
		`«{--}»`,
		`«{-}»`,
		`«{}»`,
	}},
	{test.Negation, []string{
		`{«-»decl2}`,
		`{«--»decl2}`,
		`{«--»}`,
		`{«-»}`,
		`{«-»{«-»decl2}}`,
	}},
	{test.Decl1, []string{
		`{«decl1(a.b.c.d123)»}`,
	}},
	{test.Decl2, []string{
		`«decl2»`,
	}},

	{test.MultiLineComment, []string{
		` decl2 «/* ****/» decl1(a)`,
	}},
	{test.SingleLineComment, []string{
		` decl2 «// abc»
		 decl1(a)`,
	}},
	{test.InvalidToken, []string{
		` decl2 «%» `,
	}},
	{test.Identifier, []string{
		` decl1(«abc».«def1») `,
	}},
}

func TestParser(t *testing.T) {
	l := new(test.Lexer)
	p := new(test.Parser)

	seen := map[test.NodeType]bool{}
	for _, tc := range parseTests {
		seen[tc.nt] = true
		for _, input := range tc.inputs {
			pt := parsertest.New(t, tc.nt.String(), input)
			l.Init(pt.Source())
			errHandler := func(se test.SyntaxError) bool {
				pt.ConsumeError(t, se.Offset, se.Endoffset)
				return true
			}
			p.Init(errHandler, func(tn test.NodeType, offset, endoffset int) {
				if tn == tc.nt {
					pt.Consume(t, offset, endoffset)
				}
			})
			pt.Done(t, p.Parse(l))
		}
	}
	for n := test.NodeType(1); n < test.NodeTypeMax; n++ {
		if !seen[n] {
			t.Errorf("%v is not tested", n)
		}
	}
}
