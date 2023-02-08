// Copyright (c) 2019, Drone IO Inc.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package yaml

type (
	// Variable represents an environment variable that
	// can be defined as a string literal or as a reference
	// to a secret.
	Variable struct {
		Value      string `json:"value,omitempty"`
		FromSecret string `json:"from_secret,omitempty" yaml:"from_secret"`
	}

	// variable is a tempoary type used to unmarshal
	// variables with references to secrets.
	variable struct {
		Value      string
		FromSecret string `yaml:"from_secret"`
	}
)

// UnmarshalYAML implements yaml unmarshalling.
func (v *Variable) UnmarshalYAML(unmarshal func(interface{}) error) error {
	d := new(variable)

	err := unmarshal(&d.Value)
	if err != nil {
		err = unmarshal(d)
	}

	v.Value = d.Value
	v.FromSecret = d.FromSecret

	return err
}

// MarshalYAML implements yaml marshalling.
func (v *Variable) MarshalYAML() (interface{}, error) {
	if v.FromSecret != "" {
		m := map[string]interface{}{}
		m["from_secret"] = v.FromSecret

		return m, nil
	}

	if v.Value != "" {
		return v.Value, nil
	}

	//nolint:nilnil
	return nil, nil
}
