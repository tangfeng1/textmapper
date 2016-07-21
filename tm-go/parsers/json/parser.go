package json

import (
	"fmt"
)

// Parser is a table-driven LALR parser for json.
type Parser struct {
	err      ErrorHandler
	listener Listener

	stack []node
	lexer *Lexer
	next  symbol
}

type symbol struct {
	symbol    int32
	offset    int
	endoffset int
}

type node struct {
	sym   symbol
	state int8
	value int
}

func (p *Parser) Init(err ErrorHandler, l Listener) {
	p.err = err
	p.listener = l
}

const (
	startStackSize = 512
	noToken        = int32(UNAVAILABLE)
	eoiToken       = int32(EOI)
	debugSyntax    = false
)

func (p *Parser) Parse(lexer *Lexer) (bool, int) {
	return p.parse(1, 44, lexer)
}

func (p *Parser) parse(start, end int8, lexer *Lexer) (bool, int) {
	if cap(p.stack) < startStackSize {
		p.stack = make([]node, 0, startStackSize)
	}
	state := start
	recovering := 0

	p.stack = append(p.stack[:0], node{state: state})
	p.lexer = lexer
	p.next.symbol = int32(lexer.Next())
	p.next.offset, p.next.endoffset = lexer.Pos()

	for state != end {
		action := tmAction[state]
		if action < -2 {
			// Lookahead is needed.
			if p.next.symbol == noToken {
				p.next.symbol = int32(p.lexer.Next())
				p.next.offset, p.next.endoffset = p.lexer.Pos()
			}
			action = lalr(action, p.next.symbol)
		}

		if action >= 0 {
			// Reduce.
			rule := action
			ln := int(tmRuleLen[rule])

			var node node
			node.sym.symbol = tmRuleSymbol[rule]
			if ln == 0 {
				node.sym.offset, _ = lexer.Pos()
				node.sym.endoffset = node.sym.offset
			} else {
				node.sym.offset = p.stack[len(p.stack)-ln].sym.offset
				node.sym.endoffset = p.stack[len(p.stack)-1].sym.endoffset
			}
			p.applyRule(rule, &node, p.stack[len(p.stack)-ln:])
			if debugSyntax {
				fmt.Printf("reduced to: %v\n", Symbol(node.sym.symbol))
			}
			p.stack = p.stack[:len(p.stack)-ln]
			state = gotoState(p.stack[len(p.stack)-1].state, node.sym.symbol)
			node.state = state
			p.stack = append(p.stack, node)

		} else if action == -1 {
			// Shift.
			if p.next.symbol == noToken {
				p.next.symbol = int32(lexer.Next())
				p.next.offset, p.next.endoffset = lexer.Pos()
			}
			state = gotoState(state, p.next.symbol)
			p.stack = append(p.stack, node{
				sym:   p.next,
				state: state,
			})
			if debugSyntax {
				fmt.Printf("shift: %v (%s)\n", Symbol(p.next.symbol), lexer.Text())
			}
			if state != -1 && p.next.symbol != eoiToken {
				p.next.symbol = noToken
			}
			if recovering > 0 {
				recovering--
			}
		}

		if action == -2 || state == -1 {
			if p.recover() {
				state = p.stack[len(p.stack)-1].state
				if recovering == 0 {
					offset, endoffset := lexer.Pos()
					line := lexer.Line()
					p.err(line, offset, endoffset-offset, "syntax error")
				}
				if recovering >= 3 {
					p.next.symbol = int32(p.lexer.Next())
					p.next.offset, p.next.endoffset = lexer.Pos()
				}
				recovering = 4
				continue
			}
			if len(p.stack) == 0 {
				state = start
				p.stack = append(p.stack, node{state: state})
			}
			break
		}
	}

	if state != end {
		if recovering > 0 {
			return false, 0
		}
		offset, endoffset := lexer.Pos()
		line := lexer.Line()
		p.err(line, offset, endoffset-offset, "syntax error")
		return false, 0
	}

	return true, p.stack[len(p.stack)-2].value
}

const errSymbol = 17

func (p *Parser) recover() bool {
	if p.next.symbol == noToken {
		p.next.symbol = int32(p.lexer.Next())
		p.next.offset, p.next.endoffset = p.lexer.Pos()
	}
	if p.next.symbol == eoiToken {
		return false
	}
	e, _ := p.lexer.Pos()
	s := e
	for len(p.stack) > 0 && gotoState(p.stack[len(p.stack)-1].state, errSymbol) == -1 {
		// TODO cleanup
		p.stack = p.stack[:len(p.stack)-1]
		if len(p.stack) > 0 {
			s = p.stack[len(p.stack)-1].sym.offset
		}
	}
	if len(p.stack) > 0 {
		state := gotoState(p.stack[len(p.stack)-1].state, errSymbol)
		p.stack = append(p.stack, node{
			sym:   symbol{errSymbol, s, e},
			state: state,
		})
		return true
	}
	return false
}

func lalr(action, next int32) int32 {
	a := -action - 3
	for ; tmLalr[a] >= 0; a += 2 {
		if tmLalr[a] == next {
			break
		}
	}
	return tmLalr[a+1]
}

func gotoState(state int8, symbol int32) int8 {
	min := tmGoto[symbol]
	max := tmGoto[symbol+1] - 1

	for min <= max {
		e := (min + max) >> 1
		i := tmFrom[e]
		if i == state {
			return tmTo[e]
		} else if i < state {
			min = e + 1
		} else {
			max = e - 1
		}
	}
	return -1
}

func (p *Parser) lookahead(start, end int8) bool {
	var lexer Lexer = *p.lexer
	lexer.err = IgnoreErrorsHandler

	var allocated [64]node
	state := start
	stack := append(allocated[:0], node{state: state})
	next := p.next.symbol

	for state != end {
		action := tmAction[state]
		if action < -2 {
			// Lookahead is needed.
			if next == noToken {
				next = int32(lexer.Next())
			}
			action = lalr(action, next)
		}

		if action >= 0 {
			// Reduce.
			rule := action
			ln := int(tmRuleLen[rule])

			var node node
			node.sym.symbol = tmRuleSymbol[rule]
			stack = stack[:len(stack)-ln]
			state = gotoState(stack[len(stack)-1].state, node.sym.symbol)
			node.state = state
			stack = append(stack, node)

		} else if action == -1 {
			// Shift.
			if next == noToken {
				next = int32(lexer.Next())
			}
			state = gotoState(state, next)
			stack = append(stack, node{
				sym:   symbol{symbol: next},
				state: state,
			})
			if state != -1 && next != eoiToken {
				next = noToken
			}
		}

		if action == -2 || state == -1 {
			break
		}
	}

	return state == end
}

func (p *Parser) applyRule(rule int32, node *node, rhs []node) {
	switch rule {
	case 32:
		if p.lookahead(0, 42) /* EmptyObject */ {
			node.sym.symbol = 22 /* lookahead_EmptyObject */;
		} else {
			node.sym.symbol = 24 /* lookahead_notEmptyObject */;
		}
		return
	}
	nt := ruleNodeType[rule]
	if nt == 0 {
		return
	}
	p.listener.Node(nt, node.sym.offset, node.sym.endoffset)
}
