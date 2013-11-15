/*
 * Copyright (c) 2013 Dave Collins <dave@davec.name>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package spew_test

import (
	"fmt"
)

// custom type to test Stinger interface on non-pointer receiver.
type stringer string

// String implements the Stringer interface for testing invocation of custom
// stringers on types with non-pointer receivers.
func (s stringer) String() string {
	return "stringer " + string(s)
}

// custom type to test Stinger interface on pointer receiver.
type pstringer string

// String implements the Stringer interface for testing invocation of custom
// stringers on types with only pointer receivers.
func (s *pstringer) String() string {
	return "stringer " + string(*s)
}

// xref1 and xref2 are cross referencing structs for testing circular reference
// detection.
type xref1 struct {
	ps2 *xref2
}
type xref2 struct {
	ps1 *xref1
}

// indirCir1, indirCir2, and indirCir3 are used to generate an indirect circular
// reference for testing detection.
type indirCir1 struct {
	ps2 *indirCir2
}
type indirCir2 struct {
	ps3 *indirCir3
}
type indirCir3 struct {
	ps1 *indirCir1
}

// embed is used to test embedded structures.
type embed struct {
	a string
}

// embedwrap is used to test embedded structures.
type embedwrap struct {
	*embed
	e *embed
}

// panicer is used to intentionally cause a panic for testing spew properly
// handles them
type panicer int

func (p panicer) String() string {
	panic("test panic")
}

// customError is used to test custom error interface invocation.
type customError int

func (e customError) Error() string {
	return fmt.Sprintf("error: %d", int(e))
}

// stringizeWants converts a slice of wanted test output into a format suitable
// for a test error message.
func stringizeWants(wants []string) string {
	s := ""
	for i, want := range wants {
		if i > 0 {
			s += fmt.Sprintf("want%d: %s", i+1, want)
		} else {
			s += "want: " + want
		}
	}
	return s
}

// testFailed returns whether or not a test failed by checking if the result
// of the test is in the slice of wanted strings.
func testFailed(result string, wants []string) bool {
	for _, want := range wants {
		if result == want {
			return false
		}
	}
	return true
}

// TestSortValues ensures the sort functionality for relect.Value based sorting
// works as intended.
func TestSortValues(t *testing.T) {
	v := reflect.ValueOf

	a := v("a")
	b := v("b")
	c := v("c")
	tests := []struct {
		input    []reflect.Value
		expected []reflect.Value
	}{
		{
			[]reflect.Value{v(2), v(1), v(3)},
			[]reflect.Value{v(1), v(2), v(3)},
		},
		{
			[]reflect.Value{v(2.), v(1.), v(3.)},
			[]reflect.Value{v(1.), v(2.), v(3.)},
		},
		{
			[]reflect.Value{v(false), v(true), v(false)},
			[]reflect.Value{v(false), v(false), v(true)},
		},
		{
			[]reflect.Value{b, a, c},
			[]reflect.Value{a, b, c},
		},
	}
	for _, test := range tests {
		spew.SortValues(test.input)
		if !reflect.DeepEqual(test.input, test.expected) {
			t.Errorf("Sort mismatch:\n %v != %v", test.input, test.expected)
		}
	}
}
