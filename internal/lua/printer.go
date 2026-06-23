package lua

import "strings"

type emitter struct{ toks []token }

func (e *emitter) name(s string) { e.toks = append(e.toks, token{tkName, s}) }
func (e *emitter) op(s string)   { e.toks = append(e.toks, token{tkOp, s}) }
func (e *emitter) num(s string)  { e.toks = append(e.toks, token{tkNumber, s}) }
func (e *emitter) str(s string)  { e.toks = append(e.toks, token{tkString, s}) }

func (e *emitter) string() string {
	var b strings.Builder
	var prev *token
	for i := range e.toks {
		t := &e.toks[i]
		if t.text == "" {
			continue
		}
		if prev != nil && needsSpace(*prev, *t) {
			b.WriteByte(' ')
		}
		b.WriteString(t.text)
		prev = t
	}
	return b.String()
}

// Print renders the chunk as minified Lua.
func (c *Chunk) Print() string {
	e := &emitter{}
	e.block(c.Body)
	return e.string()
}

func (e *emitter) block(stats []Stat) {
	for i, s := range stats {
		e.stat(s)
		// A ';' separator is only needed to disambiguate a statement that ends
		// with an expression from a following '(' (ambiguous call). Emit one
		// defensively between two statements when the next begins with '('.
		if i+1 < len(stats) && beginsWithParen(stats[i+1]) {
			e.op(";")
		}
	}
}

func beginsWithParen(s Stat) bool {
	switch v := s.(type) {
	case *CallStat:
		return exprBeginsWithParen(v.Call)
	case *AssignStat:
		return len(v.Targets) > 0 && exprBeginsWithParen(v.Targets[0])
	}
	return false
}

func exprBeginsWithParen(x Expr) bool {
	for {
		switch v := x.(type) {
		case *ParenExpr:
			return true
		case *CallExpr:
			x = v.Fn
		case *IndexExpr:
			x = v.Obj
		default:
			return false
		}
	}
}

func (e *emitter) stat(s Stat) {
	switch v := s.(type) {
	case *LocalStat:
		e.name("local")
		e.nameList(v.Names)
		if len(v.Values) > 0 {
			e.op("=")
			e.exprList(v.Values)
		}
	case *AssignStat:
		e.exprList(v.Targets)
		e.op(v.Op)
		e.exprList(v.Values)
	case *CallStat:
		e.expr(v.Call)
	case *DoStat:
		e.name("do")
		e.block(v.Body)
		e.name("end")
	case *WhileStat:
		e.name("while")
		e.expr(v.Cond)
		e.name("do")
		e.block(v.Body)
		e.name("end")
	case *RepeatStat:
		e.name("repeat")
		e.block(v.Body)
		e.name("until")
		e.expr(v.Cond)
	case *IfStat:
		for i, cond := range v.Conds {
			if i == 0 {
				e.name("if")
			} else {
				e.name("elseif")
			}
			e.expr(cond)
			e.name("then")
			e.block(v.Blocks[i])
		}
		if v.HasElse {
			e.name("else")
			e.block(v.Else)
		}
		e.name("end")
	case *NumericForStat:
		e.name("for")
		e.nameRef(v.Var)
		e.op("=")
		e.expr(v.Start)
		e.op(",")
		e.expr(v.Stop)
		if v.Step != nil {
			e.op(",")
			e.expr(v.Step)
		}
		e.name("do")
		e.block(v.Body)
		e.name("end")
	case *GenericForStat:
		e.name("for")
		e.nameList(v.Vars)
		e.name("in")
		e.exprList(v.Exprs)
		e.name("do")
		e.block(v.Body)
		e.name("end")
	case *FuncStat:
		e.name("function")
		e.expr(v.Target)
		e.funcBody(v.Func)
	case *LocalFuncStat:
		e.name("local")
		e.name("function")
		e.nameRef(v.Name)
		e.funcBody(v.Func)
	case *ReturnStat:
		e.name("return")
		e.exprList(v.Values)
	case *BreakStat:
		e.name("break")
	case *ContinueStat:
		e.name("continue")
	case *GotoStat:
		e.name("goto")
		e.name(v.Label)
	case *LabelStat:
		e.op("::")
		e.name(v.Name)
		e.op("::")
	case *TypeAliasStat:
		for _, tk := range v.Raw {
			e.toks = append(e.toks, token{tkName, tk}) // raw passthrough; spacing via needsSpace
		}
	}
}

func (e *emitter) funcBody(f *FuncExpr) {
	e.op("(")
	for i, pm := range f.Params {
		if i > 0 {
			e.op(",")
		}
		e.nameRef(pm)
	}
	if f.IsVararg {
		if len(f.Params) > 0 {
			e.op(",")
		}
		e.op("...")
	}
	e.op(")")
	e.block(f.Body)
	e.name("end")
}

func (e *emitter) nameList(ns []*NameExpr) {
	for i, n := range ns {
		if i > 0 {
			e.op(",")
		}
		e.nameRef(n)
	}
}

func (e *emitter) nameRef(n *NameExpr) {
	if n.Binding != nil && n.Binding.NewName != "" {
		e.name(n.Binding.NewName)
		return
	}
	e.name(n.Name)
}

func (e *emitter) exprList(xs []Expr) {
	for i, x := range xs {
		if i > 0 {
			e.op(",")
		}
		e.expr(x)
	}
}

func (e *emitter) expr(x Expr) {
	switch v := x.(type) {
	case *NameExpr:
		e.nameRef(v)
	case *NumberExpr:
		e.num(v.Text)
	case *StringExpr:
		e.str(v.Text)
	case *BoolExpr:
		if v.Val {
			e.name("true")
		} else {
			e.name("false")
		}
	case *NilExpr:
		e.name("nil")
	case *VarargExpr:
		e.op("...")
	case *ParenExpr:
		e.op("(")
		e.expr(v.E)
		e.op(")")
	case *IndexExpr:
		e.expr(v.Obj)
		if v.Key != nil {
			e.op("[")
			e.expr(v.Key)
			e.op("]")
		} else if v.IsMethod {
			e.op(":")
			e.name(v.Field)
		} else {
			e.op(".")
			e.name(v.Field)
		}
	case *CallExpr:
		e.expr(v.Fn)
		// method call already printed ':name' via Fn IndexExpr; args follow.
		e.op("(")
		e.exprList(v.Args)
		e.op(")")
	case *FuncExpr:
		e.name("function")
		e.funcBody(v)
	case *TableExpr:
		e.op("{")
		for i, f := range v.Fields {
			if i > 0 {
				e.op(",")
			}
			switch {
			case f.Key != nil:
				e.op("[")
				e.expr(f.Key)
				e.op("]")
				e.op("=")
				e.expr(f.Value)
			case f.KeyName != "":
				e.name(f.KeyName)
				e.op("=")
				e.expr(f.Value)
			default:
				e.expr(f.Value)
			}
		}
		e.op("}")
	case *BinExpr:
		e.expr(v.L)
		// Operators that are words (and/or) rely on needsSpace; symbolic ops too.
		if isWordOp(v.Op) {
			e.name(v.Op)
		} else {
			e.op(v.Op)
		}
		e.expr(v.R)
	case *UnExpr:
		if v.Op == "not" {
			e.name(v.Op)
		} else {
			e.op(v.Op)
		}
		e.expr(v.E)
	case *IfExpr:
		e.name("if")
		e.expr(v.Cond)
		e.name("then")
		e.expr(v.Then)
		for i := range v.ElifConds {
			e.name("elseif")
			e.expr(v.ElifConds[i])
			e.name("then")
			e.expr(v.ElifThen[i])
		}
		e.name("else")
		e.expr(v.Else)
	}
}

func isWordOp(op string) bool { return op == "and" || op == "or" }
