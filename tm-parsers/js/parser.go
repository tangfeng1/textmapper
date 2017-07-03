// generated by Textmapper; DO NOT EDIT

package js

func (p *Parser) Parse(lexer *Lexer) error {
	return p.parse(0, 3884, lexer)
}

func (p *Parser) applyRule(rule int32, lhs *stackEntry, rhs []stackEntry) {
	switch rule {
	case 2202: // IterationStatement : 'for' '(' 'async' 'of' AssignmentExpression_In ')' Statement
		p.listener(IdentifierReference, rhs[2].sym.offset, rhs[2].sym.endoffset)
	case 2216: // IterationStatement_Await : 'for' '(' 'async' 'of' AssignmentExpression_Await_In ')' Statement_Await
		p.listener(IdentifierReference, rhs[2].sym.offset, rhs[2].sym.endoffset)
	case 2230: // IterationStatement_Yield : 'for' '(' 'async' 'of' AssignmentExpression_In_Yield ')' Statement_Yield
		p.listener(IdentifierReference, rhs[2].sym.offset, rhs[2].sym.endoffset)
	}
	nt := ruleNodeType[rule]
	if nt == 0 {
		return
	}
	p.listener(nt, lhs.sym.offset, lhs.sym.endoffset)
}

const errSymbol = 2