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
	case t.text == "if":
		return p.parseIfExpr()
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

// blockEnd reports whether the current token closes a block.
func (p *parser) blockEnd() bool {
	switch p.peekText() {
	case "", "end", "else", "elseif", "until":
		return true
	}
	return false
}

func (p *parser) parseBlock() ([]Stat, error) {
	var body []Stat
	for !p.blockEnd() {
		if p.accept(";") {
			continue
		}
		if p.peekText() == "return" {
			ret, err := p.parseReturn()
			if err != nil {
				return nil, err
			}
			body = append(body, ret)
			break // return must be last in a block
		}
		s, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if s != nil {
			body = append(body, s)
		}
	}
	return body, nil
}

func (p *parser) parseReturn() (Stat, error) {
	p.advance() // 'return'
	ret := &ReturnStat{}
	if !p.blockEnd() && !p.isText(";") {
		for {
			e, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			ret.Values = append(ret.Values, e)
			if !p.accept(",") {
				break
			}
		}
	}
	p.accept(";")
	return ret, nil
}

func (p *parser) parseStatement() (Stat, error) {
	// Luau contextual: `type X = ...` and `export type X = ...`
	// 'type' and 'export' are not Lua keywords; detect them by lookahead.
	if (p.isText("export") && p.pos+1 < len(p.toks) && p.toks[p.pos+1].text == "type") ||
		(p.isText("type") && p.isName2(p.pos+1)) {
		return p.parseTypeAlias()
	}
	switch p.peekText() {
	case "local":
		return p.parseLocal()
	case "do":
		p.advance()
		body, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		if err := p.expect("end"); err != nil {
			return nil, err
		}
		return &DoStat{Body: body}, nil
	case "while":
		p.advance()
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expect("do"); err != nil {
			return nil, err
		}
		body, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		if err := p.expect("end"); err != nil {
			return nil, err
		}
		return &WhileStat{Cond: cond, Body: body}, nil
	case "repeat":
		p.advance()
		body, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		if err := p.expect("until"); err != nil {
			return nil, err
		}
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &RepeatStat{Body: body, Cond: cond}, nil
	case "if":
		return p.parseIf()
	case "for":
		return p.parseFor()
	case "function":
		return p.parseFuncStat()
	case "break":
		p.advance()
		return &BreakStat{}, nil
	case "continue":
		p.advance()
		return &ContinueStat{}, nil
	case "goto":
		p.advance()
		if !p.isName() {
			return nil, p.errf("expected label after goto")
		}
		return &GotoStat{Label: p.advance().text}, nil
	case "::":
		p.advance()
		if !p.isName() {
			return nil, p.errf("expected label name")
		}
		name := p.advance().text
		if err := p.expect("::"); err != nil {
			return nil, err
		}
		return &LabelStat{Name: name}, nil
	}
	// Expression statement: call or assignment.
	return p.parseExprStatement()
}

func (p *parser) parseLocal() (Stat, error) {
	p.advance() // 'local'
	if p.peekText() == "function" {
		p.advance()
		if !p.isName() {
			return nil, p.errf("expected function name")
		}
		name := &NameExpr{Name: p.advance().text}
		fn, err := p.parseFuncBody()
		if err != nil {
			return nil, err
		}
		return &LocalFuncStat{Name: name, Func: fn}, nil
	}
	ls := &LocalStat{}
	for {
		if !p.isName() {
			return nil, p.errf("expected local name, got %q", p.peekText())
		}
		ls.Names = append(ls.Names, &NameExpr{Name: p.advance().text})
		p.skipTypeAnnotation()
		if !p.accept(",") {
			break
		}
	}
	if p.accept("=") {
		for {
			e, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			ls.Values = append(ls.Values, e)
			if !p.accept(",") {
				break
			}
		}
	}
	return ls, nil
}

func (p *parser) parseIf() (Stat, error) {
	st := &IfStat{}
	p.advance() // 'if'
	for {
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expect("then"); err != nil {
			return nil, err
		}
		blk, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		st.Conds = append(st.Conds, cond)
		st.Blocks = append(st.Blocks, blk)
		if !p.accept("elseif") {
			break
		}
	}
	if p.accept("else") {
		blk, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		st.HasElse = true
		st.Else = blk
	}
	if err := p.expect("end"); err != nil {
		return nil, err
	}
	return st, nil
}

func (p *parser) parseFor() (Stat, error) {
	p.advance() // 'for'
	if !p.isName() {
		return nil, p.errf("expected loop variable")
	}
	first := &NameExpr{Name: p.advance().text}
	p.skipTypeAnnotation()
	if p.accept("=") {
		start, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expect(","); err != nil {
			return nil, err
		}
		stop, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		var step Expr
		if p.accept(",") {
			step, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
		}
		if err := p.expect("do"); err != nil {
			return nil, err
		}
		body, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		if err := p.expect("end"); err != nil {
			return nil, err
		}
		return &NumericForStat{Var: first, Start: start, Stop: stop, Step: step, Body: body}, nil
	}
	// generic for
	vars := []*NameExpr{first}
	for p.accept(",") {
		if !p.isName() {
			return nil, p.errf("expected loop variable")
		}
		vars = append(vars, &NameExpr{Name: p.advance().text})
		p.skipTypeAnnotation()
	}
	if err := p.expect("in"); err != nil {
		return nil, err
	}
	var exprs []Expr
	for {
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, e)
		if !p.accept(",") {
			break
		}
	}
	if err := p.expect("do"); err != nil {
		return nil, err
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if err := p.expect("end"); err != nil {
		return nil, err
	}
	return &GenericForStat{Vars: vars, Exprs: exprs, Body: body}, nil
}

func (p *parser) parseFuncStat() (Stat, error) {
	p.advance() // 'function'
	if !p.isName() {
		return nil, p.errf("expected function name")
	}
	var target Expr = &NameExpr{Name: p.advance().text}
	isMethod := false
	for {
		if p.accept(".") {
			if !p.isName() {
				return nil, p.errf("expected field name")
			}
			target = &IndexExpr{Obj: target, Field: p.advance().text}
			continue
		}
		if p.accept(":") {
			if !p.isName() {
				return nil, p.errf("expected method name")
			}
			target = &IndexExpr{Obj: target, Field: p.advance().text, IsMethod: true}
			isMethod = true
		}
		break
	}
	fn, err := p.parseFuncBody()
	if err != nil {
		return nil, err
	}
	_ = isMethod // self injection deferred to Task 7 scope resolver
	return &FuncStat{Target: target, IsMethod: isMethod, Func: fn}, nil
}

func (p *parser) parseExprStatement() (Stat, error) {
	first, err := p.parseSuffixed()
	if err != nil {
		return nil, err
	}
	// Assignment? (=, or a Luau compound op, or a comma for multi-assign)
	if p.isText(",") || isAssignOp(p.peekText()) {
		targets := []Expr{first}
		for p.accept(",") {
			t, err := p.parseSuffixed()
			if err != nil {
				return nil, err
			}
			targets = append(targets, t)
		}
		op := p.peekText()
		if !isAssignOp(op) {
			return nil, p.errf("expected assignment operator, got %q", op)
		}
		p.advance()
		var vals []Expr
		for {
			v, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			vals = append(vals, v)
			if !p.accept(",") {
				break
			}
		}
		return &AssignStat{Targets: targets, Op: op, Values: vals}, nil
	}
	// Otherwise it must be a call statement.
	if _, ok := first.(*CallExpr); !ok {
		return nil, p.errf("expected statement, got bare expression")
	}
	return &CallStat{Call: first}, nil
}

func isAssignOp(op string) bool {
	switch op {
	case "=", "+=", "-=", "*=", "/=", "//=", "%=", "^=", "..=":
		return true
	}
	return false
}

// skipTypeAnnotation consumes an optional ': <type>' after a name/param.
func (p *parser) skipTypeAnnotation() {
	if p.isText(":") {
		p.advance()
		p.skipType()
	}
}

// skipReturnTypeAnnotation consumes an optional ': <type>' after a ')' params list.
func (p *parser) skipReturnTypeAnnotation() {
	if p.isText(":") {
		p.advance()
		p.skipType()
	}
}

// skipType consumes a Luau type expression after ':' or '='. It consumes at
// least one type atom, then continues across type connectors (| & -> ?) and
// balanced brackets, stopping at a token that ends the type at depth 0. It is
// permissive: on an exotic type it may over-consume and cause Parse to error,
// which is the intended safety net (the obfuscator falls back to minify).
func (p *parser) skipType() {
	p.skipTypeAtom()
	for {
		switch {
		case p.isText("?"):
			p.advance() // postfix optional, no atom follows
		case p.isText("|") || p.isText("&"):
			p.advance()
			p.skipTypeAtom()
		case p.isText("-") && p.peekAt(1).text == ">":
			p.advance() // '-'
			p.advance() // '>'
			p.skipTypeAtom()
		default:
			return
		}
	}
}

// skipTypeAtom consumes one type atom: a balanced (...) or {...} group, or a
// (possibly dotted, possibly generic) name / literal / typeof(...).
func (p *parser) skipTypeAtom() {
	switch {
	case p.isText("("):
		p.skipBalanced("(", ")")
	case p.isText("{"):
		p.skipBalanced("{", "}")
	default:
		if p.atEnd() {
			return
		}
		p.advance() // name / nil / true / false / string / number / typeof / ...
		for p.isText(".") { // dotted module type a.b.C
			p.advance()
			if !p.atEnd() {
				p.advance()
			}
		}
		if p.isText("(") { // typeof(expr)
			p.skipBalanced("(", ")")
		}
		if p.isText("<") { // generic args <...>
			p.skipBalanced("<", ">")
		}
	}
}

// skipBalanced consumes a run delimited by open/close, honoring nesting.
func (p *parser) skipBalanced(open, close string) {
	if !p.accept(open) {
		return
	}
	depth := 1
	for !p.atEnd() && depth > 0 {
		t := p.peekText()
		if t == open {
			depth++
		} else if t == close {
			depth--
		}
		p.advance()
	}
}

// isName2 checks whether the token at absolute index i is a non-keyword name.
func (p *parser) isName2(i int) bool {
	return i < len(p.toks) && p.toks[i].kind == tkName && !luaKeywords[p.toks[i].text]
}

// parseTypeAlias consumes a Luau `type X = ...` or `export type X = ...`
// statement and returns a TypeAliasStat with Raw == nil so the printer emits nothing.
func (p *parser) parseTypeAlias() (Stat, error) {
	p.accept("export")  // optional 'export'
	p.advance()         // 'type'
	if !p.isName() {
		return nil, p.errf("expected type name")
	}
	p.advance() // type name
	// optional generics <T, U>
	if p.isText("<") {
		p.skipBalanced("<", ">")
	}
	if err := p.expect("="); err != nil {
		return nil, err
	}
	p.skipType()
	return &TypeAliasStat{Raw: nil}, nil
}

// parseIfExpr parses a Luau if-expression: `if cond then a [elseif c then b]* else z`.
func (p *parser) parseIfExpr() (Expr, error) {
	p.advance() // 'if'
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if err := p.expect("then"); err != nil {
		return nil, err
	}
	thenE, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	ie := &IfExpr{Cond: cond, Then: thenE}
	for p.accept("elseif") {
		c, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expect("then"); err != nil {
			return nil, err
		}
		v, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ie.ElifConds = append(ie.ElifConds, c)
		ie.ElifThen = append(ie.ElifThen, v)
	}
	if err := p.expect("else"); err != nil {
		return nil, err
	}
	elseE, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	ie.Else = elseE
	return ie, nil
}
