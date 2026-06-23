package obfuscator

import (
	"fmt"
	"os"

	"github.com/alfin-efendy/lua-bundler/internal/lua"
)

// Obfuscator transforms Lua source by level. It holds no cross-unit state:
// every Obfuscate call is independent (per-module/per-entry isolation).
type Obfuscator struct {
	level int // 1 = minify, 2 = +rename, 3 = +rename (string encryption: Phase 2)
}

// NewObfuscator creates an obfuscator clamped to levels 1..3.
func NewObfuscator(level int) *Obfuscator {
	if level < 1 {
		level = 1
	}
	if level > 3 {
		level = 3
	}
	return &Obfuscator{level: level}
}

// Obfuscate applies the configured level to one independent unit of Lua source.
// On any parse failure it falls back to lexer-only minification, so bundling
// never fails because of obfuscation.
func (o *Obfuscator) Obfuscate(code string) string {
	if o.level == 1 {
		return lua.Minify(code)
	}
	chunk, err := lua.Parse(code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  obfuscate: parse failed, falling back to minify: %v\n", err)
		return lua.Minify(code)
	}
	lua.Rename(chunk)
	// Level 3 string encryption is added in Phase 2.
	return chunk.Print()
}
