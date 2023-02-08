// Copyright (c), the Drone Authors.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package main

import (
	"bytes"
	"io"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/linter"
	"github.com/drone/drone-yaml/yaml/pretty"
)

//nolint:gochecknoglobals
var (
	format     = kingpin.Command("fmt", "format the yaml file")
	formatSave = format.Flag("save", "save result to source").Short('s').Bool()
	formatFile = format.Arg("source", "source file location").Default(".drone.yml").File()

	lint     = kingpin.Command("lint", "lint the yaml file")
	lintPriv = lint.Flag("privileged", "privileged mode").Short('p').Bool()
	lintFile = lint.Arg("source", "source file location").Default(".drone.yml").File()
)

const DefaultFilePerm = 0o640

func main() {
	switch kingpin.Parse() {
	case format.FullCommand():
		kingpin.FatalIfError(runFormat(), "")
	case lint.FullCommand():
		kingpin.FatalIfError(runLint(), "")
	}
}

func runFormat() error {
	f := *formatFile

	m, err := yaml.Parse(f)
	if err != nil {
		return err
	}

	b := new(bytes.Buffer)
	pretty.Print(b, m)

	if *formatSave {
		return os.WriteFile(f.Name(), b.Bytes(), DefaultFilePerm)
	}

	_, err = io.Copy(os.Stderr, b)

	return err
}

func runLint() error {
	f := *lintFile

	m, err := yaml.Parse(f)
	if err != nil {
		return err
	}

	for _, r := range m.Resources {
		err := linter.Lint(r, *lintPriv)
		if err != nil {
			return err
		}
	}

	return nil
}
