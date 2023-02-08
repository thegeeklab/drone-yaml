// Copyright (c) 2019, Drone IO Inc.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package yaml

import "errors"

var ErrInvlaidCronBranch = errors.New("yaml: invalid cron branch")

type (
	// Cron is a resource that defines a cron job, used
	// to execute pipelines at scheduled intervals.
	Cron struct {
		Version string `json:"version,omitempty"`
		Kind    string `json:"kind,omitempty"`
		Type    string `json:"type,omitempty"`
		Name    string `json:"name,omitempty"`

		Spec CronSpec `json:"spec,omitempty"`
	}

	// CronSpec defines the cron job.
	CronSpec struct {
		Schedule   string         `json:"schedule,omitempty"`
		Branch     string         `json:"branch,omitempty"`
		Deployment CronDeployment `json:"deployment,omitempty" yaml:"deployment"`
	}

	// CronDeployment defines a cron job deployment.
	CronDeployment struct {
		Target string `json:"target,omitempty"`
	}
)

// GetVersion returns the resource version.
func (c *Cron) GetVersion() string { return c.Version }

// GetKind returns the resource kind.
func (c *Cron) GetKind() string { return c.Kind }

// Validate returns an error if the cron is invalid.
func (c Cron) Validate() error {
	switch {
	case c.Spec.Branch == "":
		return ErrInvlaidCronBranch
	default:
		return nil
	}
}
