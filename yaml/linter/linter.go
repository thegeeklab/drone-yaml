// Copyright (c) 2019, Drone IO Inc.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package linter

import (
	"errors"
	"fmt"

	"github.com/drone/drone-yaml/yaml"
)

//nolint:gochecknoglobals
var (
	os = map[string]struct{}{
		"linux":   {},
		"windows": {},
	}
	arch = map[string]struct{}{
		"arm":   {},
		"arm64": {},
		"amd64": {},
	}
)

var (
	// ErrDuplicateStepName is returned when two Pipeline steps
	// have the same name.
	ErrDuplicateStepName = errors.New("linter: duplicate step names")

	// ErrMissingDependency is returned when a Pipeline step
	// defines dependencies that are invalid or unknown.
	ErrMissingDependency = errors.New("linter: invalid or unknown step dependency")

	// ErrCyclicalDependency is returned when a Pipeline step
	// defines a cyclical dependency, which would result in an
	// infinite execution loop.
	ErrCyclicalDependency = errors.New("linter: cyclical step dependency detected")

	ErrUnsupportedOS         = errors.New("linter: unsupported os")
	ErrUnsupportedArch       = errors.New("linter: unsupported architecture")
	ErrInvalidImage          = errors.New("linter: invalid or missing image")
	ErrInvalidBuildImage     = errors.New("linter: invalid or missing build image")
	ErrInvalidName           = errors.New("linter: invalid or missing name")
	ErrPrivilegedNotAllowed  = errors.New("linter: untrusted repositories cannot enable privileged mode")
	ErrMountNotAllowed       = errors.New("linter: untrusted repositories cannot mount devices")
	ErrDNSNotAllowed         = errors.New("linter: untrusted repositories cannot configure dns")
	ErrDNSSearchNotAllowed   = errors.New("linter: untrusted repositories cannot configure dns_search")
	ErrExtraHostsNotAllowed  = errors.New("linter: untrusted repositories cannot configure extra_hosts")
	ErrNetworkModeNotAllowed = errors.New("linter: untrusted repositories cannot configure network_mode")
	ErrInvalidVolumeName     = errors.New("linter: invalid volume name")
	ErrHostPortNotAllowed    = errors.New("linter: untrusted repositories cannot map to a host port")
	ErrHostVolumeNotAllowed  = errors.New("linter: untrusted repositories cannot mount host volumes")
	ErrTempVolumeNotAllowed  = errors.New("linter: untrusted repositories cannot mount in-memory volumes")
)

// Lint performs lint operations for a resource.
func Lint(resource yaml.Resource, trusted bool) error {
	switch v := resource.(type) {
	case *yaml.Cron:
		return v.Validate()
	case *yaml.Pipeline:
		return checkPipeline(v, trusted)
	case *yaml.Secret:
		return v.Validate()
	case *yaml.Registry:
		return v.Validate()
	case *yaml.Signature:
		return v.Validate()
	default:
		return nil
	}
}

func checkPipeline(pipeline *yaml.Pipeline, trusted bool) error {
	err := checkVolumes(pipeline, trusted)
	if err != nil {
		return err
	}

	err = checkPlatform(pipeline.Platform)
	if err != nil {
		return err
	}

	names := map[string]struct{}{}
	if !pipeline.Clone.Disable {
		names["clone"] = struct{}{}
	}

	for _, container := range pipeline.Steps {
		_, ok := names[container.Name]
		if ok {
			return ErrDuplicateStepName
		}

		names[container.Name] = struct{}{}

		err := checkContainer(container, trusted)
		if err != nil {
			return err
		}

		err = checkDeps(container, names)
		if err != nil {
			return err
		}
	}

	for _, container := range pipeline.Services {
		_, ok := names[container.Name]
		if ok {
			return ErrDuplicateStepName
		}

		names[container.Name] = struct{}{}

		err := checkContainer(container, trusted)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkPlatform(platform yaml.Platform) error {
	if v := platform.OS; v != "" {
		_, ok := os[v]
		if !ok {
			return fmt.Errorf("%w: %s", ErrUnsupportedOS, v)
		}
	}

	if v := platform.Arch; v != "" {
		_, ok := arch[v]
		if !ok {
			return fmt.Errorf("%w: %s", ErrUnsupportedArch, v)
		}
	}

	return nil
}

func checkContainer(container *yaml.Container, trusted bool) error {
	err := checkPorts(container.Ports, trusted)
	if err != nil {
		return err
	}

	if container.Build == nil && container.Image == "" {
		return ErrInvalidImage
	}

	if container.Build != nil && container.Build.Image == "" {
		return ErrInvalidBuildImage
	}

	if container.Name == "" {
		return ErrInvalidName
	}

	if trusted && container.Privileged {
		return ErrPrivilegedNotAllowed
	}

	if trusted && len(container.Devices) > 0 {
		return ErrMountNotAllowed
	}

	if trusted && len(container.DNS) > 0 {
		return ErrDNSNotAllowed
	}

	if trusted && len(container.DNSSearch) > 0 {
		return ErrDNSSearchNotAllowed
	}

	if trusted && len(container.ExtraHosts) > 0 {
		return ErrExtraHostsNotAllowed
	}

	if trusted && len(container.NetworkMode) > 0 {
		return ErrNetworkModeNotAllowed
	}

	for _, mount := range container.Volumes {
		switch mount.Name {
		case "workspace", "_workspace", "_docker_socket":
			return fmt.Errorf("%w: %s", ErrInvalidVolumeName, mount.Name)
		}
	}

	return nil
}

func checkPorts(ports []*yaml.Port, trusted bool) error {
	for _, port := range ports {
		err := checkPort(port, trusted)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkPort(port *yaml.Port, trusted bool) error {
	if trusted && port.Host != 0 {
		return ErrHostPortNotAllowed
	}

	return nil
}

func checkVolumes(pipeline *yaml.Pipeline, trusted bool) error {
	for _, volume := range pipeline.Volumes {
		if volume.Temp != nil {
			err := checkEmptyDirVolume(volume.Temp, trusted)
			if err != nil {
				return err
			}
		}

		if volume.Host != nil {
			err := checkHostPathVolume(trusted)
			if err != nil {
				return err
			}
		}

		switch volume.Name {
		case "workspace", "_workspace", "_docker_socket":
			return fmt.Errorf("%w: %s", ErrInvalidVolumeName, volume.Name)
		}
	}

	return nil
}

func checkHostPathVolume(trusted bool) error {
	if trusted {
		return ErrHostVolumeNotAllowed
	}

	return nil
}

func checkEmptyDirVolume(volume *yaml.VolumeEmptyDir, trusted bool) error {
	if trusted && volume.Medium == "memory" {
		return ErrTempVolumeNotAllowed
	}

	return nil
}

func checkDeps(container *yaml.Container, deps map[string]struct{}) error {
	for _, dep := range container.DependsOn {
		_, ok := deps[dep]
		if !ok {
			return ErrMissingDependency
		}

		if container.Name == dep {
			return ErrCyclicalDependency
		}
	}

	return nil
}
