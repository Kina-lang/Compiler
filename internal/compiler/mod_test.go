package compiler

import (
	"regexp"
	"testing"
)

// TestGetGreeting calls compiler.GetGreeting with a name, checking
// for a valid return value.
func TestGetGreeting(t *testing.T) {
    name := "Ada"
    want := regexp.MustCompile(`Hello, ` + name + "!")
    msg := GetGreeting(name)

    if !want.MatchString(msg) {
        t.Errorf(`GetGreeting("%s") = %q, want match for %#q`, name, msg, want)
    }
}
