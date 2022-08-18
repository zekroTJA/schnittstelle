package schnittstelle

import (
	"fmt"
	"strings"
)

// MultiError is an array of errors which can be handled
// and one error object.
type MultiError []error

// Error assembles a single error message out of all
// errors in the array.
//
// If the array is empty, an empty string is returned.
func (t MultiError) Error() string {
	if len(t) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, err := range t {
		fmt.Fprintf(&sb, "[%d] %s\n", i, err.Error())
	}

	return sb.String()
}

// ToError returns nil if the length of the error
// array is 0 and returns the array if not.
//
// This should be called before returning the
// MultiError in a function to ensure that
// subsequent nil checks perform as expected.
func (t MultiError) ToError() error {
	if len(t) == 0 {
		return nil
	}
	return t
}
