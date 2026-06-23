package lua

import "fmt"

type parser struct {
	toks []token
	pos  int
}

// Parse scans and parses src into a Chunk. It returns an error on malformed or
// unsupported input; callers fall back to Minify in that case.
func Parse(src string) (*Chunk, error) {
	raw := lex(src)
	toks := make([]token, 0, len(raw))
	for _, t := range raw {
		if t.kind != tkComment {
			toks = append(toks, t)
		}
	}
	p := &parser{toks: toks}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if !p.atEnd() {
		return nil, p.errf("unexpected %q", p.cur().text)
	}
	return &Chunk{Body: body}, nil
}

func (p *parser) atEnd() bool { return p.pos >= len(p.toks) }

func (p *parser) cur() token {
	if p.pos < len(p.toks) {
		return p.toks[p.pos]
	}
	return token{tkOp, ""}
}

func (p *parser) peekText() string { return p.cur().text }

func (p *parser) advance() token {
	t := p.cur()
	p.pos++
	return t
}

func (p *parser) isText(s string) bool { return !p.atEnd() && p.toks[p.pos].text == s }

func (p *parser) accept(s string) bool {
	if p.isText(s) {
		p.pos++
		return true
	}
	return false
}

func (p *parser) expect(s string) error {
	if p.accept(s) {
		return nil
	}
	return p.errf("expected %q, got %q", s, p.peekText())
}

func (p *parser) errf(format string, a ...any) error {
	return fmt.Errorf("lua parse: "+format, a...)
}

// keyword set, used to distinguish names from keywords.
var luaKeywords = map[string]bool{
	"and": true, "break": true, "do": true, "else": true, "elseif": true,
	"end": true, "false": true, "for": true, "function": true, "goto": true,
	"if": true, "in": true, "local": true, "nil": true, "not": true, "or": true,
	"repeat": true, "return": true, "then": true, "true": true, "until": true,
	"while": true, "continue": true,
}

func (p *parser) isName() bool {
	t := p.cur()
	return t.kind == tkName && !luaKeywords[t.text]
}

// peekAt returns the token at offset positions ahead of p.pos, or a zero token
// if out of bounds. Used for lookahead without advancing.
func (p *parser) peekAt(offset int) token {
	idx := p.pos + offset
	if idx < len(p.toks) {
		return p.toks[idx]
	}
	return token{tkOp, ""}
}

// ---- Expressions (precedence climbing) ----

// binPrec returns left/right binding powers for a binary operator, or ok=false.
func binPrec(op string) (left, right int, ok bool) {
	switch op {
	case "or":
		return 1, 2, true
	case "and":
		return 3, 4, true
	case "<", ">", "<=", ">=", "~=", "==":
		return 5, 6, true
	case "|":
		return 7, 8, true
	case "~":
		return 9, 10, true
	case "&":
		return 11, 12, true
	case "<<", ">>":
		return 13, 14, true
	case "..": // right associative
		return 18, 17, true
	case "+", "-":
		return 19, 20, true
	case "*", "/", "//", "%":
		return 21, 22, true
	case "^": // right associative
		return 28, 27, true
	}
	return 0, 0, false
}

const unaryPrec = 25 // between * / and ^

func (p *parser) parseExpr() (Expr, error) { return p.parseBinExpr(0) }

func (p *parser) parseBinExpr(minBP int) (Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for {
		op := p.peekText()
		lbp, rbp, ok := binPrec(op)
		if !ok || lbp < minBP {
			break
		}
		p.advance()
		right, err := p.parseBinExpr(rbp)
		if err != nil {
			return nil, err
		}
		left = &BinExpr{Op: op, L: left, R: right}
	}
	return left, nil
}

func (p *parser) parseUnary() (Expr, error) {
	op := p.peekText()
	if op == "not" || op == "-" || op == "#" || op == "~" {
		p.advance()
		e, err := p.parseBinExpr(unaryPrec)
		if err != nil {
			return nil, err
		}
		return &UnExpr{Op: op, E: e}, nil
	}
	return p.parseSuffixed()
}

// parseSuffixed parses a primary expression followed by any chain of
// .field, :method(args), [key], (args), and string/table call sugar.
func (p *parser) parseSuffixed() (Expr, error) {
	e, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for {
		switch {
		case p.accept("."):
			if !p.isName() {
				return nil, p.errf("expected field name after '.'")
			}
			e = &IndexExpr{Obj: e, Field: p.advance().text}
		case p.accept(":"):
			if !p.isName() {
				return nil, p.errf("expected method name after ':'")
			}
			name := p.advance().text
			args, err := p.parseCallArgs()
			if err != nil {
				return nil, err
			}
			e = &CallExpr{Fn: &IndexExpr{Obj: e, Field: name, IsMethod: true}, Args: args}
		case p.accept("["):
			k, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if err := p.expect("]"); err != nil {
				return nil, err
			}
			e = &IndexExpr{Obj: e, Key: k}
		case p.isText("("), p.cur().kind == tkString, p.isText("{"):
			args, err := p.parseCallArgs()
			if err != nil {
				return nil, err
			}
			e = &CallExpr{Fn: e, Args: args}
		default:
			return e, nil
		}
	}
}

// parseCallArgs handles (a, b), "str", and {table} call forms.
func (p *parser) parseCallArgs() ([]Expr, error) {
	if p.cur().kind == tkString {
		return []Expr{&StringExpr{Text: p.advance().text}}, nil
	}
	if p.isText("{") {
		tbl, err := p.parseTable()
		if err != nil {
			return nil, err
		}
		return []Expr{tbl}, nil
	}
	if err := p.expect("("); err != nil {
		return nil, err
	}
	var args []Expr
	if !p.isText(")") {
		for {
			a, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			args = append(args, a)
			if !p.accept(",") {
				break
			}
		}
	}
	if err := p.expect(")"); err != nil {
		return nil, err
	}
	return args, nil
}

func (p *parser) parsePrimary() (Expr, error) {
	t := p.cur()
	switch {
	case t.kind == tkNumber:
		p.advance()
		return &NumberExpr{Text: t.text}, nil
	case t.kind == tkString:
		p.advance()
		return &StringExpr{Text: t.text}, nil
	case t.text == "nil":
		p.advance()
		return &NilExpr{}, nil
	case t.text == "true":
		p.advance()
		return &BoolExpr{Val: true}, nil
	case t.text == "false":
		p.advance()
		return &BoolExpr{Val: false}, nil
	case t.text == "...":
		p.advance()
		return &VarargExpr{}, nil
	case t.text == "function":
		p.advance()
		return p.parseFuncBody()
	case t.text == "{":
		return p.parseTable()
	case t.text == "(":
		p.advance()
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expect(")"); err != nil {
			return nil, err
		}
		return &ParenExpr{E: e}, nil
	case p.isName():
		p.advance()
		return &NameExpr{Name: t.text}, nil
	}
	return nil, p.errf("unexpected token %q", t.text)
}

// parseFuncBody parses the part after the 'function' keyword: (params) body end.
func (p *parser) parseFuncBody() (*FuncExpr, error) {
	if err := p.expect("("); err != nil {
		return nil, err
	}
	fn := &FuncExpr{}
	if !p.isText(")") {
		for {
			if p.accept("...") {
				fn.IsVararg = true
				break
			}
			if !p.isName() {
				return nil, p.errf("expected parameter name")
			}
			fn.Params = append(fn.Params, &NameExpr{Name: p.advance().text})
			p.skipTypeAnnotation() // Luau: param: T
			if !p.accept(",") {
				break
			}
		}
	}
	if err := p.expect(")"); err != nil {
		return nil, err
	}
	p.skipReturnTypeAnnotation() // Luau: ): T
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if err := p.expect("end"); err != nil {
		return nil, err
	}
	fn.Body = body
	return fn, nil
}

func (p *parser) parseTable() (Expr, error) {
	if err := p.expect("{"); err != nil {
		return nil, err
	}
	tbl := &TableExpr{}
	for !p.isText("}") && !p.atEnd() {
		var f TableField
		switch {
		case p.accept("["):
			k, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if err := p.expect("]"); err != nil {
				return nil, err
			}
			if err := p.expect("="); err != nil {
				return nil, err
			}
			v, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			f.Key, f.Value = k, v
		case p.isName() && p.peekAt(1).text == "=":
			f.KeyName = p.advance().text
			p.advance() // '='
			v, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			f.Value = v
		default:
			v, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			f.Value = v
		}
		tbl.Fields = append(tbl.Fields, f)
		if !p.accept(",") && !p.accept(";") {
			break
		}
	}
	if err := p.expect("}"); err != nil {
		return nil, err
	}
	return tbl, nil
}

// parseBlock is defined in Task 5. For Task 3 a temporary version supports only
// a single `return <expr>` so expression tests can run; Task 5 replaces it.
func (p *parser) parseBlock() ([]Stat, error) {
	if p.accept("return") {
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return []Stat{&ReturnStat{Values: []Expr{e}}}, nil
	}
	return nil, nil
}

// Type-annotation skippers are implemented fully in Task 6. Temporary no-ops:
func (p *parser) skipTypeAnnotation()       {}
func (p *parser) skipReturnTypeAnnotation() {}
