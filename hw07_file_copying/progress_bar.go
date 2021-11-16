package main

import (
	"fmt"
	"strings"
)

type Bar struct {
	l     int
	acc   int
	total int
}

func NewBar(total int) *Bar {
	return &Bar{total: total}
}

func (b *Bar) Add(acc int) {
	if b.l > 0 {
		fmt.Print(strings.Repeat("\r", b.l))
	}
	b.acc += acc
	out := "progress: " + fmt.Sprintf("%02d", b.acc*100/b.total) + "%"
	b.l = len(out)
	fmt.Print(out)
}

func (b *Bar) Finish() {
	fmt.Println()
}
