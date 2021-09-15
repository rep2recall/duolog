package duolog

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
)

type Duolog struct {
	// Set filename, if you need to log to file
	Filename string
	// Set to true to disable coloring
	NoColor bool
	// Lexer is syntax detector - consider JSON
	Lexer string
	// Formatter defines how to format and how many colors - consider terminal256
	Formatter string
	// Theme - consider paraiso-dark
	Theme string

	target *os.File
}

func (f Duolog) Write(p []byte) (n int, err error) {
	segs := strings.Split(strings.TrimRight(string(p), "\n"), "\t")

	s := ""
	if len(segs) == 3 {
		if !isObject(segs[1]) {
			segs[1] = ""
		}
		if !isObject(segs[2]) {
			segs[2] = ""
		}
		s = strings.Join(segs, "\t")
	} else {
		s = segs[0]
	}
	if s[len(s)-1] != '\n' {
		s += "\n"
	}

	if s0, e := url.PathUnescape(s); e == nil {
		s = s0
	}

	if f.NoColor {
		fmt.Print(s)
	} else {
		if err := quick.Highlight(os.Stdout, s, f.Lexer, f.Formatter, f.Theme); err != nil {
			log.New(os.Stderr, "error: ", log.LstdFlags).Println(err)
		}
	}

	if f.target != nil {
		return f.target.Write([]byte(s))
	}

	return len([]byte(s)), nil
}

func isObject(s string) bool {
	return len(s) >= 2 && s[0] == '{' && s[len(s)-1] == '}'
}

func (w Duolog) Logger() *log.Logger {
	return log.New(w, "", log.LstdFlags)
}

func New(filename string, cfg Duolog) (Duolog, error) {
	if !cfg.NoColor {
		if cfg.Formatter == "" {
			cfg.Formatter = "terminal256"
		}

		if cfg.Lexer == "" {
			cfg.Lexer = "JSON"
		}

		if cfg.Theme == "" {
			cfg.Theme = "paraiso-dark"
		}
	}

	f, err := os.Create(filename)
	cfg.target = f

	return cfg, err
}
