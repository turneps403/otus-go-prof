package hw04lrucache_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	hw04lrucache "github.com/turneps403/otus-go-prof/hw04lrucache"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := hw04lrucache.NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := hw04lrucache.NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())
		require.Equal(t, 30, l.Back().Val)
		require.Equal(t, 20, l.Front().Next.Val)
		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Val)
		require.Equal(t, 70, l.Back().Val)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Val.(int))
		}
		require.Equal(t, []int{80, 60, 40, 10, 30, 50, 70}, elems)
	})
}

func TestListMy(t *testing.T) {
	l := hw04lrucache.NewList()
	require.Nil(t, l.Front())
	require.Nil(t, l.Back())

	l.PushFront(100500)
	require.Equal(t, l.Front(), l.Back())

	l.Remove(l.Back())
	require.Nil(t, l.Front())
	require.Nil(t, l.Back())

	l.PushBack(100501)
	require.Equal(t, l.Front(), l.Back())

	l.PushBack(100502)
	require.Equal(t, l.Front().Val, 100501)
	require.Equal(t, l.Back().Val, 100502)
}
