package json_test

import (
	"testing"

	"github.com/inspirer/textmapper/tm-go/parsers/json"
	"github.com/inspirer/textmapper/tm-parsers/parsertest"
)

var jsParseTests = []struct {
	nt     json.NodeType
	inputs []string
}{

	{json.EmptyObject, []string{
		`«{}»`,
		`«{ /* comment */ }»`,
		`{"aa": «{}» }`,
	}},
	{json.JSONObject, []string{
		`«{ "a" : "b" }»`,
		`«{ "a" : ["b"] }»`,
		`«{ "a" : {} }»`,
		`«{ "a" : «{"q":B}» }»`,
	}},
	{json.JSONArray, []string{
		`{ "a" : «["b"]» }`,
		` «[]» `,
	}},
	{json.JSONText, []string{
		`«{ "a" : ["b", A] }»`,
		` «"aa"» `,
		` «A» `,
	}},
	{json.JSONMember, []string{
		`[{ «"a" : ["b"]», «"q":[]» }]`,
	}},
	{json.JSONValue, []string{
		`«{ "a" : «[«"b"»]» }»`,
		` «"aa"» `,
	}},
	{json.InvalidToken, []string{
		`  «%» null `,
	}},
	{json.NonExistingType, []string{}},
	{json.MultiLineComment, []string{
		`{ "a"«/* abc */» : [] }`,
	}},
	{json.JsonString, []string{
		`{ «"a"» : [«"b"»] }`,
	}},
}

func TestParser(t *testing.T) {
	l := new(json.Lexer)
	p := new(json.Parser)

	seen := map[json.NodeType]bool{}
	for _, tc := range jsParseTests {
		seen[tc.nt] = true
		for _, input := range tc.inputs {
			test := parsertest.New(t, tc.nt.String(), input)
			l.Init(test.Source())
			errHandler := func(se json.SyntaxError) bool {
				test.ConsumeError(t, se.Offset, se.Endoffset)
				return true
			}
			p.Init(errHandler, func(nt json.NodeType, offset, endoffset int) {
				if nt == tc.nt {
					test.Consume(t, offset, endoffset)
				}
			})
			test.Done(t, p.Parse(l))
		}
	}
	for n := json.NodeType(1); n < json.NodeTypeMax; n++ {
		if !seen[n] {
			t.Errorf("%v is not tested", n)
		}
	}
}

func BenchmarkParser(b *testing.B) {
	l := new(json.Lexer)
	p := new(json.Parser)
	onError := func(se json.SyntaxError) bool {
		b.Errorf("unexpected: %v", se)
		return false
	}

	p.Init(onError, func(t json.NodeType, offset, endoffset int) {})
	for i := 0; i < b.N; i++ {
		l.Init(jsonExample)
		p.Parse(l)
	}
	b.SetBytes(int64(len(jsonExample)))
}
