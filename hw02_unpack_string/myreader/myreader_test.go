package myreader_test

// go test -v -test.run=TestReaderValid

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/turneps403/otus-go-prof/hw02_unpack_string/myreader"
)

func TestReaderEmptyContructor(t *testing.T) {
	t.Log("run TestReaderEmptyContructor")

	r := myreader.NewMyReader("")
	require.NotNilf(t, r, "instance")
	require.Falsef(t, r.HasNext(), "hasNext")
}

func TestReaderNonEmptyContructor(t *testing.T) {
	t.Log("run TestReaderNonEmptyContructor")

	r := myreader.NewMyReader("foobar")
	require.NotNilf(t, r, "instance")
	require.Truef(t, r.HasNext(), "hasNext")
	ru, rep, err := r.Next()
	require.NoError(t, err)
	require.Equal(t, ru, 'f')
	require.Equal(t, rep, 1)
}

func TestReaderValid(t *testing.T) {
	t.Log("run TestReaderValid")

	type res struct {
		r   rune
		rep int
	}

	tests := []struct {
		input    string
		expected []res
	}{
		{
			input: "fo",
			expected: []res{
				{r: 'f', rep: 1},
				{r: 'o', rep: 1},
			},
		},
		{
			input: "fo2",
			expected: []res{
				{r: 'f', rep: 1},
				{r: 'o', rep: 2},
			},
		},
		{
			input: "fo4b",
			expected: []res{
				{r: 'f', rep: 1},
				{r: 'o', rep: 4},
				{r: 'b', rep: 1},
			},
		},
		{
			input: "fo4b2",
			expected: []res{
				{r: 'f', rep: 1},
				{r: 'o', rep: 4},
				{r: 'b', rep: 2},
			},
		},
		{
			input: "fo0b",
			expected: []res{
				{r: 'f', rep: 1},
				{r: 'o', rep: 0},
				{r: 'b', rep: 1},
			},
		},
		{
			input: "fЯ0b",
			expected: []res{
				{r: 'f', rep: 1},
				{r: 'Я', rep: 0},
				{r: 'b', rep: 1},
			},
		},
		{
			input: `fo\4b2`,
			expected: []res{
				{r: 'f', rep: 1},
				{r: 'o', rep: 1},
				{r: '4', rep: 1},
				{r: 'b', rep: 2},
			},
		},
		{
			input: `\4\53b2`,
			expected: []res{
				{r: '4', rep: 1},
				{r: '5', rep: 3},
				{r: 'b', rep: 2},
			},
		},
		{
			input: `\n5\t`,
			expected: []res{
				{r: '\n', rep: 5},
				{r: '\t', rep: 1},
			},
		},
		{
			input:    "",
			expected: []res{},
		},
		{
			input: `e\\5a`,
			expected: []res{
				{r: 'e', rep: 1},
				{r: '\\', rep: 5},
				{r: 'a', rep: 1},
			},
		},
	}

	for _, v := range tests {
		v := v
		t.Run(v.input, func(t *testing.T) {
			t.Logf("run test: %s", t.Name())
			r := myreader.NewMyReader(v.input)
			for _, res := range v.expected {
				require.Truef(t, r.HasNext(), "expected longer than have from reader")
				ru, rep, err := r.Next()
				require.Equal(t, ru, res.r)
				require.Equal(t, rep, res.rep)
				require.NoError(t, err)
			}
			require.Falsef(t, r.HasNext(), "outcome from reader longer than expected")
		})
	}

}

func TestReaderInValid(t *testing.T) {
	t.Log("run TestReaderInValid")

	type res struct {
		r   rune
		rep int
		err bool
	}

	tests := []struct {
		input    string
		expected []res
	}{
		{
			input: `fo\`,
			expected: []res{
				{r: 'f', rep: 1, err: false},
				{r: 'o', rep: 1, err: false},
				{r: 0, rep: 0, err: true},
			},
		},
		{
			input: "4fo",
			expected: []res{
				{r: 0, rep: 0, err: true},
			},
		},
	}

	for _, v := range tests {
		v := v
		t.Run(v.input, func(t *testing.T) {
			t.Logf("run test: %s", t.Name())
			r := myreader.NewMyReader(v.input)
			for _, res := range v.expected {
				ru, rep, err := r.Next()
				require.Equal(t, ru, res.r)
				require.Equal(t, rep, res.rep)
				if res.err {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}

}
