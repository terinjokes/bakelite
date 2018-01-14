package flag

import (
	"fmt"
	"strings"
)

// A StringsValue is a command-line flag that interprets its argument
// as a space-separated list of strings.
type StringsValue []string

// Set implements the flag.Value interface by spliting the provided string at
// spaces.
func (v *StringsValue) Set(s string) error {
	if s == "" {
		*v = []string{}
		return nil
	}

	*v = strings.Fields(s)

	return nil
}

// Get implements the flag.Getter interface by returning the contents of this
// value.
func (v *StringsValue) Get() interface{} {
	return []string(*v)
}

// String returns a string
func (v *StringsValue) String() string {
	return fmt.Sprintf("%s", *v)
}
