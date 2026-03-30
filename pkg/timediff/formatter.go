package timediff

import "time"

// Formatter pre-computes configuration for repeated formatting calls.
// Use NewFormatter to create one, then call Format or AppendFormat per input.
// It captures time.Now() once at construction; all calls use that snapshot.
type Formatter struct {
	now    time.Time
	syntax Syntax
	short  bool
	parts  int
	opts   Options
	skip   []Unit
}

// NewFormatter creates a Formatter with the given options.
func NewFormatter(opts ...Option) *Formatter {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}

	f := &Formatter{
		now:    time.Now(),
		syntax: cfg.syntax,
		short:  cfg.short,
		parts:  cfg.parts,
		opts:   cfg.options,
		skip:   cfg.skip,
	}

	if cfg.nowOverride != nil {
		f.now = *cfg.nowOverride
	}
	if cfg.other != nil {
		f.now = *cfg.other
	}

	if f.syntax == SyntaxRelativeAuto {
		if cfg.other != nil {
			f.syntax = SyntaxRelativeToOther
		} else {
			f.syntax = SyntaxRelativeToNow
		}
	}

	return f
}

// Format returns the human-readable time difference between from and the
// Formatter's reference time.
func (f *Formatter) Format(from time.Time) string {
	iv := Diff(from, f.now)
	return formatInterval(iv, f.syntax, f.short, f.parts, f.opts, f.skip)
}

// AppendFormat appends the formatted time difference to b and returns the
// extended buffer.
func (f *Formatter) AppendFormat(b []byte, from time.Time) []byte {
	iv := Diff(from, f.now)
	return appendInterval(b, iv, f.syntax, f.short, f.parts, f.opts, f.skip)
}
