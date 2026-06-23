package lua

// Node is any AST node.
type Node interface{ node() }

// Stat is a statement node.
type Stat interface {
	Node
	stat()
}

// Expr is an expression node.
type Expr interface {
	Node
	expr()
}

// Chunk is a parsed unit (a file or module body).
type Chunk struct{ Body []Stat }

func (*Chunk) node() {}

// ---- Statements ----

type LocalStat struct {
	Names   []*NameExpr // declared local names (renameable)
	Attribs []string    // Luau attributes per name ("" if none), e.g. "const"; len == len(Names) or 0
	Values  []Expr      // RHS expressions (may be empty)
}

type AssignStat struct {
	Targets []Expr // lvalues (NameExpr or IndexExpr)
	Op      string // "=", "+=", "-=", "*=", "/=", "//=", "%=", "^=", "..="
	Values  []Expr
}

type CallStat struct{ Call Expr } // a function/method call used as a statement

type DoStat struct{ Body []Stat }

type WhileStat struct {
	Cond Expr
	Body []Stat
}

type RepeatStat struct {
	Body []Stat
	Cond Expr // sees locals declared in Body
}

type IfStat struct {
	Conds   []Expr   // len == len(Blocks) for if/elseif; else has no cond
	Blocks  [][]Stat // one per cond, plus optionally one trailing else block
	HasElse bool
	Else    []Stat
}

type NumericForStat struct {
	Var         *NameExpr // renameable, scoped to Body
	Start, Stop Expr
	Step        Expr // nil if absent
	Body        []Stat
}

type GenericForStat struct {
	Vars  []*NameExpr // renameable, scoped to Body
	Exprs []Expr
	Body  []Stat
}

type FuncStat struct {
	// function a.b.c:d() ... end  — Target is the dotted/colon path, not renameable
	// (it is a field assignment). Name resolution treats the leading name as a use.
	Target   Expr // NameExpr or IndexExpr chain
	IsMethod bool // true if defined with ':'
	Func     *FuncExpr
}

type LocalFuncStat struct {
	Name *NameExpr // renameable; visible inside Func (for recursion)
	Func *FuncExpr
}

type ReturnStat struct{ Values []Expr }

type BreakStat struct{}

type ContinueStat struct{} // Luau

type GotoStat struct{ Label string }

type LabelStat struct{ Name string }

// TypeAliasStat is a Luau `type X = ...` / `export type X = ...`. The whole
// statement is preserved as raw tokens and re-emitted verbatim minus spacing,
// because Phase 1 does not transform types. Renaming skips it entirely.
type TypeAliasStat struct{ Raw []string } // token texts in order

func (*LocalStat) node()      {}
func (*LocalStat) stat()      {}
func (*AssignStat) node()     {}
func (*AssignStat) stat()     {}
func (*CallStat) node()       {}
func (*CallStat) stat()       {}
func (*DoStat) node()         {}
func (*DoStat) stat()         {}
func (*WhileStat) node()      {}
func (*WhileStat) stat()      {}
func (*RepeatStat) node()     {}
func (*RepeatStat) stat()     {}
func (*IfStat) node()         {}
func (*IfStat) stat()         {}
func (*NumericForStat) node() {}
func (*NumericForStat) stat() {}
func (*GenericForStat) node() {}
func (*GenericForStat) stat() {}
func (*FuncStat) node()       {}
func (*FuncStat) stat()       {}
func (*LocalFuncStat) node()  {}
func (*LocalFuncStat) stat()  {}
func (*ReturnStat) node()     {}
func (*ReturnStat) stat()     {}
func (*BreakStat) node()      {}
func (*BreakStat) stat()      {}
func (*ContinueStat) node()   {}
func (*ContinueStat) stat()   {}
func (*GotoStat) node()       {}
func (*GotoStat) stat()       {}
func (*LabelStat) node()      {}
func (*LabelStat) stat()      {}
func (*TypeAliasStat) node()  {}
func (*TypeAliasStat) stat()  {}

// ---- Expressions ----

type NameExpr struct {
	Name string
	// Resolved binding, set by the scope resolver. nil => free (global). When
	// set and Binding.NewName != "", the printer emits NewName.
	Binding *Binding
}

type NumberExpr struct{ Text string } // verbatim token text

type StringExpr struct{ Text string } // verbatim token text incl. quotes/brackets

type BoolExpr struct{ Val bool }

type NilExpr struct{}

type VarargExpr struct{} // ...

type IndexExpr struct {
	Obj Expr
	// Dot/Method field access uses Field (not renameable). Bracket access uses Key.
	Field    string // non-empty for a.b ; method name for a:b()
	IsMethod bool   // true for a:b
	Key      Expr   // non-nil for a[expr]
}

type CallExpr struct {
	Fn   Expr
	Args []Expr
	// Method call a:b(args) is represented as Fn = IndexExpr{IsMethod:true}.
}

type FuncExpr struct {
	Params   []*NameExpr // renameable, scoped to Body
	IsVararg bool
	Body     []Stat
}

type TableExpr struct{ Fields []TableField }

// TableField is one entry in a table constructor.
//   - Key == nil, KeyName == "": positional ({v})
//   - KeyName != "": name key ({k = v}) — KeyName is NOT renameable
//   - Key != nil: computed key ({[e] = v})
type TableField struct {
	KeyName string
	Key     Expr
	Value   Expr
}

type BinExpr struct {
	Op   string
	L, R Expr
}

type UnExpr struct {
	Op string // "not", "-", "#", "~"
	E  Expr
}

// ParenExpr preserves explicit parentheses where they are semantically
// significant (e.g. truncating a multi-value call to one value).
type ParenExpr struct{ E Expr }

func (*NameExpr) node()   {}
func (*NameExpr) expr()   {}
func (*NumberExpr) node() {}
func (*NumberExpr) expr() {}
func (*StringExpr) node() {}
func (*StringExpr) expr() {}
func (*BoolExpr) node()   {}
func (*BoolExpr) expr()   {}
func (*NilExpr) node()    {}
func (*NilExpr) expr()    {}
func (*VarargExpr) node() {}
func (*VarargExpr) expr() {}
func (*IndexExpr) node()  {}
func (*IndexExpr) expr()  {}
func (*CallExpr) node()   {}
func (*CallExpr) expr()   {}
func (*FuncExpr) node()   {}
func (*FuncExpr) expr()   {}
func (*TableExpr) node()  {}
func (*TableExpr) expr()  {}
func (*BinExpr) node()    {}
func (*BinExpr) expr()    {}
func (*UnExpr) node()     {}
func (*UnExpr) expr()     {}
func (*ParenExpr) node()  {}
func (*ParenExpr) expr()  {}

// Binding is a resolved local/param/forvar declaration shared by its
// declaration NameExpr and all referencing NameExprs.
type Binding struct {
	OrigName string
	NewName  string // set by the renamer; empty means unchanged
}
