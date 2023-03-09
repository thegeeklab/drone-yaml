// Copyright (c) 2019, Drone IO Inc.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package yaml

import (
	"bytes"
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v3"
)

// Resource enums.
const (
	KindCron      = "cron"
	KindPipeline  = "pipeline"
	KindRegistry  = "registry"
	KindSecret    = "secret"
	KindSignature = "signature"
)

var ErrMarshalNotImplemented = errors.New("yaml: marshal not implemented")

type (
	// Manifest is a collection of Drone resources.
	Manifest struct {
		Resources []Resource
	}

	// Resource represents a Drone resource.
	Resource interface {
		// GetVersion returns the resource version.
		GetVersion() string

		// GetKind returns the resource kind.
		GetKind() string
	}

	// RawResource is a raw encoded resource with the
	// resource kind and type extracted.
	RawResource struct {
		Version string
		Kind    string
		Type    string
		Data    []byte `yaml:"-"`
	}

	//nolint:musttag
	resource struct {
		Version string
		Kind    string `json:"kind"`
		Type    string `json:"type"`
	}
)

// UnmarshalJSON implement the json.Unmarshaler.
func (m *Manifest) UnmarshalJSON(b []byte) error {
	messages := []json.RawMessage{}

	err := json.Unmarshal(b, &messages)
	if err != nil {
		return err
	}

	for _, message := range messages {
		res := new(resource)

		err := json.Unmarshal(message, res)
		if err != nil {
			return err
		}

		var obj Resource

		switch res.Kind {
		case "cron":
			obj = new(Cron)
		case "secret":
			obj = new(Secret)
		case "signature":
			obj = new(Signature)
		case "registry":
			obj = new(Registry)
		default:
			obj = new(Pipeline)
		}

		err = json.Unmarshal(message, obj)
		if err != nil {
			return err
		}

		m.Resources = append(m.Resources, obj)
	}

	return nil
}

// MarshalJSON implement the json.Marshaler.
func (m *Manifest) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Resources)
}

// MarshalYAML is not implemented and returns an error. This is
// because the Manifest is a representation of multiple yaml
// documents, and MarshalYAML would otherwise attempt to marshal
// as a single Yaml document. Use the Encode method instead.
func (m *Manifest) MarshalYAML() (interface{}, error) {
	return nil, ErrMarshalNotImplemented
}

// Encode encodes the manifest in Yaml format.
func (m *Manifest) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := yaml.NewEncoder(buf)

	for _, res := range m.Resources {
		if err := enc.Encode(res); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
