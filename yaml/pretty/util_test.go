// Copyright (c), the Drone Authors.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package pretty

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestQuoted(t *testing.T) {
	tests := []struct {
		before, after string
	}{
		{"", `""`},
		{"foo", "foo"},

		// special characters only quoted when followed
		// by whitespace.
		{"&foo", "&foo"},
		{"!foo", "!foo"},
		{"-foo", "-foo"},
		{":foo", ":foo"},

		{"& foo", `"& foo"`},
		{"! foo", `"! foo"`},
		{"- foo", `"- foo"`},
		{": foo", `": foo"`},

		{" & foo", `" & foo"`},
		{" ! foo", `" ! foo"`},
		{" - foo", `" - foo"`},
		{" : foo", `" : foo"`},

		// special characters only quoted when it is the
		// first character in the string.
		{",foo", `",foo"`},
		{"[foo", `"[foo"`},
		{"]foo", `"]foo"`},
		{"{foo", `"{foo"`},
		{"}foo", `"}foo"`},
		{"*foo", `"*foo"`},
		{`"foo`, `"\"foo"`},
		{`'foo`, `"'foo"`},
		{`%foo`, `"%foo"`},
		{`@foo`, `"@foo"`},
		{`|foo`, `"|foo"`},
		{`>foo`, `">foo"`},
		{`#foo`, `"#foo"`},

		{`foo:bar`, `foo:bar`},
		{`foo :bar`, `foo :bar`},
		{`foo: bar`, `"foo: bar"`},
		{`foo:`, `"foo:"`},
		{`alpine:3.8`, `alpine:3.8`}, // verify docker image names are ok

		// comments should be escaped. A comment is a pound
		// sybol preceded by a space.
		{`foo#bar`, `foo#bar`},
		{`foo #bar`, `"foo #bar"`},

		// strings with newlines and control characters
		// should be escaped
		{"foo\nbar", "\"foo\\nbar\""},
	}

	for _, test := range tests {
		buf := new(baseWriter)
		writeEncode(buf, test.before)
		a := test.after
		b := buf.String()

		if b != a {
			t.Errorf("Want %q, got %q", a, b)
		}
	}
}

func TestChunk(t *testing.T) {
	testChunk := []string{
		"ZDllMjFjZDg3Zjk0ZWFjZDRhMjdhMTA1ZDQ1OTVkYTA1ODBjMTk0ZWVlZjQyNmU4",
		"N2RiNTIwZjg0NWQwYjcyYjE3MmFmZDIyYzg3NTQ1N2YyYzgxODhjYjJmNDhhOTFj",
		"ZjdhMzA0YjEzYWFlMmYxMTIwMmEyM2Q1YjQ5Yjg2ZmMK",
	}

	s := strings.Join(testChunk, "")
	got, want := chunk(s, 64), testChunk

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected chunk value")
		t.Log(diff)
	}
}
