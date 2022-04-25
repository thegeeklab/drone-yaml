// Copyright (c) 2019, Drone IO Inc.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package pretty

import (
	"strings"

	"github.com/drone/drone-yaml/yaml"
)

// TODO consider "!!binary |" for secret value

// helper function to pretty prints the signature resource.
func printSecret(w writer, v *yaml.Secret) {
	_, _ = w.WriteString("---")
	w.WriteTagValue("version", v.Version)
	w.WriteTagValue("kind", v.Kind)
	w.WriteTagValue("type", v.Type)

	if len(v.Data) > 0 {
		w.WriteTagValue("name", v.Name)
		printData(w, v.Data)
	}
	if !isSecretGetEmpty(v.Get) {
		w.WriteTagValue("name", v.Name)
		_ = w.WriteByte('\n')
		printGet(w, v.Get)
	}
	_ = w.WriteByte('\n')
	_ = w.WriteByte('\n')
}

// helper function prints the get block.
func printGet(w writer, v yaml.SecretGet) {
	w.WriteTag("get")
	w.IndentIncrease()
	w.WriteTagValue("path", v.Path)
	w.WriteTagValue("name", v.Name)
	w.WriteTagValue("key", v.Key)
	w.IndentDecrease()
}

func printData(w writer, d string) {
	w.WriteTag("data")
	_ = w.WriteByte(' ')
	_ = w.WriteByte('>')
	w.IndentIncrease()
	d = spaceReplacer.Replace(d)
	for _, s := range chunk(d, 60) {
		_ = w.WriteByte('\n')
		w.Indent()
		_, _ = w.WriteString(s)
	}
	w.IndentDecrease()
}

// replace spaces and newlines.
var spaceReplacer = strings.NewReplacer(" ", "", "\n", "")

// helper function returns true if the secret get
// object is empty.
func isSecretGetEmpty(v yaml.SecretGet) bool {
	return v.Key == "" &&
		v.Name == "" &&
		v.Path == ""
}
