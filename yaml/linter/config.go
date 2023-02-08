// Copyright (c) 2019, Drone IO Inc.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package linter

import (
	"errors"

	"github.com/drone/drone-yaml/yaml"
)

// ErrDuplicatePipelineName is returned when two Pipeline
// resources have the same name.
var ErrDuplicatePipelineName = errors.New("linter: duplicate pipeline names")

// ErrMissingPipelineDependency is returned when a Pipeline
// defines dependencies that are invalid or unknown.
var ErrMissingPipelineDependency = errors.New("linter: invalid or unknown pipeline dependency")

// ErrCyclicalPipelineDependency is returned when a Pipeline
// defines a cyclical dependency, which would result in an
// infinite execution loop.
var ErrCyclicalPipelineDependency = errors.New("linter: cyclical pipeline dependency detected")

// ErrPipelineSelfDependency is returned when a Pipeline
// defines a dependency on itself.
var ErrPipelineSelfDependency = errors.New("linter: pipeline cannot have a dependency on itself")

// Manifest performs lint operations for a manifest.
func Manifest(manifest *yaml.Manifest, trusted bool) error {
	return checkPipelines(manifest)
}

func checkPipelines(manifest *yaml.Manifest) error {
	names := map[string]struct{}{}

	for _, resource := range manifest.Resources {
		switch v := resource.(type) {
		case *yaml.Pipeline:
			_, ok := names[v.Name]
			if ok {
				return ErrDuplicatePipelineName
			}

			names[v.Name] = struct{}{}

			err := checkPipelineDeps(v, names)
			if err != nil {
				return err
			}

			if (v.Kind == "pipeline" || v.Kind == "") && (v.Type == "" || v.Type == "docker") {
				err = checkPlatform(v.Platform)
				if err != nil {
					return err
				}
			}
		default:
			continue
		}
	}

	return nil
}

func checkPipelineDeps(pipeline *yaml.Pipeline, deps map[string]struct{}) error {
	for _, dep := range pipeline.DependsOn {
		_, ok := deps[dep]
		if !ok {
			return ErrMissingPipelineDependency
		}

		if pipeline.Name == dep {
			return ErrPipelineSelfDependency
		}
	}

	return nil
}
