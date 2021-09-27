package hw04lrucache

import (
	"fmt"
	"strings"
	"sync"
)

// FOR SUCC PASS: go test -v -count=1 -race -timeout=1m .
var (
	listMutex = sync.RWMutex{}
)

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	fmt.Stringer // for some debug purposes
}

type ListItem struct {
	Val  interface{}
	Next *ListItem
	Prev *ListItem
}

type list struct {
	head   *ListItem
	tail   *ListItem
	length int
}

// Explanation: https://sentry.io/answers/interface-pointer-receiver/
func NewList() List {
	return List(new(list))
}

func (l *list) Len() int {
	listMutex.RLock()
	ret := l.length
	listMutex.RUnlock()
	return ret
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	listMutex.Lock()
	el := &ListItem{Val: v, Next: l.head}
	if l.head != nil {
		l.head.Prev = el
	}
	l.head = el
	if l.Back() == nil {
		l.tail = l.head
	}
	l.length++
	listMutex.Unlock()
	return l.Front()
}

func (l *list) PushBack(v interface{}) *ListItem {
	listMutex.Lock()
	if l.Back() == nil {
		l.tail = &ListItem{Val: v}
		l.head = l.tail
	} else {
		l.tail.Next = &ListItem{Val: v, Prev: l.tail}
		l.tail = l.tail.Next
	}
	l.length++
	listMutex.Unlock()
	return l.Back()
}

func (l *list) Remove(it *ListItem) {
	listMutex.Lock()
	if it.Prev != nil {
		it.Prev.Next = it.Next
	}
	if it.Next != nil {
		it.Next.Prev = it.Prev
	}
	l.length--
	if l.head == it {
		l.head = it.Next
	}
	if l.tail == it {
		l.tail = it.Prev
	}
	listMutex.Unlock()
}

func (l *list) String() string { // reuse knowledge from the `hw01_hello_otus`
	var sb strings.Builder
	fmt.Fprint(&sb, "[")
	t := l.head
	for t != nil {
		fmt.Fprint(&sb, t.Val)
		if t.Next != nil {
			fmt.Fprint(&sb, " -> ")
		}
		t = t.Next
	}
	fmt.Fprint(&sb, "]")
	return sb.String()
}
