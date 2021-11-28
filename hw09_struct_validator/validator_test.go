package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	AppLen struct {
		Version string `validate:"len:5"`
	}

	AppAlias struct {
		Version UserRole `validate:"len:5"`
	}

	AppMin struct {
		Version         string   `validate:"min:5"`
		VersionInt      int      `validate:"min:5"`
		VersionFloat    float64  `validate:"min:5"`
		VersionIntP     *int     `validate:"min:5"`
		VersionStringP  *string  `validate:"min:5"`
		VersionStringPP **string `validate:"min:5"`
	}

	AppMax struct {
		Version      string  `validate:"max:5"`
		VersionInt   int     `validate:"max:5"`
		VersionFloat float64 `validate:"max:5"`
	}

	AppMinMax struct {
		Version      string  `validate:"max:6|min:4"`
		VersionInt   int     `validate:"min:5|max:5"`
		VersionFloat float64 `validate:"max:6|min:4"`
	}

	AppEmpty struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	AppInt struct {
		VersionStr string `validate:"in:foo,bar,baz"`
		VersionInt int    `validate:"in:1,2,3"`
	}

	AppRegexp struct {
		Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	}

	AppInvalid struct {
		Email string `validate:"FOO:100500"`
	}
)

func TestValidate(t *testing.T) {
	intp := 5
	strp := "12345"
	strpp := &strp
	tests := []struct {
		in          interface{}
		expectedErr bool
	}{
		{AppLen{Version: "12345"}, false},
		{AppLen{Version: "1234"}, true},
		{AppLen{Version: "123456"}, true},

		{AppAlias{Version: "12345"}, false},
		{AppAlias{Version: "1234"}, true},
		{AppAlias{Version: "123456"}, true},

		{AppEmpty{}, false},

		{AppMin{"12345", 5, 5.0, &intp, &strp, &strpp}, false},
		{AppMin{"1234", 5, 5.0, &intp, &strp, &strpp}, true},
		{AppMin{"12345", 4, 5.0, &intp, &strp, &strpp}, true},
		{AppMin{"12345", 5, 4.9, &intp, &strp, &strpp}, true},

		{AppMax{"12345", 5, 5.0}, false},
		{AppMax{"123456", 5, 5.0}, true},
		{AppMax{"12345", 6, 5.0}, true},
		{AppMax{"12345", 5, 5.1}, true},

		{AppMinMax{"1234", 5, 5.0}, false},
		{AppMinMax{"123456", 5, 5.0}, false},
		{AppMinMax{"12345", 5, 4.0}, false},
		{AppMinMax{"12345", 5, 6.0}, false},
		{AppMinMax{"1234567", 5, 5.0}, true},
		{AppMinMax{"123", 5, 5.0}, true},
		{AppMinMax{"12345", 6, 5.0}, true},
		{AppMinMax{"12345", 4, 5.0}, true},
		{AppMinMax{"12345", 5, 6.1}, true},

		{AppInt{"foo", 1}, false},
		{AppInt{"bar", 1}, false},
		{AppInt{"baz", 1}, false},
		{AppInt{"foo", 2}, false},
		{AppInt{"foo", 3}, false},
		{AppInt{"moo", 1}, true},
		{AppInt{"foo", 0}, true},
		{AppInt{"foo", 4}, true},

		{AppRegexp{"foo@example.com"}, false},
		{AppRegexp{"foo@"}, true},
		{AppRegexp{"@example.com"}, true},
		{AppRegexp{"fooexample.com"}, true},

		{AppInvalid{"foo"}, true},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
