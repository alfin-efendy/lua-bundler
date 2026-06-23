package obfuscator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/alfin-efendy/lua-bundler/internal/lua"
)

// Obfuscator transforms Lua source by level. It is stateless across units
// except for the per-instance string-encryption key (level 3), which must be
// shared by every unit and the injected decoder.
type Obfuscator struct {
	level int
	key   byte // string-encryption key, set when level >= 3
}

// NewObfuscator creates an obfuscator clamped to levels 1..3.
func NewObfuscator(level int) *Obfuscator {
	if level < 1 {
		level = 1
	}
	if level > 3 {
		level = 3
	}
	o := &Obfuscator{level: level}
	if level >= 3 {
		o.key = randomKey()
	}
	return o
}

// Obfuscate applies the configured level to one independent unit of Lua source.
// On any parse failure it falls back to lexer-only minification.
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
	if o.level >= 3 {
		lua.EncryptStrings(chunk, o.key)
	}
	return chunk.Print()
}

// DecoderPrelude returns the `local _d=...` definition that decodes strings
// encrypted at level 3, or "" for levels 1-2. The bundler injects it once at
// the top of the bundle.
func (o *Obfuscator) DecoderPrelude() string {
	if o.level < 3 {
		return ""
	}
	return lua.DecoderPrelude(o.key)
}

// randomKey returns a non-zero byte for XOR string encryption.
func randomKey() byte {
	n, err := rand.Int(rand.Reader, big.NewInt(255))
	if err != nil {
		return 0x5a // deterministic fallback; crypto/rand does not fail on Linux
	}
	return byte(n.Int64()) + 1 // 1..255
}
