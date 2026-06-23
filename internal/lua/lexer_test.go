package lua

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []token
	}{
		{
			name: "names numbers operators",
			src:  "local x = 1",
			want: []token{
				{tkName, "local"}, {tkName, "x"}, {tkOp, "="}, {tkNumber, "1"},
			},
		},
		{
			name: "double quoted string with keyword",
			src:  `local x = "function"`,
			want: []token{
				{tkName, "local"}, {tkName, "x"}, {tkOp, "="}, {tkString, `"function"`},
			},
		},
		{
			name: "single quoted string with escape",
			src:  `s = 'a\'b'`,
			want: []token{
				{tkName, "s"}, {tkOp, "="}, {tkString, `'a\'b'`},
			},
		},
		{
			name: "comma inside string",
			src:  `"a, b"`,
			want: []token{{tkString, `"a, b"`}},
		},
		{
			name: "long string level zero",
			src:  `[[a, b end]]`,
			want: []token{{tkString, `[[a, b end]]`}},
		},
		{
			name: "long string level two",
			src:  `[==[x]==]`,
			want: []token{{tkString, `[==[x]==]`}},
		},
		{
			name: "line comment",
			src:  "a -- hi",
			want: []token{{tkName, "a"}, {tkComment, "-- hi"}},
		},
		{
			name: "block comment",
			src:  `--[[ c ]]`,
			want: []token{{tkComment, `--[[ c ]]`}},
		},
		{
			name: "block comment nonzero level",
			src:  `--[==[ c ]==]`,
			want: []token{{tkComment, `--[==[ c ]==]`}},
		},
		{
			name: "numbers hex float exp dotlead",
			src:  "0xFF 1.5 1e3 .5",
			want: []token{
				{tkNumber, "0xFF"}, {tkNumber, "1.5"}, {tkNumber, "1e3"}, {tkNumber, ".5"},
			},
		},
		{
			name: "multi char operators",
			src:  "== ~= <= >= .. ... :: // << >>",
			want: []token{
				{tkOp, "=="}, {tkOp, "~="}, {tkOp, "<="}, {tkOp, ">="}, {tkOp, ".."},
				{tkOp, "..."}, {tkOp, "::"}, {tkOp, "//"}, {tkOp, "<<"}, {tkOp, ">>"},
			},
		},
		{
			name: "index bracket not long string",
			src:  "t[1]",
			want: []token{
				{tkName, "t"}, {tkOp, "["}, {tkNumber, "1"}, {tkOp, "]"},
			},
		},
		{
			name: "greedy ellipsis matches over dotdot",
			src:  "...",
			want: []token{{tkOp, "..."}},
		},
		{
			name: "concat between names no spaces",
			src:  "x..y",
			want: []token{{tkName, "x"}, {tkOp, ".."}, {tkName, "y"}},
		},
		{
			name: "ellipsis between names no spaces",
			src:  "a...b",
			want: []token{{tkName, "a"}, {tkOp, "..."}, {tkName, "b"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, lex(tt.src))
		})
	}
}
