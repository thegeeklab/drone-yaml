// Copyright (c) 2019, Drone IO Inc.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package pretty

import (
	"github.com/drone/drone-yaml/yaml"
)

// helper function pretty prints the signature resource.
func printSignature(w writer, v *yaml.Signature) {
	w.WriteString("---")
	w.WriteTagValue("version", v.Version)
	w.WriteTagValue("kind", v.Kind)
	w.WriteTagValue("hmac", v.Hmac)
	w.WriteByte('\n')
	w.WriteByte('\n')
}
