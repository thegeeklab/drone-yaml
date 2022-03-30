// Copyright (c), the Drone Authors.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package pretty

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

// this unit tests pretty prints a complex yaml structure
// to ensure we have common use cases covered.
func TestWriteComplexValue(t *testing.T) {
	block := map[interface{}]interface{}{}
	err := yaml.Unmarshal([]byte(testComplexValue), &block)
	if err != nil {
		t.Error(err)
		return
	}

	b := new(baseWriter)
	writeValue(b, block)
	got, want := b.String(), strings.TrimSpace(testComplexValue)
	if got != want {
		t.Errorf("Unexpected block format")
		println(got)
		println("---")
		println(want)
	}
}

var testComplexValue = `
a: b
c:
- d
- e
f:
  g: h
  i:
  - j
  - k
  - l: m
    o: p
    q:
    - r
    - s: ~
  - {}
  - []
  - ~
t: {}
u: []
v: 1
w: true
x: ~
z: "#y"
zz: "\nz\n"
"{z}": z`
