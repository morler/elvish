package eval

import (
	"testing"

	"src.elv.sh/pkg/parse"
)

var benchmarks = []struct {
	name string
	code string
}{
	{"empty", ""},
	{"nop", "nop"},
	{"nop-nop", "nop | nop"},
	{"put-x", "put x"},
	{"for-100", "for x [(range 100)] { }"},
	{"range-100", "range 100 | each {|_| }"},
	{"read-local", "var x = val; nop $x"},
	{"read-upval", "var x = val; { nop $x }"},
}

// Benchmarks specifically for compilation performance optimizations
var compilationBenchmarks = []struct {
	name string
	code string
}{
	{"command-with-options", "echo &option=value"},
	{"command-multiple-options", "echo &a=1 &b=2 &c=3 &d=4 &e=5"},
	{"tilde-simple", "put ~"},
	{"tilde-username", "put ~user"},
	{"tilde-path", "put ~/path/to/file"},
	{"tilde-multiple", "put ~ ~/a ~/b ~/c ~/d"},
	{"tilde-complex", "put ~/a/b/c ~/x/y/z ~/foo/bar/baz"},
	{"mixed-tilde-options", "echo &debug=$true ~/config/app.conf"},
}

func BenchmarkEval(b *testing.B) {
	for _, bench := range benchmarks {
		b.Run(bench.name, func(b *testing.B) {
			ev := NewEvaler()
			src := parse.Source{Name: "[benchmark]", Code: bench.code}

			tree, err := parse.Parse(src, parse.Config{})
			if err != nil {
				panic(err)
			}
			op, _, err := compile(ev.builtin.static(), ev.global.static(), nil, tree, nil)
			if err != nil {
				panic(err)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				fm, cleanup := ev.prepareFrame(src, EvalCfg{Global: ev.Global()})
				_, exec := op.prepare(fm)
				_ = exec()
				cleanup()
			}
		})
	}
}

// BenchmarkCompilation focuses specifically on compilation performance
func BenchmarkCompilation(b *testing.B) {
	for _, bench := range compilationBenchmarks {
		b.Run(bench.name, func(b *testing.B) {
			ev := NewEvaler()
			src := parse.Source{Name: "[compilation-benchmark]", Code: bench.code}

			for i := 0; i < b.N; i++ {
				tree, err := parse.Parse(src, parse.Config{})
				if err != nil {
					b.Fatal(err)
				}
				
				b.StartTimer()
				_, _, err = compile(ev.builtin.static(), ev.global.static(), nil, tree, nil)
				b.StopTimer()
				
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkTildeExpansion focuses specifically on tilde expansion optimization
func BenchmarkTildeExpansion(b *testing.B) {
	tildeTests := []struct {
		name string
		code string
	}{
		{"single-tilde", "put ~"},
		{"tilde-array-small", "put [~ ~/a ~/b]"},
		{"tilde-array-large", "put [~ ~/a ~/b ~/c ~/d ~/e ~/f ~/g ~/h ~/i ~/j]"},
		{"no-tilde", "put /abs/path /another/path"},
		{"mixed-paths", "put ~ /abs ~/rel /another"},
	}
	
	for _, bench := range tildeTests {
		b.Run(bench.name, func(b *testing.B) {
			ev := NewEvaler()
			src := parse.Source{Name: "[tilde-benchmark]", Code: bench.code}

			tree, err := parse.Parse(src, parse.Config{})
			if err != nil {
				b.Fatal(err)
			}
			op, _, err := compile(ev.builtin.static(), ev.global.static(), nil, tree, nil)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				fm, cleanup := ev.prepareFrame(src, EvalCfg{Global: ev.Global()})
				_, exec := op.prepare(fm)
				_ = exec()
				cleanup()
			}
		})
	}
}
