package grammar

import (
	"github.com/inspirer/textmapper/tm-go/lex"
	"github.com/inspirer/textmapper/tm-go/status"
	"github.com/inspirer/textmapper/tm-go/util/container"
	"github.com/inspirer/textmapper/tm-parsers/tm/ast"
	"sort"
	"strconv"
)

// Compile validates and compiles grammar files.
func Compile(file ast.File) (*Grammar, error) {
	targetLang, _ := file.Header().Target()
	c := &compiler{
		file: file,
		out: &Grammar{
			Name:       file.Header().Name().Text(),
			TargetLang: targetLang.Text(),
			Lexer:      &Lexer{},
			Options: &Options{
				TokenLine: true,
			},
		},
		syms:       make(map[string]int),
		symAction:  make(map[int]int),
		codeAction: make(map[symCode]int),
	}
	c.compileLexer()
	return c.out, c.s.Err()
}

type compiler struct {
	file ast.File
	out  *Grammar
	s    status.Status

	syms        map[string]int
	conds       map[string]int
	inclusiveSC []int
	patterns    []*patterns // to keep track of unused patterns
	classRules  []*lex.Rule
	rules       []*lex.Rule
	symAction   map[int]int
	codeAction  map[symCode]int
}

type symCode struct {
	code  string
	space bool
	sym   int
}

func (c *compiler) compileLexer() {
	c.addToken(Eoi, ast.RawType{}, nil)
	c.out.InvalidToken = c.addToken(InvalidToken, ast.RawType{}, nil)

	c.collectStartConds()
	lexer, _ := c.file.Lexer()
	c.traverseLexer(lexer.LexerPart(), c.inclusiveSC, nil /*parent patterns*/)
	c.resolveClasses()

	var err error
	c.out.Tables, err = lex.Compile(c.rules)
	c.s.AddError(err)

	for _, p := range c.patterns {
		for name, unused := range p.unused {
			c.errorf(unused, "unused pattern '%v'", name)
		}
	}
}

func (c *compiler) collectStartConds() {
	conds := make(map[string]bool)
	var names []string

	insert := func(n *ast.Node, excl bool) {
		name := n.Text()
		if _, exists := conds[name]; exists {
			c.errorf(n, "redeclaration of %v", name)
			return
		}
		conds[name] = excl
		names = append(names, name)
	}
	lexer, _ := c.file.Lexer()
	for _, p := range lexer.LexerPart() {
		switch p := p.(type) {
		case *ast.ExclusiveStartConds:
			for _, s := range p.States() {
				insert(s.Node, true)
			}
		case *ast.InclusiveStartConds:
			for _, s := range p.States() {
				insert(s.Node, false)
			}
		}
	}

	if _, exists := conds[lex.Initial]; !exists {
		names = append(names, lex.Initial)
		copy(names[1:], names)
		names[0] = lex.Initial
	}

	c.out.StartConditions = names
	c.conds = make(map[string]int)
	for i, name := range names {
		c.conds[name] = i
		if !conds[name] {
			c.inclusiveSC = append(c.inclusiveSC, i)
		}
	}
}

func (c *compiler) resolveSC(sc ast.StartConditions) []int {
	var ret []int
	refs := sc.Stateref()
	if len(refs) == 0 {
		// <*>
		for i := range c.out.StartConditions {
			ret = append(ret, i)
		}
		return ret
	}
	for _, ref := range refs {
		name := ref.Name().Text()
		if val, ok := c.conds[name]; ok {
			ret = append(ret, val)
			continue
		}
		c.errorf(ref.Name(), "unresolved reference %v", name)
	}
	sort.Ints(ret)
	return ret
}

func (c *compiler) addToken(name string, t ast.RawType, n status.SourceNode) int {
	var rawType string
	if t.IsValid() {
		rawType = t.Text()
	}
	if i, exists := c.syms[name]; exists {
		sym := c.out.Syms[i]
		if sym.Type != rawType {
			anchor := n
			if t.IsValid() {
				anchor = t
			}
			c.errorf(anchor, "terminal type redeclaration for %v, was %v", name, sym.PrettyType())
		}
		return sym.Index
	}

	sym := Symbol{
		Index:  len(c.syms),
		Type:   rawType,
		Name:   name,
		Origin: n,
	}
	c.syms[name] = sym.Index
	c.out.Syms = append(c.out.Syms, sym)
	return sym.Index
}

func (c *compiler) addLexerAction(cmd ast.Command, space ast.LexemeAttribute, sym int) int {
	if !cmd.IsValid() && !space.IsValid() {
		if a, ok := c.symAction[sym]; ok {
			return a
		}
		a := len(c.out.RuleToken)
		c.out.RuleToken = append(c.out.RuleToken, sym)
		c.symAction[sym] = a
		return a
	}

	key := symCode{cmd.Text(), space.IsValid(), sym}
	if a, ok := c.codeAction[key]; ok {
		return a
	}
	a := len(c.out.RuleToken)
	c.out.RuleToken = append(c.out.RuleToken, sym)
	c.codeAction[key] = a

	act := SemanticAction{Action: a, Code: key.code, Space: space.IsValid()}
	if cmd.IsValid() {
		act.Origin = cmd
	} else {
		act.Origin = space
	}
	c.out.Actions = append(c.out.Actions, act)
	return a
}

func (c *compiler) traverseLexer(parts []ast.LexerPart, defaultSCs []int, p *patterns) {
	inClause := p != nil
	ps := &patterns{
		parent: p,
		set:    make(map[string]*lex.Regexp),
		unused: make(map[string]status.SourceNode),
	}
	c.patterns = append(c.patterns, ps)

	for _, p := range parts {
		switch p := p.(type) {
		case *ast.Lexeme:
			rawType, _ := p.RawType()
			tok := c.addToken(p.Name().Text(), rawType, p.Name())
			pat, ok := p.Pattern()
			if !ok {
				break
			}

			re, err := parsePattern(pat)
			c.s.AddError(err)
			rule := &lex.Rule{
				RE:              re,
				StartConditions: defaultSCs,
				Resolver:        ps,
				Origin:          p,
				OriginName:      p.Name().Text(),
			}

			if prio, ok := p.Priority(); ok {
				rule.Precedence, _ = strconv.Atoi(prio.Text())
			}
			if sc, ok := p.StartConditions(); ok {
				rule.StartConditions = c.resolveSC(sc)
			}

			var class bool
			var space ast.LexemeAttribute
			if attrs, ok := p.Attrs(); ok {
				switch name := attrs.LexemeAttribute().Text(); name {
				case "class":
					class = true
				case "space":
					space = attrs.LexemeAttribute()
				default:
					c.errorf(attrs.LexemeAttribute(), "unsupported attribute")
				}
			}

			cmd, _ := p.Command()
			rule.Action = c.addLexerAction(cmd, space, tok)
			if class {
				c.classRules = append(c.classRules, rule)
			} else {
				c.rules = append(c.rules, rule)
			}
		case *ast.NamedPattern:
			c.s.AddError(ps.add(p))
		case *ast.StartConditionsScope:
			newDefaults := c.resolveSC(p.StartConditions())
			c.traverseLexer(p.LexerPart(), newDefaults, ps)
		case *ast.SyntaxProblem, *ast.DirectiveBrackets:
			c.errorf(p.TmNode(), "syntax error")
		case *ast.ExclusiveStartConds, *ast.InclusiveStartConds:
			if inClause {
				// %s and %x are not allowed inside start condition clauses.
				c.errorf(p.TmNode(), "syntax error")
			}
		}
	}
}

func (c *compiler) resolveClasses() {
	if len(c.classRules) == 0 {
		return
	}

	var rewritten []*lex.Rule
	for index, r := range c.classRules {
		fork := new(lex.Rule)
		*fork = *r
		fork.Action = index
		rewritten = append(rewritten, fork)
	}
	tables, err := lex.Compile(rewritten)
	c.s.AddError(err)
	if err != nil {
		// Pretend that these class rules do not exist in the grammar and keep going.
		return
	}

	// Pre-create class actions.
	for _, r := range c.classRules {
		ca := ClassAction{
			Action: r.Action,
			Custom: make(map[string]int),
		}
		c.out.ClassActions = append(c.out.ClassActions, ca)
	}

	out := c.rules[:0]
	for _, r := range c.rules {
		val, isConst := r.RE.Constant()
		if !isConst {
			out = append(out, r)
			continue
		}

		classRule := -1
		for _, start := range r.StartConditions {
			size, result := tables.Scan(start, val)
			if size == len(val) && result >= 0 {
				classRule = result
				break
			}
		}
		if classRule == -1 {
			out = append(out, r)
			continue
		}
		class := c.classRules[classRule]

		if !container.SliceEqual(class.StartConditions, r.StartConditions) {
			c.errorf(r.Origin, "%v must be applicable in the same set of start conditions as %v", r.OriginName, class.OriginName)

			// Fixing the problem for now and keep going.
			r.StartConditions = class.StartConditions
		}

		// Move the rule under its class rule.
		c.out.ClassActions[classRule].Custom[val] = r.Action
	}

	for i, r := range c.classRules {
		if len(c.out.ClassActions[i].Custom) == 0 {
			c.errorf(r.Origin, "class rule without specializations '%v'", r.OriginName)
		}
	}

	c.rules = append(out, c.classRules...)
	c.classRules = nil
}

func (c *compiler) errorf(n status.SourceNode, format string, a ...interface{}) {
	c.s.Errorf(n, format, a...)
}

type patterns struct {
	parent *patterns
	set    map[string]*lex.Regexp
	unused map[string]status.SourceNode
}

func (p *patterns) Resolve(name string) *lex.Regexp {
	if v, ok := p.set[name]; ok {
		delete(p.unused, name)
		return v
	}
	if p.parent != nil {
		return p.parent.Resolve(name)
	}
	return nil
}

var emptyRE = lex.MustParse("")

func (p *patterns) add(np *ast.NamedPattern) error {
	name := np.Name().Text()
	if _, exists := p.set[name]; exists {
		return status.Errorf(np.Name(), "redeclaration of %v", name)
	}

	re, err := parsePattern(np.Pattern())
	p.set[name] = re
	p.unused[name] = np.Name()
	return err
}

func parsePattern(p ast.Pattern) (*lex.Regexp, error) {
	text := p.Text()
	text = text[1 : len(text)-1]
	re, err := lex.ParseRegexp(text)
	if err != nil {
		rng := p.SourceRange()
		err := err.(lex.ParseError)
		if err.Offset <= err.EndOffset && err.EndOffset <= len(text) && err.Offset < len(text) {
			if err.Offset < err.EndOffset {
				rng.EndOffset = rng.Offset + err.EndOffset + 1
			} else {
				rng.EndOffset--
			}
			rng.Offset += err.Offset + 1
			rng.Column += err.Offset + 1
		}
		return emptyRE, &status.Error{Origin: rng, Msg: err.Error()}
	}
	return re, nil
}
