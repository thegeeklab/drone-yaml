// Copyright (c) 2019, Drone IO Inc.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package pretty

import (
	"sort"

	"github.com/drone/drone-yaml/yaml"
)

// helper function pretty prints the container mapping.
func printContainer(w writer, v *yaml.Container) {
	w.IndentIncrease()
	w.WriteTagValue("name", v.Name)
	w.WriteTagValue("pull", v.Pull)
	w.WriteTagValue("image", v.Image)

	if v.Build != nil {
		printBuild(w, v.Build)
	}

	if v.Push != nil {
		w.WriteTagValue("push", v.Push.Image)
	}

	w.WriteTagValue("detach", v.Detach)
	w.WriteTagValue("user", v.User)
	w.WriteTagValue("shell", v.Shell)
	w.WriteTagValue("entrypoint", v.Entrypoint)
	w.WriteTagValue("command", v.Command)
	w.WriteTagValue("commands", v.Commands)
	w.WriteTagValue("dns", v.DNS)
	w.WriteTagValue("dns_search", v.DNSSearch)
	w.WriteTagValue("extra_hosts", v.ExtraHosts)
	w.WriteTagValue("network_mode", v.NetworkMode)

	if len(v.Settings) > 0 {
		printSettings(w, v.Settings)
	}

	if len(v.Environment) > 0 {
		printEnviron(w, v.Environment)
	}

	w.WriteTagValue("failure", v.Failure)
	w.WriteTagValue("privileged", v.Privileged)
	w.WriteTagValue("working_dir", v.WorkingDir)

	if len(v.Devices) > 0 {
		printDeviceMounts(w, v.Devices)
	}

	if len(v.Ports) > 0 {
		printPorts(w, v.Ports)
	}

	if v.Resources != nil {
		printResources(w, v.Resources)
	}

	if len(v.Volumes) > 0 {
		printVolumeMounts(w, v.Volumes)
	}

	if !isConditionsEmpty(v.When) {
		printConditions(w, "when", v.When)
	}

	if len(v.DependsOn) > 0 {
		printDependsOn(w, v.DependsOn)
	}

	_ = w.WriteByte('\n')
	w.IndentDecrease()
}

// helper function pretty prints the build node.
func printBuild(w writer, v *yaml.Build) {
	if shortBuild(v) {
		w.WriteTagValue("build", v.Image)
	} else {
		w.WriteTag("build")
		w.IndentIncrease()
		w.WriteTagValue("image", v.Image)
		w.WriteTagValue("args", v.Args)
		w.WriteTagValue("cache_from", v.CacheFrom)
		w.WriteTagValue("context", v.Context)
		w.WriteTagValue("dockerfile", v.Dockerfile)
		w.WriteTagValue("labels", v.Labels)
		w.IndentDecrease()
	}
}

// helper function pretty prints the depends_on sequence.
func printDependsOn(w writer, v []string) {
	w.WriteTagValue("depends_on", v)
}

// helper function pretty prints the device sequence.
func printDeviceMounts(w writer, v []*yaml.VolumeDevice) {
	w.WriteTag("devices")

	for _, v := range v {
		s := new(indexWriter)
		s.writer = w
		s.IndentIncrease()
		s.WriteTagValue("name", v.Name)
		s.WriteTagValue("path", v.Path)
		s.IndentDecrease()
	}
}

// helper function pretty prints the environment mapping.
func printEnviron(w writer, v map[string]*yaml.Variable) {
	keys := make([]string, 0)

	for k := range v {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	w.WriteTag("environment")
	w.IndentIncrease()

	for _, k := range keys {
		v := v[k]
		if v.FromSecret == "" {
			w.WriteTagValue(k, v.Value)
		} else {
			w.WriteTag(k)
			w.IndentIncrease()
			w.WriteTagValue("from_secret", v.FromSecret)
			w.IndentDecrease()
		}
	}

	w.IndentDecrease()
}

// helper function pretty prints the port sequence.
func printPorts(w writer, v []*yaml.Port) {
	w.WriteTag("ports")

	for _, v := range v {
		if shortPort(v) {
			_ = w.WriteByte('\n')
			w.Indent()
			_ = w.WriteByte('-')
			_ = w.WriteByte(' ')
			writeInt(w, v.Port)

			continue
		}

		s := new(indexWriter)
		s.writer = w
		s.IndentIncrease()
		s.WriteTagValue("port", v.Port)
		s.WriteTagValue("host", v.Host)
		s.WriteTagValue("protocol", v.Protocol)
		s.IndentDecrease()
	}
}

// helper function pretty prints the resource mapping.
func printResources(w writer, v *yaml.Resources) {
	w.WriteTag("resources")
	w.IndentIncrease()

	if v.Limits != nil {
		w.WriteTag("limits")
		w.IndentIncrease()
		w.WriteTagValue("cpu", v.Limits.CPU)
		w.WriteTagValue("memory", v.Limits.Memory)
		w.IndentDecrease()
	}

	if v.Requests != nil {
		w.WriteTag("requests")
		w.IndentIncrease()
		w.WriteTagValue("cpu", v.Requests.CPU)
		w.WriteTagValue("memory", v.Requests.Memory)
		w.IndentDecrease()
	}

	w.IndentDecrease()
}

// helper function pretty prints the resoure mapping.
func printSettings(w writer, v map[string]*yaml.Parameter) {
	keys := make([]string, 0)

	for k := range v {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	w.WriteTag("settings")
	w.IndentIncrease()

	for _, k := range keys {
		v := v[k]
		if v.FromSecret == "" {
			w.IncludeZero()
			w.WriteTagValue(k, v.Value)
			w.ExcludeZero()
		} else {
			w.WriteTag(k)
			w.IndentIncrease()
			w.WriteTagValue("from_secret", v.FromSecret)
			w.IndentDecrease()
		}
	}

	w.IndentDecrease()
}

// helper function pretty prints the volume sequence.
func printVolumeMounts(w writer, v []*yaml.VolumeMount) {
	w.WriteTag("volumes")

	for _, v := range v {
		s := new(indexWriter)
		s.writer = w
		w.IndentIncrease()
		s.IndentIncrease()

		s.WriteTagValue("name", v.Name)
		s.WriteTagValue("path", v.Path)

		s.IndentDecrease()
		w.IndentDecrease()
	}
}

// helper function returns true if the Build block should
// be printed in short form.
func shortBuild(b *yaml.Build) bool {
	return len(b.Args) == 0 &&
		len(b.CacheFrom) == 0 &&
		len(b.Context) == 0 &&
		len(b.Dockerfile) == 0 &&
		len(b.Labels) == 0
}

// helper function returns true if the Port block should
// be printed in short form.
func shortPort(p *yaml.Port) bool {
	return p.Host == 0 && len(p.Protocol) == 0
}
