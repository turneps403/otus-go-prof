package myreader

var (
	converterToNonPrinted = map[rune]rune{'n': '\n', 'r': '\r', 't': '\t'}
	NonPrintedToLiteral   = map[rune]string{'\n': "\n", '\r': "\r", '\t': "\t"}
)

type myreader struct {
	chrs []rune
	pos  int
}

func NewMyReader(s string) *myreader {
	return &myreader{[]rune(s), 0}
}

func (r *myreader) HasNext() bool {
	return r.pos < len(r.chrs)
}

func (r *myreader) Next() (rune, int, error) {
	if !r.HasNext() {
		return 0, 0, &ReaderError{reason: "Next() was called after the end of the line or line was empty"}
	}
	chr := r.chrs[r.pos]
	r.pos++

	if '0' <= chr && chr <= '9' {
		return 0, 0, &ReaderError{reason: "sequence has been started with a numeric value"}
	}

	if chr == '\\' {
		if r.HasNext() {
			chr = r.chrs[r.pos]
			r.pos++

			if nonPr, exists := converterToNonPrinted[chr]; exists {
				chr = nonPr
			}

		} else {
			return 0, 0, &ReaderError{reason: "ambiguous use of \\ symbol"}
		}
	}

	rep := 1
	if r.HasNext() {
		repTmp := r.chrs[r.pos]
		if '0' <= repTmp && repTmp <= '9' {
			rep = int(repTmp - '0')
			r.pos++
		}
	}

	return chr, rep, nil
}
