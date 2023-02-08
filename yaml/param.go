// Copyright (c) 2019, Drone IO Inc.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package yaml

type (
	// Parameter represents an configuration parameter that
	// can be defined as a literal or as a reference
	// to a secret.
	Parameter struct {
		Value      interface{} `json:"value,omitempty"`
		FromSecret string      `json:"from_secret,omitempty" yaml:"from_secret"`
	}

	// parameter is a tempoary type used to unmarshal
	// parameters with references to secrets.
	parameter struct {
		FromSecret string `yaml:"from_secret"`
	}
)

// UnmarshalYAML implements yaml unmarshalling.
func (p *Parameter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	d := new(parameter)
	err := unmarshal(d)

	if err == nil && d.FromSecret != "" {
		p.FromSecret = d.FromSecret

		return nil
	}

	var i interface{}
	err = unmarshal(&i)
	p.Value = i

	return err
}

// MarshalYAML implements yaml marshalling.
func (p *Parameter) MarshalYAML() (interface{}, error) {
	if p.FromSecret != "" {
		m := map[string]interface{}{}
		m["from_secret"] = p.FromSecret

		return m, nil
	}

	if p.Value != "" {
		return p.Value, nil
	}

	//nolint:nilnil
	return nil, nil
}
